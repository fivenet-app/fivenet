package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"golang.org/x/exp/slices"
)

func (s *Server) handleTemplateAccessChanges(ctx context.Context, tx *sql.Tx, templateId uint64, access []*documents.TemplateJobAccess) error {
	userId := auth.GetUserIDFromContext(ctx)

	// Get existing job and user accesses from database
	current, err := s.getTemplateJobAccess(ctx, templateId)
	if err != nil {
		return err
	}

	toCreate, toUpdate, toDelete := s.compareTemplateJobAccess(tx, current, access)

	if err := s.createTemplateJobAccess(ctx, tx, templateId, userId, toCreate); err != nil {
		return err
	}

	if err := s.updateTemplateJobAccess(ctx, tx, templateId, userId, toUpdate); err != nil {
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
			uCreator.ID,
			uCreator.Identifier,
			uCreator.Job,
			uCreator.JobGrade,
			uCreator.Firstname,
			uCreator.Lastname,
		).
		FROM(
			dTemplatesJobAccess.
				LEFT_JOIN(uCreator,
					uCreator.ID.EQ(dTemplatesJobAccess.CreatorID),
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

	if current == nil {
		return
	}

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

func (s *Server) checkIfUserHasAccessToTemplate(ctx context.Context, templateId uint64, userId int32, job string, jobGrade int32, publicOk bool, access documents.ACCESS_LEVEL) (bool, error) {
	out, err := s.checkIfUserHasAccessToTemplateIDs(ctx, userId, job, jobGrade, publicOk, access, templateId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToTemplateIDs(ctx context.Context, userId int32, job string, jobGrade int32, publicOk bool, access documents.ACCESS_LEVEL, templateIds ...uint64) ([]uint64, error) {
	if len(templateIds) == 0 {
		return templateIds, nil
	}

	// Allow superusers access to any templates
	if s.p.Can(userId, common.SuperuserAnyAccess) {
		return templateIds, nil
	}

	ids := make([]jet.Expression, len(templateIds))
	for i := 0; i < len(templateIds); i++ {
		ids[i] = jet.Uint64(templateIds[i])
	}

	condition := jet.AND(
		dTemplates.ID.IN(ids...),
		dTemplates.DeletedAt.IS_NULL(),
		jet.OR(
			dTemplates.CreatorID.EQ(jet.Int32(userId)),
			jet.AND(
				dTemplatesJobAccess.Access.IS_NOT_NULL(),
				dTemplatesJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
			),
		),
	)

	stmt := dTemplates.
		SELECT(
			dTemplates.ID,
		).
		FROM(
			dTemplates.
				LEFT_JOIN(dTemplatesJobAccess,
					dTemplatesJobAccess.TemplateID.EQ(dTemplates.ID).
						AND(dTemplatesJobAccess.Job.EQ(jet.String(job))).
						AND(dTemplatesJobAccess.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				),
		).
		WHERE(condition).
		GROUP_BY(dTemplates.ID).
		ORDER_BY(dTemplates.ID.DESC(), dTemplatesJobAccess.MinimumGrade)

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
