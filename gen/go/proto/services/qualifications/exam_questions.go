package qualifications

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/filestore"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/storage"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getExamQuestions(ctx context.Context, tx qrm.DB, qualificationId uint64, withAnswers bool) (*qualifications.ExamQuestions, error) {
	columns := []jet.Projection{
		tExamQuestions.QualificationID,
		tExamQuestions.CreatedAt,
		tExamQuestions.UpdatedAt,
		tExamQuestions.Title,
		tExamQuestions.Description,
		tExamQuestions.Data,
		tExamQuestions.Points,
	}

	if withAnswers {
		columns = append(columns, tExamQuestions.Answer)
	}

	stmt := tExamQuestions.
		SELECT(
			tExamQuestions.ID,
			columns...,
		).
		FROM(tExamQuestions).
		WHERE(jet.AND(
			tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
		))

	var dest qualifications.ExamQuestions
	if err := stmt.QueryContext(ctx, tx, &dest.Questions); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) countExamQuestions(ctx context.Context, qualificationid uint64) (int32, error) {
	stmt := tExamQuestions.
		SELECT(
			jet.COUNT(jet.DISTINCT(tExamQuestions.ID)).AS("datacount.totalcount"),
		).
		FROM(tExamQuestions).
		WHERE(
			tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationid)),
		)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		return 0, err
	}

	return int32(count.TotalCount), nil
}

func (s *Server) handleExamQuestionsChanges(ctx context.Context, tx qrm.DB, qualificationId uint64, questions *qualifications.ExamQuestions) error {
	tExamQuestions := table.FivenetQualificationsExamQuestions
	if len(questions.Questions) == 0 {
		stmt := tExamQuestions.
			DELETE().
			WHERE(tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId))).
			LIMIT(100)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}

		return nil
	}
	current, err := s.getExamQuestions(ctx, tx, qualificationId, false)
	if err != nil {
		return err
	}

	toCreate, toUpdate, toDelete := s.compareExamQuestions(current.Questions, questions.Questions)

	for _, question := range toCreate {
		if question.Data == nil {
			continue
		}

		switch data := question.Data.Data.(type) {
		case *qualifications.ExamQuestionData_Image:
			if data.Image.Image == nil || data.Image.Image.Url == nil {
				continue
			}

			if len(data.Image.Image.Data) == 0 {
				continue
			}

			if !data.Image.Image.IsImage() {
				return errorsqualifications.ErrFailedQuery
			}

			if err := data.Image.Image.Optimize(ctx); err != nil {
				return err
			}

			if err := data.Image.Image.Upload(ctx, s.st, filestore.QualificationExamAssets, storage.FileNameSplitter(data.Image.Image.GetHash())); err != nil {
				return err
			}
		}

		stmt := tExamQuestions.
			INSERT(
				tExamQuestions.QualificationID,
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
			).
			VALUES(
				qualificationId,
				question.Title,
				question.Description,
				question.Data,
				question.Answer,
				question.Points,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	for _, question := range toUpdate {
		if question.Data != nil {
			switch data := question.Data.Data.(type) {
			case *qualifications.ExamQuestionData_Image:
				if data.Image.Image == nil {
					continue
				}

				if len(data.Image.Image.Data) == 0 {
					continue
				}

				if !data.Image.Image.IsImage() {
					return errorsqualifications.ErrFailedQuery
				}

				if err := data.Image.Image.Optimize(ctx); err != nil {
					return err
				}

				if err := data.Image.Image.Upload(ctx, s.st, filestore.QualificationExamAssets, storage.FileNameSplitter(data.Image.Image.GetHash())); err != nil {
					return err
				}
			}
		}

		stmt := tExamQuestions.
			UPDATE(
				tExamQuestions.Title,
				tExamQuestions.Description,
				tExamQuestions.Data,
				tExamQuestions.Answer,
				tExamQuestions.Points,
			).
			SET(
				question.Title,
				question.Description,
				question.Data,
				question.Answer,
				question.Points,
			).
			WHERE(jet.AND(
				tExamQuestions.ID.EQ(jet.Uint64(question.Id)),
				tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if len(toDelete) > 0 {
		questionIds := []jet.Expression{}
		for _, question := range toDelete {
			questionIds = append(questionIds, jet.Uint64(question.Id))
		}

		stmt := tExamQuestions.
			DELETE().
			WHERE(jet.AND(
				tExamQuestions.ID.IN(questionIds...),
				tExamQuestions.QualificationID.EQ(jet.Uint64(qualificationId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}

		for _, question := range toDelete {
			if question.Data == nil {
				continue
			}

			switch data := question.Data.Data.(type) {
			case *qualifications.ExamQuestionData_Image:
				if data.Image.Image == nil || data.Image.Image.Url == nil {
					continue
				}

				if err := s.st.Delete(ctx, filestore.StripURLPrefix(*data.Image.Image.Url)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Server) compareExamQuestions(current, in []*qualifications.ExamQuestion) (toCreate []*qualifications.ExamQuestion, toUpdate []*qualifications.ExamQuestion, toDelete []*qualifications.ExamQuestion) {
	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b *qualifications.ExamQuestion) int {
		return int(a.Id - b.Id)
	})

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cj := range current {
			idx := slices.IndexFunc(in, func(a *qualifications.ExamQuestion) bool {
				return cj.Id == a.Id
			})
			// No match in incoming questions, needs to be deleted
			if idx == -1 {
				toDelete = append(toDelete, cj)
				continue
			}

			foundTracker = append(foundTracker, idx)
			toUpdate = append(toUpdate, in[idx])
		}

		for i, eq := range in {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, eq)
			}
		}
	}

	return
}

type examResponses struct {
	ExamResponses *qualifications.ExamResponses `alias:"responses"`
}

func (s *Server) getExamResponses(ctx context.Context, qualificationId uint64, userId int32) (*qualifications.ExamResponses, error) {
	tExamResponses := tExamResponses.AS("examresponses")
	stmt := tExamResponses.
		SELECT(
			tExamResponses.QualificationID,
			tExamResponses.UserID,
			tExamResponses.Responses,
		).
		FROM(tExamResponses).
		WHERE(jet.AND(
			tExamResponses.QualificationID.EQ(jet.Uint64(qualificationId)),
			tExamResponses.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	dest := examResponses{
		ExamResponses: &qualifications.ExamResponses{},
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	dest.ExamResponses.QualificationId = qualificationId
	dest.ExamResponses.UserId = userId

	return dest.ExamResponses, nil
}
