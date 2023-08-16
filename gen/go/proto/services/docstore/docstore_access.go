package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document access!")
	}

	access, err := s.getDocumentAccess(ctx, req.DocumentId)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(access.Jobs); i++ {
		s.c.EnrichJobInfo(access.Jobs[i])
	}

	for i := 0; i < len(access.Users); i++ {
		if access.Users[i].User != nil {
			s.c.EnrichJobInfo(access.Users[i].User)
		}
	}

	resp := &GetDocumentAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *SetDocumentAccessRequest) (*SetDocumentAccessResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "SetDocumentAccess",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.ACCESS_LEVEL_ACCESS)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit the document access!")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.handleDocumentAccessChanges(ctx, tx, req.Mode, req.DocumentId, req.Access); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &SetDocumentAccessResponse{}, nil
}

func (s *Server) handleDocumentAccessChanges(ctx context.Context, tx *sql.Tx, mode ACCESS_LEVEL_UPDATE_MODE, documentId uint64, access *documents.DocumentAccess) error {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Get existing job and user accesses from database
	current, err := s.getDocumentAccess(ctx, documentId)
	if err != nil {
		return err
	}

	switch mode {
	case ACCESS_LEVEL_UPDATE_MODE_UPDATE:
		toCreate, toUpdate, toDelete := s.compareDocumentAccess(tx, current, access)

		if err := s.createDocumentAccess(ctx, tx, documentId, userInfo.UserId, toCreate); err != nil {
			return err
		}

		if err := s.updateDocumentAccess(ctx, tx, documentId, userInfo.UserId, toUpdate); err != nil {
			return err
		}

		if err := s.deleteDocumentAccess(ctx, tx, documentId, toDelete); err != nil {
			return err
		}

	case ACCESS_LEVEL_UPDATE_MODE_DELETE:
		if err := s.deleteDocumentAccess(ctx, tx, documentId, access); err != nil {
			return err
		}

	case ACCESS_LEVEL_UPDATE_MODE_CLEAR:
		if err := s.clearDocumentAccess(ctx, tx, documentId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareDocumentAccess(tx *sql.Tx, current, in *documents.DocumentAccess) (toCreate *documents.DocumentAccess, toUpdate *documents.DocumentAccess, toDelete *documents.DocumentAccess) {
	toCreate = &documents.DocumentAccess{}
	toUpdate = &documents.DocumentAccess{}
	toDelete = &documents.DocumentAccess{}

	if current == nil || (len(current.Jobs) == 0 && len(current.Users) == 0) {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current.Jobs, func(a, b *documents.DocumentJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current.Jobs) == 0 {
		toCreate.Jobs = in.Jobs
	} else {
		foundTracker := []int{}
		for _, cj := range current.Jobs {
			var found *documents.DocumentJobAccess
			var foundIdx int
			for k, uj := range in.Jobs {
				if cj.Job != uj.Job {
					continue
				}
				if cj.MinimumGrade != uj.MinimumGrade {
					continue
				}
				found = uj
				foundIdx = k
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

		for k, uj := range in.Jobs {
			idx := slices.Index(foundTracker, k)
			if idx == -1 {
				toCreate.Jobs = append(toCreate.Jobs, uj)
			}
		}
	}

	if len(current.Users) == 0 {
		toCreate.Users = in.Users
	} else {
		foundTracker := []int{}
		for _, cj := range current.Users {
			var found *documents.DocumentUserAccess
			var foundIdx int
			for k, uj := range in.Users {
				if cj.UserId != uj.UserId {
					continue
				}
				found = uj
				foundIdx = k
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete.Users = append(toDelete.Users, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.Access != found.Access {
				cj.Access = found.Access
				changed = true
			}

			if changed {
				toUpdate.Users = append(toUpdate.Users, cj)
			}
		}

		for k, uj := range in.Users {
			idx := slices.Index(foundTracker, k)
			if idx == -1 {
				toCreate.Users = append(toCreate.Users, uj)
			}
		}
	}

	return
}

func (s *Server) getDocumentAccess(ctx context.Context, documentId uint64) (*documents.DocumentAccess, error) {
	tDJobAccess := table.FivenetDocumentsJobAccess.AS("documentjobaccess")
	jobStmt := tDJobAccess.
		SELECT(
			tDJobAccess.ID,
			tDJobAccess.DocumentID,
			tDJobAccess.Job,
			tDJobAccess.MinimumGrade,
			tDJobAccess.Access,
		).
		FROM(
			tDJobAccess,
		).
		WHERE(
			tDJobAccess.DocumentID.EQ(jet.Uint64(documentId)),
		).
		ORDER_BY(
			tDJobAccess.ID.ASC(),
		)

	var jobAccess []*documents.DocumentJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	tUsers := tUsers.AS("usershort")
	tDUserAccess := table.FivenetDocumentsUserAccess.AS("documentuseraccess")
	userStmt := tDUserAccess.
		SELECT(
			tDUserAccess.ID,
			tDUserAccess.DocumentID,
			tDUserAccess.UserID,
			tDUserAccess.Access,
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
		).
		FROM(
			tDUserAccess.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDUserAccess.UserID),
				),
		).
		WHERE(
			tDUserAccess.DocumentID.EQ(jet.Uint64(documentId)),
		).
		ORDER_BY(
			tDUserAccess.ID.ASC(),
		)

	var userAccess []*documents.DocumentUserAccess
	if err := userStmt.QueryContext(ctx, s.db, &userAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &documents.DocumentAccess{
		Jobs:  jobAccess,
		Users: userAccess,
	}, nil
}

func (s *Server) createDocumentAccess(ctx context.Context, tx *sql.Tx, documentId uint64, userId int32, access *documents.DocumentAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tDJobAccess := table.FivenetDocumentsJobAccess
			stmt := tDJobAccess.
				INSERT(
					tDJobAccess.DocumentID,
					tDJobAccess.Job,
					tDJobAccess.MinimumGrade,
					tDJobAccess.Access,
				).
				VALUES(
					documentId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	if access.Users != nil {
		for k := 0; k < len(access.Users); k++ {
			// Create document user access
			tDUserAccess := table.FivenetDocumentsUserAccess
			stmt := tDUserAccess.
				INSERT(
					tDUserAccess.DocumentID,
					tDUserAccess.UserID,
					tDUserAccess.Access,
				).
				VALUES(
					documentId,
					access.Users[k].UserId,
					access.Users[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) updateDocumentAccess(ctx context.Context, tx *sql.Tx, documentId uint64, userId int32, access *documents.DocumentAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tDJobAccess := table.FivenetDocumentsJobAccess
			stmt := tDJobAccess.
				UPDATE(
					tDJobAccess.DocumentID,
					tDJobAccess.Job,
					tDJobAccess.MinimumGrade,
					tDJobAccess.Access,
				).
				SET(
					documentId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				).
				WHERE(
					tDJobAccess.ID.EQ(jet.Uint64(access.Jobs[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	if access.Users != nil {
		for k := 0; k < len(access.Users); k++ {
			// Create document user access
			tDUserAccess := table.FivenetDocumentsUserAccess
			stmt := tDUserAccess.
				UPDATE(
					tDUserAccess.DocumentID,
					tDUserAccess.UserID,
					tDUserAccess.Access,
				).
				SET(
					documentId,
					access.Users[k].UserId,
					access.Users[k].Access,
				).
				WHERE(
					tDUserAccess.ID.EQ(jet.Uint64(access.Users[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) deleteDocumentAccess(ctx context.Context, tx *sql.Tx, documentId uint64, access *documents.DocumentAccess) error {
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

		tDJobAccess := table.FivenetDocumentsJobAccess
		jobStmt := tDJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					tDJobAccess.ID.IN(jobIds...),
					tDJobAccess.DocumentID.EQ(jet.Uint64(documentId)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if access.Users != nil && len(access.Users) > 0 {
		uaIds := []jet.Expression{}
		for i := 0; i < len(access.Users); i++ {
			if access.Users[i].Id == 0 {
				continue
			}
			uaIds = append(uaIds, jet.Uint64(access.Users[i].Id))
		}

		tDUserAccess := table.FivenetDocumentsUserAccess
		userStmt := tDUserAccess.
			DELETE().
			WHERE(
				jet.AND(
					tDUserAccess.ID.IN(uaIds...),
					tDUserAccess.DocumentID.EQ(jet.Uint64(documentId)),
				),
			).
			LIMIT(25)

		if _, err := userStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearDocumentAccess(ctx context.Context, tx *sql.Tx, documentId uint64) error {
	jobStmt := tDJobAccess.
		DELETE().
		WHERE(tDJobAccess.DocumentID.EQ(jet.Uint64(documentId)))

	if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	userStmt := tDUserAccess.
		DELETE().
		WHERE(tDUserAccess.DocumentID.EQ(jet.Uint64(documentId)))

	if _, err := userStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
