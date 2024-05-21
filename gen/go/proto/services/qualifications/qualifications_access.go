package qualifications

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) GetQualificationAccess(ctx context.Context, req *GetQualificationAccessRequest) (*GetQualificationAccessResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	access, err := s.getQualificationAccess(ctx, req.QualificationId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	for i := 0; i < len(access.Jobs); i++ {
		s.enricher.EnrichJobInfo(access.Jobs[i])
	}

	resp := &GetQualificationAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetQualificationAccess(ctx context.Context, req *SetQualificationAccessRequest) (*SetQualificationAccessResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "SetQualificationAccess",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.handleQualificationAccessChanges(ctx, tx, req.Mode, req.QualificationId, req.Access); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &SetQualificationAccessResponse{}, nil
}

func (s *Server) getQualificationAccess(ctx context.Context, qualificationId uint64) (*qualifications.QualificationAccess, error) {
	tQJobAccess := table.FivenetQualificationsJobAccess.AS("qualificationjobaccess")
	jobStmt := tQJobAccess.
		SELECT(
			tQJobAccess.ID,
			tQJobAccess.QualificationID,
			tQJobAccess.Job,
			tQJobAccess.MinimumGrade,
			tQJobAccess.Access,
		).
		FROM(
			tQJobAccess,
		).
		WHERE(
			tQJobAccess.QualificationID.EQ(jet.Uint64(qualificationId)),
		).
		ORDER_BY(
			tQJobAccess.ID.ASC(),
		)

	var jobAccess []*qualifications.QualificationJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &qualifications.QualificationAccess{
		Jobs: jobAccess,
	}, nil
}

func (s *Server) handleQualificationAccessChanges(ctx context.Context, tx qrm.DB, mode qualifications.AccessLevelUpdateMode, qualificationId uint64, access *qualifications.QualificationAccess) error {

	switch mode {
	case qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED:
		fallthrough
	case qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE:
		// Get existing job and user accesses from database
		current, err := s.getQualificationAccess(ctx, qualificationId)
		if err != nil {
			return err
		}

		toCreate, toUpdate, toDelete := s.compareQualificationAccess(current, access)

		if err := s.createQualificationAccess(ctx, tx, qualificationId, toCreate); err != nil {
			return err
		}

		if err := s.updateQualificationAccess(ctx, tx, qualificationId, toUpdate); err != nil {
			return err
		}

		if err := s.deleteQualificationAccess(ctx, tx, qualificationId, toDelete); err != nil {
			return err
		}

	case qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_DELETE:
		if err := s.deleteQualificationAccess(ctx, tx, qualificationId, access); err != nil {
			return err
		}

	case qualifications.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_CLEAR:
		if err := s.clearQualificationAccess(ctx, tx, qualificationId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareQualificationAccess(current, in *qualifications.QualificationAccess) (toCreate *qualifications.QualificationAccess, toUpdate *qualifications.QualificationAccess, toDelete *qualifications.QualificationAccess) {
	toCreate = &qualifications.QualificationAccess{}
	toUpdate = &qualifications.QualificationAccess{}
	toDelete = &qualifications.QualificationAccess{}

	if current == nil || len(current.Jobs) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current.Jobs, func(a, b *qualifications.QualificationJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current.Jobs) == 0 {
		toCreate.Jobs = in.Jobs
	} else {
		foundTracker := []int{}
		for _, cj := range current.Jobs {
			var found *qualifications.QualificationJobAccess
			var foundIdx int
			for i, uj := range in.Jobs {
				if cj.Job != uj.Job {
					continue
				}
				if cj.MinimumGrade != uj.MinimumGrade {
					continue
				}
				found = uj
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete.Jobs = append(toDelete.Jobs, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.MinimumGrade != found.MinimumGrade {
				cj.MinimumGrade = found.MinimumGrade
				changed = true
			}
			if cj.Access != found.Access {
				cj.Access = found.Access
				changed = true
			}

			if changed {
				toUpdate.Jobs = append(toUpdate.Jobs, cj)
			}
		}

		for i, uj := range in.Jobs {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate.Jobs = append(toCreate.Jobs, uj)
			}
		}
	}

	return
}

func (s *Server) createQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *qualifications.QualificationAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tQJobAccess := table.FivenetQualificationsJobAccess
			stmt := tQJobAccess.
				INSERT(
					tQJobAccess.QualificationID,
					tQJobAccess.Job,
					tQJobAccess.MinimumGrade,
					tQJobAccess.Access,
				).
				VALUES(
					qualificationId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) updateQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *qualifications.QualificationAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tQJobAccess := table.FivenetQualificationsJobAccess
			stmt := tQJobAccess.
				UPDATE(
					tQJobAccess.QualificationID,
					tQJobAccess.Job,
					tQJobAccess.MinimumGrade,
					tQJobAccess.Access,
				).
				SET(
					qualificationId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				).
				WHERE(
					tQJobAccess.ID.EQ(jet.Uint64(access.Jobs[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) deleteQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *qualifications.QualificationAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil && len(access.Jobs) > 0 {
		jobIds := []jet.Expression{}
		for i := 0; i < len(access.Jobs); i++ {
			if access.Jobs[i].Id == 0 {
				continue
			}
			jobIds = append(jobIds, jet.Uint64(access.Jobs[i].Id))
		}

		tQJobAccess := table.FivenetQualificationsJobAccess
		stmt := tQJobAccess.
			DELETE().
			WHERE(jet.AND(
				tQJobAccess.ID.IN(jobIds...),
				tQJobAccess.QualificationID.EQ(jet.Uint64(qualificationId)),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64) error {
	jobStmt := tQJobAccess.
		DELETE().
		WHERE(tQJobAccess.QualificationID.EQ(jet.Uint64(qualificationId)))

	if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) checkIfHasAccess(levels []string, userInfo *userinfo.UserInfo, creatorJob string, creator *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
	// and not the rank attributes checks
	if creatorJob != userInfo.Job {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if creator == nil {
		return true
	}

	// If no levels set, assume "Own" as default
	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	if slices.Contains(levels, "Any") {
		return true
	}
	if slices.Contains(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
