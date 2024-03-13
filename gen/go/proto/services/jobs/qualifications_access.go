package jobs

import (
	"context"
	"errors"
	"slices"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) GetQualificationAccess(ctx context.Context, req *GetQualificationAccessRequest) (*GetQualificationAccessResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsjobs.ErrFailedQuery
	}

	access, err := s.getQualificationAccess(ctx, req.QualificationId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
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
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.qualifications.id", int64(req.QualificationId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "SetQualificationAccess",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	ok, err := s.checkIfUserHasAccessToQuali(ctx, req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.handleQualificationAccessChanges(ctx, tx, req.Mode, req.QualificationId, req.Access); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &SetQualificationAccessResponse{}, nil
}

func (s *Server) handleQualificationAccessChanges(ctx context.Context, tx qrm.DB, mode jobs.AccessLevelUpdateMode, qualificationId uint64, access *jobs.QualificationAccess) error {
	// Get existing job and user accesses from database
	current, err := s.getQualificationAccess(ctx, qualificationId)
	if err != nil {
		return err
	}

	switch mode {
	case jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED:
		fallthrough
	case jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE:
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

	case jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_DELETE:
		if err := s.deleteQualificationAccess(ctx, tx, qualificationId, access); err != nil {
			return err
		}

	case jobs.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_CLEAR:
		if err := s.clearQualificationAccess(ctx, tx, qualificationId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareQualificationAccess(current, in *jobs.QualificationAccess) (toCreate *jobs.QualificationAccess, toUpdate *jobs.QualificationAccess, toDelete *jobs.QualificationAccess) {
	toCreate = &jobs.QualificationAccess{}
	toUpdate = &jobs.QualificationAccess{}
	toDelete = &jobs.QualificationAccess{}

	if current == nil || (len(current.Jobs) == 0 && len(current.Requirements) == 0) {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current.Jobs, func(a, b *jobs.QualificationJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current.Jobs) == 0 {
		toCreate.Jobs = in.Jobs
	} else {
		foundTracker := []int{}
		for _, cj := range current.Jobs {
			var found *jobs.QualificationJobAccess
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

	if len(current.Requirements) == 0 {
		toCreate.Requirements = in.Requirements
	} else {
		foundTracker := []int{}
		for _, cj := range current.Requirements {
			var found *jobs.QualificationRequirementsAccess
			var foundIdx int
			for i, uj := range in.Requirements {
				if cj.QualificationId != uj.QualificationId {
					continue
				}
				found = uj
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete.Requirements = append(toDelete.Requirements, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.Access != found.Access {
				cj.Access = found.Access
				changed = true
			}

			if changed {
				toUpdate.Requirements = append(toUpdate.Requirements, cj)
			}
		}

		for i, uj := range in.Requirements {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate.Requirements = append(toCreate.Requirements, uj)
			}
		}
	}

	return
}

func (s *Server) getQualificationAccess(ctx context.Context, qualificationId uint64) (*jobs.QualificationAccess, error) {
	tQJobAccess := table.FivenetJobsQualificationsJobAccess.AS("qualificationjobaccess")
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

	var jobAccess []*jobs.QualificationJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	tQReqAccess := table.FivenetJobsQualificationsReqsAccess.AS("qualificationrequirementsaccess")
	userStmt := tQReqAccess.
		SELECT(
			tQReqAccess.ID,
			tQReqAccess.QualificationID,
			tQReqAccess.TargetQualificationID,
			tQReqAccess.Access,
			tQuali.ID,
			tQuali.Title,
			tQuali.Abbreviation,
		).
		FROM(
			tQReqAccess.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQReqAccess.TargetQualificationID),
				),
		).
		WHERE(
			tQReqAccess.QualificationID.EQ(jet.Uint64(qualificationId)),
		).
		ORDER_BY(
			tQReqAccess.ID.ASC(),
		)

	var reqsAccess []*jobs.QualificationRequirementsAccess
	if err := userStmt.QueryContext(ctx, s.db, &reqsAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &jobs.QualificationAccess{
		Jobs:         jobAccess,
		Requirements: reqsAccess,
	}, nil
}

func (s *Server) createQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *jobs.QualificationAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tQJobAccess := table.FivenetJobsQualificationsJobAccess
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

	if access.Requirements != nil {
		for k := 0; k < len(access.Requirements); k++ {
			// Create document user access
			tQReqAccess := table.FivenetJobsQualificationsReqsAccess
			stmt := tQReqAccess.
				INSERT(
					tQReqAccess.QualificationID,
					tQReqAccess.Access,
				).
				VALUES(
					qualificationId,
					access.Requirements[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) updateQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *jobs.QualificationAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tQJobAccess := table.FivenetJobsQualificationsJobAccess
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

	if access.Requirements != nil {
		for k := 0; k < len(access.Requirements); k++ {
			// Create document user access
			tQReqAccess := table.FivenetJobsQualificationsReqsAccess
			stmt := tQReqAccess.
				UPDATE(
					tQReqAccess.QualificationID,
					tQReqAccess.Access,
				).
				SET(
					qualificationId,
					access.Requirements[k].Access,
				).
				WHERE(
					tQReqAccess.ID.EQ(jet.Uint64(access.Requirements[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) deleteQualificationAccess(ctx context.Context, tx qrm.DB, qualificationId uint64, access *jobs.QualificationAccess) error {
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

		tQJobAccess := table.FivenetJobsQualificationsJobAccess
		jobStmt := tQJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					tQJobAccess.ID.IN(jobIds...),
					tQJobAccess.QualificationID.EQ(jet.Uint64(qualificationId)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if access.Requirements != nil && len(access.Requirements) > 0 {
		uaIds := []jet.Expression{}
		for i := 0; i < len(access.Requirements); i++ {
			if access.Requirements[i].Id == 0 {
				continue
			}
			uaIds = append(uaIds, jet.Uint64(access.Requirements[i].Id))
		}

		tQReqAccess := table.FivenetJobsQualificationsReqsAccess
		userStmt := tQReqAccess.
			DELETE().
			WHERE(
				jet.AND(
					tQReqAccess.ID.IN(uaIds...),
					tQReqAccess.QualificationID.EQ(jet.Uint64(qualificationId)),
				),
			).
			LIMIT(25)

		if _, err := userStmt.ExecContext(ctx, tx); err != nil {
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

	userStmt := tQReqAccess.
		DELETE().
		WHERE(tQReqAccess.QualificationID.EQ(jet.Uint64(qualificationId)))

	if _, err := userStmt.ExecContext(ctx, tx); err != nil {
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
