package qualifications

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	expireExamsCronjob = "qualifications.exams.expire"
	retainExamsCronjob = "qualifications.exams.retain"

	// examMaxRetentionDays bounds stored attempt snapshots, answers, and grading.
	// Qualification results themselves are intentionally retained.
	examMaxRetentionDays   = 60
	examExpiryBatchSize    = 100
	examRetentionBatchSize = 1000
)

type ExamHousekeeper struct {
	logger *zap.Logger
	store  qualificationsstore.IStore
	server *Server
}

type ExamHousekeeperResult struct {
	fx.Out

	Housekeeper  *ExamHousekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewExamHousekeeper(p Params, server *Server) ExamHousekeeperResult {
	h := &ExamHousekeeper{
		logger: p.Logger.Named("qualifications.exam_housekeeper"),
		store:  p.Store,
		server: server,
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
	return registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     retainExamsCronjob,
		Schedule: "15 3 * * *", // Daily at 03:15.
	})
}

func (h *ExamHousekeeper) RegisterCronjobHandlers(handlers *croner.Handlers) error {
	handlers.Add(expireExamsCronjob, func(ctx context.Context, _ *cron.CronjobData) error {
		attempts, err := h.store.ListExpiredExamUsers(ctx, examExpiryBatchSize)
		if err != nil {
			return fmt.Errorf("expire exam attempts: %w", err)
		}
		for _, attempt := range attempts {
			if err := h.completeExpiredExam(ctx, attempt); err != nil {
				h.logger.Error(
					"failed to complete expired exam",
					zap.Int64("qualification_id", attempt.GetQualificationId()),
					zap.Int32("user_id", attempt.GetUserId()),
					zap.Error(err),
				)
			}
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
		for _, attempt := range attempts {
			if err := h.deleteRetainedExamData(ctx, attempt); err != nil {
				h.logger.Error(
					"failed to delete retained exam data",
					zap.Int64("qualification_id", attempt.GetQualificationId()),
					zap.Int32("user_id", attempt.GetUserId()),
					zap.Error(err),
				)
			}
		}
		return nil
	})
	return nil
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
