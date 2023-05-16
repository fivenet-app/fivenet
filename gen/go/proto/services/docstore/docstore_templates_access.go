package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
)

func (s *Server) handleTemplateAccessChanges(ctx context.Context, tx *sql.Tx, templateId uint64, access []*documents.TemplateJobAccess) error {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Get existing job and user accesses from database
	current, err := s.getTemplateJobAccess(ctx, templateId)
	if err != nil {
		return err
	}

	toCreate, toUpdate, toDelete := s.compareTemplateJobAccess(tx, current, access)

	if err := s.createTemplateJobAccess(ctx, tx, templateId, userInfo.UserId, toCreate); err != nil {
		return err
	}

	if err := s.updateTemplateJobAccess(ctx, tx, templateId, userInfo.UserId, toUpdate); err != nil {
		return err
	}

	if err := s.deleteTemplateJobAccess(ctx, tx, templateId, toDelete); err != nil {
		return err
	}

	return nil
}

func (s *Server) getTemplateJobAccess(ctx context.Context, templateId uint64) ([]*documents.TemplateJobAccess, error) {
	dTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess.AS("templatejobaccess")
	jobStmt := dTemplatesJobAccess.
		SELECT(
			dTemplatesJobAccess.AllColumns,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			dTemplatesJobAccess.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(dTemplatesJobAccess.CreatorID),
				),
		).
		WHERE(
			dTemplatesJobAccess.TemplateID.EQ(jet.Uint64(templateId)),
		).
		ORDER_BY(
			dTemplatesJobAccess.ID.ASC(),
		)

	var jobAccess []*documents.TemplateJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return jobAccess, nil
}

func (s *Server) compareTemplateJobAccess(tx *sql.Tx, current, in []*documents.TemplateJobAccess) (toCreate []*documents.TemplateJobAccess, toUpdate []*documents.TemplateJobAccess, toDelete []*documents.TemplateJobAccess) {
	toCreate = []*documents.TemplateJobAccess{}
	toUpdate = []*documents.TemplateJobAccess{}
	toDelete = []*documents.TemplateJobAccess{}

	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *documents.TemplateJobAccess) bool {
		return a.Id > b.Id
	})

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cj := range current {
			var found *documents.TemplateJobAccess
			var foundIdx int
			for k, uj := range in {
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
				toDelete = append(toDelete, cj)
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
				toUpdate = append(toUpdate, cj)
			}
		}

		for k, uj := range in {
			idx := slices.Index(foundTracker, k)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	return
}

func (s *Server) createTemplateJobAccess(ctx context.Context, tx *sql.Tx, templateId uint64, userId int32, access []*documents.TemplateJobAccess) error {
	if access == nil {
		return nil
	}

	for k := 0; k < len(access); k++ {
		// Create document job access
		dTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		stmt := dTemplatesJobAccess.
			INSERT(
				dTemplatesJobAccess.TemplateID,
				dTemplatesJobAccess.Job,
				dTemplatesJobAccess.MinimumGrade,
				dTemplatesJobAccess.Access,
				dTemplatesJobAccess.CreatorID,
			).
			VALUES(
				templateId,
				access[k].Job,
				access[k].MinimumGrade,
				access[k].Access,
				userId,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) updateTemplateJobAccess(ctx context.Context, tx *sql.Tx, templateId uint64, userId int32, access []*documents.TemplateJobAccess) error {
	if access == nil {
		return nil
	}

	for k := 0; k < len(access); k++ {
		// Create document job access
		dTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		stmt := dTemplatesJobAccess.
			UPDATE(
				dTemplatesJobAccess.TemplateID,
				dTemplatesJobAccess.Job,
				dTemplatesJobAccess.MinimumGrade,
				dTemplatesJobAccess.Access,
				dTemplatesJobAccess.CreatorID,
			).
			SET(
				templateId,
				access[k].Job,
				access[k].MinimumGrade,
				access[k].Access,
				userId,
			).
			WHERE(
				dTemplatesJobAccess.ID.EQ(jet.Uint64(access[k].Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) deleteTemplateJobAccess(ctx context.Context, tx *sql.Tx, templateId uint64, access []*documents.TemplateJobAccess) error {
	if access == nil {
		return nil
	}

	if len(access) > 0 {
		jobIds := []jet.Expression{}
		for i := 0; i < len(access); i++ {
			if access[i].Id == 0 {
				continue
			}
			jobIds = append(jobIds, jet.Uint64(access[i].Id))
		}

		dTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		jobStmt := dTemplatesJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					dTemplatesJobAccess.ID.IN(jobIds...),
					dTemplatesJobAccess.TemplateID.EQ(jet.Uint64(templateId)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) checkIfUserHasAccessToTemplate(ctx context.Context, templateId uint64, userInfo *userinfo.UserInfo, publicOk bool, access documents.ACCESS_LEVEL) (bool, error) {
	out, err := s.checkIfUserHasAccessToTemplateIDs(ctx, userInfo, publicOk, access, templateId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToTemplateIDs(ctx context.Context, userInfo *userinfo.UserInfo, publicOk bool, access documents.ACCESS_LEVEL, templateIds ...uint64) ([]uint64, error) {
	if len(templateIds) == 0 {
		return templateIds, nil
	}

	// Allow superusers access to any templates
	if userInfo.SuperUser {
		return templateIds, nil
	}

	ids := make([]jet.Expression, len(templateIds))
	for i := 0; i < len(templateIds); i++ {
		ids[i] = jet.Uint64(templateIds[i])
	}

	condition := jet.AND(
		tDTemplates.ID.IN(ids...),
		tDTemplates.DeletedAt.IS_NULL(),
		jet.OR(
			tDTemplates.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				tDTemplatesJobAccess.Access.IS_NOT_NULL(),
				tDTemplatesJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
			),
		),
	)

	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID,
		).
		FROM(
			tDTemplates.
				LEFT_JOIN(tDTemplatesJobAccess,
					tDTemplatesJobAccess.TemplateID.EQ(tDTemplates.ID).
						AND(tDTemplatesJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDTemplatesJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition).
		GROUP_BY(tDTemplates.ID).
		ORDER_BY(tDTemplates.ID.DESC(), tDTemplatesJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest.IDs, nil
}
