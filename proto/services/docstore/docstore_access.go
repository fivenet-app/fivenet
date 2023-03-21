package docstore

import (
	context "context"
	"errors"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_ACCESS)
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
		s.c.EnrichJobInfo(access.Users[i].User)
	}

	resp := &GetDocumentAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *SetDocumentAccessRequest) (*SetDocumentAccessResponse, error) {
	resp := &SetDocumentAccessResponse{}

	if err := s.handleDocumentAccessChanges(ctx, req.Mode, req.DocumentId, req.Access); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) handleDocumentAccessChanges(ctx context.Context, mode DOC_ACCESS_UPDATE_MODE, documentID uint64, access *documents.DocumentAccess) error {
	userId := auth.GetUserIDFromContext(ctx)

	// Get existing job and user accesses from database
	current, err := s.getDocumentAccess(ctx, documentID)
	if err != nil {
		return err
	}

	switch mode {
	case DOC_ACCESS_UPDATE_MODE_UPDATE:
		toCreate, toUpdate, toDelete := s.compareDocumentAccess(current, access)

		if err := s.createDocumentAccess(ctx, documentID, userId, toCreate); err != nil {
			return err
		}

		if err := s.updateDocumentAccess(ctx, documentID, userId, toUpdate); err != nil {
			return err
		}

		if err := s.deleteDocumentAccess(ctx, documentID, toDelete); err != nil {
			return err
		}

	case DOC_ACCESS_UPDATE_MODE_DELETE:
		if err := s.deleteDocumentAccess(ctx, documentID, access); err != nil {
			return err
		}

	case DOC_ACCESS_UPDATE_MODE_CLEAR:
		if err := s.clearDocumentAccess(ctx, documentID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareDocumentAccess(current, in *documents.DocumentAccess) (*documents.DocumentAccess, *documents.DocumentAccess, *documents.DocumentAccess) {
	toCreate := &documents.DocumentAccess{}
	toUpdate := &documents.DocumentAccess{}
	toDelete := &documents.DocumentAccess{}

	if len(current.Jobs) == 0 && len(current.Users) == 0 {
		return in, nil, nil
	}

	slices.SortFunc(current.Jobs, func(a, b *documents.DocumentJobAccess) bool {
		return a.Id > b.Id
	})
	// TODO compare existing access rules with new access rules

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

	return toCreate, toUpdate, toDelete
}

func (s *Server) getDocumentAccess(ctx context.Context, documentID uint64) (*documents.DocumentAccess, error) {
	dJobAccess := table.ArpanetDocumentsJobAccess.AS("documentjobaccess")
	jobStmt := dJobAccess.
		SELECT(
			dJobAccess.AllColumns,
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
		).
		FROM(
			dJobAccess.
				LEFT_JOIN(uCreator,
					uCreator.ID.EQ(dJobAccess.CreatorID),
				),
		).
		WHERE(
			dJobAccess.DocumentID.EQ(jet.Uint64(documentID)),
		).
		ORDER_BY(
			dJobAccess.ID.ASC(),
		)

	var jobAccess []*documents.DocumentJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	user := user.AS("user")
	dUserAccess := table.ArpanetDocumentsUserAccess.AS("documentuseraccess")
	userStmt := dUserAccess.
		SELECT(
			dUserAccess.AllColumns,
			user.ID,
			user.Identifier,
			user.Job,
			user.JobGrade,
			user.Firstname,
			user.Lastname,
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
		).
		FROM(
			dUserAccess.
				LEFT_JOIN(user,
					user.ID.EQ(dUserAccess.UserID),
				).
				LEFT_JOIN(uCreator,
					uCreator.ID.EQ(dUserAccess.CreatorID),
				),
		).
		WHERE(
			dUserAccess.DocumentID.EQ(jet.Uint64(documentID)),
		).
		ORDER_BY(
			dUserAccess.ID.ASC(),
		)

	var userAccess []*documents.DocumentUserAccess
	if err := userStmt.QueryContext(ctx, s.db, &userAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	resp := &documents.DocumentAccess{
		Jobs:  jobAccess,
		Users: userAccess,
	}

	return resp, nil
}

func (s *Server) createDocumentAccess(ctx context.Context, documentID uint64, userId int32, access *documents.DocumentAccess) error {
	for k := 0; k < len(access.Jobs); k++ {
		// Create document job access
		dJobAccess := table.ArpanetDocumentsJobAccess
		stmt := dJobAccess.
			INSERT(
				dJobAccess.DocumentID,
				dJobAccess.Job,
				dJobAccess.MinimumGrade,
				dJobAccess.Access,
				dJobAccess.CreatorID,
			).
			VALUES(
				documentID,
				access.Jobs[k].Job,
				access.Jobs[k].MinimumGrade,
				access.Jobs[k].Access,
				userId,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	for k := 0; k < len(access.Users); k++ {
		// Create document user access
		dUserAccess := table.ArpanetDocumentsUserAccess
		stmt := dUserAccess.
			INSERT(
				dUserAccess.DocumentID,
				dUserAccess.UserID,
				dUserAccess.Access,
				dUserAccess.CreatorID,
			).
			VALUES(
				documentID,
				access.Users[k].UserId,
				access.Users[k].Access,
				userId,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) updateDocumentAccess(ctx context.Context, documentID uint64, userId int32, access *documents.DocumentAccess) error {
	for k := 0; k < len(access.Jobs); k++ {
		// Create document job access
		dJobAccess := table.ArpanetDocumentsJobAccess
		stmt := dJobAccess.
			UPDATE(
				dJobAccess.DocumentID,
				dJobAccess.Job,
				dJobAccess.MinimumGrade,
				dJobAccess.Access,
				dJobAccess.CreatorID,
			).
			SET(
				documentID,
				access.Jobs[k].Job,
				access.Jobs[k].MinimumGrade,
				access.Jobs[k].Access,
				userId,
			).
			WHERE(
				dJobAccess.ID.EQ(jet.Uint64(access.Jobs[k].Id)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	for k := 0; k < len(access.Users); k++ {
		// Create document user access
		dUserAccess := table.ArpanetDocumentsUserAccess
		stmt := dUserAccess.
			UPDATE(
				dUserAccess.DocumentID,
				dUserAccess.UserID,
				dUserAccess.Access,
				dUserAccess.CreatorID,
			).
			SET(
				documentID,
				access.Users[k].UserId,
				access.Users[k].Access,
				userId,
			).
			WHERE(
				dUserAccess.ID.EQ(jet.Uint64(access.Users[k].Id)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) deleteDocumentAccess(ctx context.Context, documentID uint64, access *documents.DocumentAccess) error {
	if len(access.Jobs) > 0 {
		jobIds := []jet.Expression{}
		for i := 0; i < len(access.Jobs); i++ {
			if access.Jobs[i].Id == 0 {
				continue
			}
			jobIds = append(jobIds, jet.Uint64(access.Jobs[i].Id))
		}

		dJobAccess := table.ArpanetDocumentsJobAccess
		jobStmt := dJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					dJobAccess.ID.IN(jobIds...),
					dJobAccess.DocumentID.EQ(jet.Uint64(documentID)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	if len(access.Users) > 0 {
		uaIds := []jet.Expression{}
		for i := 0; i < len(access.Users); i++ {
			if access.Users[i].Id == 0 {
				continue
			}
			uaIds = append(uaIds, jet.Uint64(access.Users[i].Id))
		}

		dUserAccess := table.ArpanetDocumentsUserAccess
		userStmt := dUserAccess.
			DELETE().
			WHERE(
				jet.AND(
					dUserAccess.ID.IN(uaIds...),
					dUserAccess.DocumentID.EQ(jet.Uint64(documentID)),
				),
			).
			LIMIT(25)

		if _, err := userStmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearDocumentAccess(ctx context.Context, documentID uint64) error {
	jobStmt := dJobAccess.
		DELETE().
		WHERE(dJobAccess.DocumentID.EQ(jet.Uint64(documentID)))

	if _, err := jobStmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	userStmt := dUserAccess.
		DELETE().
		WHERE(dUserAccess.DocumentID.EQ(jet.Uint64(documentID)))

	if _, err := userStmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
