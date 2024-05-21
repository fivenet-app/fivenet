package docstore

import (
	context "context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getTemplateJobAccess(ctx context.Context, templateId uint64) ([]*documents.TemplateJobAccess, error) {
	tDTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess.AS("templatejobaccess")
	jobStmt := tDTemplatesJobAccess.
		SELECT(
			tDTemplatesJobAccess.ID,
			tDTemplatesJobAccess.CreatedAt,
			tDTemplatesJobAccess.TemplateID,
			tDTemplatesJobAccess.Job,
			tDTemplatesJobAccess.MinimumGrade,
			tDTemplatesJobAccess.Access,
		).
		FROM(
			tDTemplatesJobAccess,
		).
		WHERE(
			tDTemplatesJobAccess.TemplateID.EQ(jet.Uint64(templateId)),
		).
		ORDER_BY(
			tDTemplatesJobAccess.ID.ASC(),
		)

	var jobAccess []*documents.TemplateJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return jobAccess, nil
}

func (s *Server) handleTemplateAccessChanges(ctx context.Context, tx qrm.DB, templateId uint64, job string, access []*documents.TemplateJobAccess) error {
	// Get existing job and user accesses from database
	current, err := s.getTemplateJobAccess(ctx, templateId)
	if err != nil {
		return err
	}

	// Make sure the template has at least someone able to edit the template, if not add the job's highest grade with edit access
	if len(access) == 0 || !slices.ContainsFunc(access, func(in *documents.TemplateJobAccess) bool {
		return in.Access == documents.AccessLevel_ACCESS_LEVEL_EDIT
	}) {
		job := s.enricher.GetJobByName(job)
		if job == nil {
			return fmt.Errorf("failed to get user job for template access")
		}

		if len(job.Grades) == 0 || job.Grades[len(job.Grades)-1] == nil {
			return fmt.Errorf("failed to get highest grade for job for template access")
		}

		grade := job.Grades[len(job.Grades)-1]

		access = append(access, &documents.TemplateJobAccess{
			TemplateId:   templateId,
			Job:          job.Name,
			MinimumGrade: grade.Grade,
			Access:       documents.AccessLevel_ACCESS_LEVEL_EDIT,
		})
	}

	toCreate, toUpdate, toDelete := s.compareTemplateJobAccess(current, access)

	if err := s.createTemplateJobAccess(ctx, tx, templateId, toCreate); err != nil {
		return err
	}

	if err := s.updateTemplateJobAccess(ctx, tx, templateId, toUpdate); err != nil {
		return err
	}

	if err := s.deleteTemplateJobAccess(ctx, tx, templateId, toDelete); err != nil {
		return err
	}

	return nil
}

func (s *Server) compareTemplateJobAccess(current, in []*documents.TemplateJobAccess) (toCreate []*documents.TemplateJobAccess, toUpdate []*documents.TemplateJobAccess, toDelete []*documents.TemplateJobAccess) {
	toCreate = []*documents.TemplateJobAccess{}
	toUpdate = []*documents.TemplateJobAccess{}
	toDelete = []*documents.TemplateJobAccess{}

	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *documents.TemplateJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cj := range current {
			var found *documents.TemplateJobAccess
			var foundIdx int
			for i, uj := range in {
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

		for i, uj := range in {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	return
}

func (s *Server) createTemplateJobAccess(ctx context.Context, tx qrm.DB, templateId uint64, access []*documents.TemplateJobAccess) error {
	if access == nil {
		return nil
	}

	for k := 0; k < len(access); k++ {
		// Create document job access
		tDTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		stmt := tDTemplatesJobAccess.
			INSERT(
				tDTemplatesJobAccess.TemplateID,
				tDTemplatesJobAccess.Job,
				tDTemplatesJobAccess.MinimumGrade,
				tDTemplatesJobAccess.Access,
			).
			VALUES(
				templateId,
				access[k].Job,
				access[k].MinimumGrade,
				access[k].Access,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) updateTemplateJobAccess(ctx context.Context, tx qrm.DB, templateId uint64, access []*documents.TemplateJobAccess) error {
	if access == nil {
		return nil
	}

	for k := 0; k < len(access); k++ {
		// Create document job access
		tDTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		stmt := tDTemplatesJobAccess.
			UPDATE(
				tDTemplatesJobAccess.TemplateID,
				tDTemplatesJobAccess.Job,
				tDTemplatesJobAccess.MinimumGrade,
				tDTemplatesJobAccess.Access,
			).
			SET(
				templateId,
				access[k].Job,
				access[k].MinimumGrade,
				access[k].Access,
			).
			WHERE(
				tDTemplatesJobAccess.ID.EQ(jet.Uint64(access[k].Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) deleteTemplateJobAccess(ctx context.Context, tx qrm.DB, templateId uint64, access []*documents.TemplateJobAccess) error {
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

		tDTemplatesJobAccess := table.FivenetDocumentsTemplatesJobAccess
		jobStmt := tDTemplatesJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					tDTemplatesJobAccess.ID.IN(jobIds...),
					tDTemplatesJobAccess.TemplateID.EQ(jet.Uint64(templateId)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) checkIfUserHasAccessToTemplate(ctx context.Context, templateId uint64, userInfo *userinfo.UserInfo, access documents.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToTemplateIDs(ctx, userInfo, access, templateId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToTemplateIDs(ctx context.Context, userInfo *userinfo.UserInfo, access documents.AccessLevel, templateIds ...uint64) ([]uint64, error) {
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
		jet.AND(
			tDTemplatesJobAccess.Access.IS_NOT_NULL(),
			tDTemplatesJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
		),
	)

	stmt := tDTemplates.
		SELECT(
			tDTemplates.ID.AS("id"),
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

	var dest []uint64
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
