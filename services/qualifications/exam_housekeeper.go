package qualifications

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	expireExamsCronjob              = "qualifications.exams.expire"
	retainExamsCronjob              = "qualifications.exams.retain"
	cleanupExamQuestionFilesCronjob = "qualifications.exams.cleanup_files"

	// examMaxRetentionDays bounds stored attempt snapshots, answers, and grading.
	// Qualification results themselves are intentionally retained.
	examMaxRetentionDays                            = 60
	examExpiryBatchSize                             = 100
	examRetentionBatchSize                          = 1000
	examQuestionFileCleanupBatchSize                = 100
	examQuestionFileCleanupAge                      = 24 * time.Hour
	cleanupExamQuestionFilesLastQualificationIDAttr = "last_qualification_id"
)

type ExamHousekeeper struct {
	logger *zap.Logger
	store  qualificationsstore.IStore
	server *Server
}

type ExamHousekeeperParams struct {
	fx.In

	Logger *zap.Logger
	Store  qualificationsstore.IStore
	Server *Server
}

type ExamHousekeeperResult struct {
	fx.Out

	Housekeeper  *ExamHousekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewExamHousekeeper(p ExamHousekeeperParams) ExamHousekeeperResult {
	h := &ExamHousekeeper{
		logger: p.Logger.Named("qualifications.exam_housekeeper"),
		store:  p.Store,
		server: p.Server,
	}
	return ExamHousekeeperResult{Housekeeper: h, CronRegister: h}
}

func (h *ExamHousekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     expireExamsCronjob,
		Schedule: "* * * * *", // Every minute.
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     retainExamsCronjob,
		Schedule: "15 3 * * *", // Daily at 03:15.
	}); err != nil {
		return err
	}
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     cleanupExamQuestionFilesCronjob,
		Schedule: "45 * * * *", // Hourly.
	}); err != nil {
		return err
	}
	return nil
}

func (h *ExamHousekeeper) RegisterCronjobHandlers(handlers *croner.Handlers) error {
	handlers.Add(expireExamsCronjob, func(ctx context.Context, _ *cron.CronjobData) error {
		attempts, err := h.store.ListExpiredExamUsers(ctx, examExpiryBatchSize)
		if err != nil {
			return fmt.Errorf("expire exam attempts: %w", err)
		}
		var failed error
		failedCount := 0
		for _, attempt := range attempts {
			if err := h.completeExpiredExam(ctx, attempt); err != nil {
				failedCount++
				failed = errors.Join(failed, err)
				h.logger.Error(
					"failed to complete expired exam",
					zap.Int64("qualification_id", attempt.GetQualificationId()),
					zap.Int32("user_id", attempt.GetUserId()),
					zap.Error(err),
				)
			}
		}
		if failed != nil {
			return fmt.Errorf(
				"failed to complete %d expired exam attempts: %w",
				failedCount,
				failed,
			)
		}
		return nil
	})
	handlers.Add(retainExamsCronjob, func(ctx context.Context, _ *cron.CronjobData) error {
		attempts, err := h.store.ListExamUsersPastRetention(
			ctx,
			time.Now().AddDate(0, 0, -examMaxRetentionDays),
			examRetentionBatchSize,
		)
		if err != nil {
			return fmt.Errorf("list expired exam retention data: %w", err)
		}
		var failed error
		failedCount := 0
		for _, attempt := range attempts {
			if err := h.deleteRetainedExamData(ctx, attempt); err != nil {
				failedCount++
				failed = errors.Join(failed, err)
				h.logger.Error(
					"failed to delete retained exam data",
					zap.Int64("qualification_id", attempt.GetQualificationId()),
					zap.Int32("user_id", attempt.GetUserId()),
					zap.Error(err),
				)
			}
		}
		if failed != nil {
			return fmt.Errorf("failed to delete %d retained exam attempts: %w", failedCount, failed)
		}
		return nil
	})
	handlers.Add(
		cleanupExamQuestionFilesCronjob,
		func(ctx context.Context, data *cron.CronjobData) error {
			dest := &cron.GenericCronData{Attributes: map[string]string{}}
			if err := data.Unmarshal(dest); err != nil {
				h.logger.Warn(
					"failed to unmarshal exam question file cleanup cron data",
					zap.Error(err),
				)
			}

			lastQualificationID := int64(0)
			if raw := dest.GetAttribute(
				cleanupExamQuestionFilesLastQualificationIDAttr,
			); raw != "" {
				if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
					lastQualificationID = parsed
				}
			}

			ids, err := h.store.ListActiveQualificationIDs(
				ctx,
				lastQualificationID,
				examQuestionFileCleanupBatchSize,
			)
			if err != nil {
				return fmt.Errorf(
					"list active qualifications for exam question file cleanup: %w",
					err,
				)
			}
			if len(ids) == 0 && lastQualificationID > 0 {
				lastQualificationID = 0
				ids, err = h.store.ListActiveQualificationIDs(
					ctx,
					0,
					examQuestionFileCleanupBatchSize,
				)
				if err != nil {
					return fmt.Errorf(
						"wrap active qualifications for exam question file cleanup: %w",
						err,
					)
				}
			}

			cutoff := time.Now().Add(-examQuestionFileCleanupAge)
			for _, qualificationID := range ids {
				if err := h.cleanupExamQuestionFiles(ctx, qualificationID, cutoff); err != nil {
					return fmt.Errorf(
						"clean up exam question files for qualification %d: %w",
						qualificationID,
						err,
					)
				}
				lastQualificationID = qualificationID
			}

			if len(ids) > 0 {
				dest.SetAttribute(
					cleanupExamQuestionFilesLastQualificationIDAttr,
					strconv.FormatInt(lastQualificationID, 10),
				)
			} else {
				dest.SetAttribute(cleanupExamQuestionFilesLastQualificationIDAttr, "0")
			}
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("marshal exam question file cleanup cron data: %w", err)
			}
			return nil
		},
	)
	return nil
}

func (h *ExamHousekeeper) cleanupExamQuestionFiles(
	ctx context.Context,
	qualificationID int64,
	cutoff time.Time,
) error {
	tx, err := h.server.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	exam, err := h.store.GetExamQuestions(ctx, tx, qualificationID, false)
	if err != nil {
		return err
	}

	referencedFileIDs := make([]int64, 0, len(exam.GetQuestions()))
	for _, question := range exam.GetQuestions() {
		if image := question.GetData().GetImage(); image != nil && image.GetImage() != nil {
			referencedFileIDs = append(referencedFileIDs, image.GetImage().GetId())
		}
	}

	_, err = h.store.UnlinkStaleExamQuestionFiles(
		ctx,
		tx,
		qualificationID,
		referencedFileIDs,
		cutoff,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (h *ExamHousekeeper) deleteRetainedExamData(
	ctx context.Context,
	attempt *qualificationsexam.ExamUser,
) error {
	tx, err := h.server.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := h.store.DeleteExamResponses(
		ctx,
		tx,
		attempt.GetAttemptId(),
	); err != nil {
		return err
	}
	if err := h.store.DeleteExamUser(
		ctx,
		tx,
		attempt.GetAttemptId(),
	); err != nil {
		return err
	}
	return tx.Commit()
}

func (h *ExamHousekeeper) completeExpiredExam(
	ctx context.Context,
	attempt *qualificationsexam.ExamUser,
) error {
	quali, err := h.store.GetQualification(
		ctx,
		attempt.GetQualificationId(),
		&userinfo.UserInfo{Superuser: true},
		false,
	)
	if err != nil {
		return err
	}
	tx, err := h.server.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	expired, err := h.store.ExpireExamUser(
		ctx,
		tx,
		attempt.GetQualificationId(),
		attempt.GetUserId(),
		attempt.GetAttemptId(),
	)
	if err != nil || !expired {
		return err
	}
	if err := h.server.addQualificationActivityForAttempt(
		ctx,
		tx,
		attempt.GetQualificationId(),
		qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_EXAM_EXPIRED,
		0,
		attempt.GetUserId(),
		attempt.GetAttemptId(),
		nil,
	); err != nil {
		return err
	}
	// Expiry holds the attempt row lock, so no partial submission can update the
	// responses between this read and the grading transaction.
	responses, _, err := h.store.GetExamResponses(
		ctx,
		tx,
		attempt.GetAttemptId(),
	)
	if err != nil {
		return err
	}
	publishNotifications := make([]func(context.Context) error, 0, 1)
	if err := h.server.gradeExam(
		ctx,
		tx,
		attempt.GetQualificationId(),
		attempt.GetUserId(),
		quali,
		attempt.GetSnapshot(),
		responses,
		&publishNotifications,
	); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	notifi.PublishAfterCommit(
		ctx,
		h.logger,
		"qualification_result_updated",
		publishNotifications...)
	// A per-user "your exam has expired" notification can be prepared here.
	return nil
}
