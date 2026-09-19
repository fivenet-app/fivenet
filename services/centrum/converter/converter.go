package centrumconverter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

const maxDispatchConvertCount = 15

const (
	converterCronName       = "centrum.dispatch.converter"
	converterLastIDAttr     = "last_id"
	converterSuccessfulAttr = "successful"
	converterFailedAttr     = "failed"
)

type runStats struct {
	lastID     int32
	successful int
	failed     int
}

type metrics struct {
	processed *prometheus.CounterVec
}

var (
	converterMetricsOnce sync.Once
	converterMetricsInst *metrics
)

func getConverterMetrics() *metrics {
	converterMetricsOnce.Do(func() {
		converterMetricsInst = &metrics{
			processed: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "centrum",
				Name:      "dispatch_converter_processed_total",
				Help:      "Phone dispatches processed by the Centrum converter.",
			}, []string{"converter", "outcome"}),
		}
		prometheus.MustRegister(converterMetricsInst.processed)
	})
	return converterMetricsInst
}

type Converter struct {
	logger *zap.Logger
	db     *sql.DB

	metrics *metrics

	dispatchCreateFn func(ctx context.Context, dsp *centrumdispatches.Dispatch) (*centrumdispatches.Dispatch, error)

	enabled       bool
	converterType string
	convertJobs   []string
}

type Params struct {
	fx.In

	Logger *zap.Logger
	DB     *sql.DB
	Config *config.Config

	Dispatches *dispatches.DispatchDB
}

type Result struct {
	fx.Out

	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func New(p Params) Result {
	convertJobs := make([]string, 0, len(p.Config.DispatchCenter.ConvertJobs))
	for _, job := range p.Config.DispatchCenter.ConvertJobs {
		job = strings.TrimSpace(job)
		if job != "" {
			convertJobs = append(convertJobs, job)
		}
	}

	if p.Config.DispatchCenter.Enabled && len(convertJobs) == 0 {
		p.Logger.Warn(
			"dispatch center converter is enabled but no valid convert jobs are configured",
		)
	}

	c := &Converter{
		logger: p.Logger.Named("centrum.converter"),
		db:     p.DB,

		metrics: getConverterMetrics(),

		dispatchCreateFn: p.Dispatches.Create,

		enabled:       p.Config.DispatchCenter.Enabled,
		converterType: p.Config.DispatchCenter.Type,
		convertJobs:   convertJobs,
	}

	return Result{
		CronRegister: c,
	}
}

func (s *Converter) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	return registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     converterCronName,
		Schedule: "*/5 * * * * *",
		Timeout:  durationpb.New(3 * time.Second),
	})
}

func (s *Converter) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add(converterCronName, func(ctx context.Context, data *cron.CronjobData) error {
		if !s.enabled {
			s.logger.Debug("dispatch converter is disabled, skipping cron job")
			return nil
		}
		if len(s.convertJobs) == 0 {
			s.logger.Debug("no convert jobs configured, skipping dispatch converter cron job")
			return nil
		}

		stored := &cron.GenericCronData{Attributes: map[string]string{}}
		if err := data.Unmarshal(stored); err != nil {
			s.logger.Warn("failed to unmarshal dispatch converter cron data", zap.Error(err))
		}

		lastID, _ := strconv.ParseInt(stored.GetAttribute(converterLastIDAttr), 10, 32)
		var stats runStats
		var err error
		switch s.converterType {
		case "lbphone":
			stats, err = s.convertLBPhoneJobMsgToDispatchWithCursor(ctx, int32(lastID))
		case "gksphone":
			stats, err = s.convertGKSPhoneJobMsgToDispatchWithCursor(ctx, int32(lastID))
		default:
			return fmt.Errorf("unknown phone dispatch converter type %q", s.converterType)
		}
		if err != nil {
			return err
		}

		stored.SetAttribute(converterLastIDAttr, strconv.FormatInt(int64(stats.lastID), 10))
		stored.SetAttribute(converterSuccessfulAttr, strconv.Itoa(stats.successful))
		stored.SetAttribute(converterFailedAttr, strconv.Itoa(stats.failed))
		if err := data.MarshalFrom(stored); err != nil {
			return fmt.Errorf("failed to marshal dispatch converter cron data. %w", err)
		}
		return nil
	})
	return nil
}

func (s *Converter) convertGKSPhoneJobMsgToDispatchWithCursor(
	ctx context.Context,
	lastID int32,
) (runStats, error) {
	stats := runStats{lastID: lastID}
	tGksPhoneJMsg := table.GksphoneJobMessage
	tGksPhoneSettings := table.GksphoneSettings
	tUsers := table.FivenetUser

	jobs := make([]string, 0, len(s.convertJobs))
	for _, job := range s.convertJobs {
		jobs = append(jobs, regexp.QuoteMeta(job))
	}

	stmt := tGksPhoneJMsg.
		SELECT(
			tGksPhoneJMsg.ID,
			tGksPhoneJMsg.Jobm,
			tGksPhoneJMsg.Anon,
			tGksPhoneJMsg.Gps,
			tGksPhoneJMsg.Message,
			tUsers.ID.AS("userid"),
		).
		FROM(
			tGksPhoneJMsg.
				INNER_JOIN(tGksPhoneSettings,
					tGksPhoneSettings.PhoneNumber.EQ(tGksPhoneJMsg.Number),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.EQ(tGksPhoneSettings.Identifier),
				),
		).
		WHERE(mysql.AND(
			tGksPhoneJMsg.ID.GT(mysql.Int32(lastID)),
			// Target job(s) are stored as JSON array like this: `["ambulance"]`
			tGksPhoneJMsg.Jobm.REGEXP_LIKE(
				mysql.String("\\[\"("+strings.Join(jobs, "|")+")\"\\]"),
			),
			tGksPhoneJMsg.Owner.EQ(mysql.Int32(0)),
		)).
		ORDER_BY(tGksPhoneJMsg.ID.ASC()).
		LIMIT(maxDispatchConvertCount)

	var dest []struct {
		*model.GksphoneJobMessage

		UserId int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return stats, fmt.Errorf("failed to query gksphone dispatches. %w", err)
		}
	}

	s.logger.Debug("converting gksphone dispatch to fivenet", zap.Int("dispatch_count", len(dest)))
	var errs error
	for _, msg := range dest {
		if msg.ID > stats.lastID {
			stats.lastID = msg.ID
		}
		job := ""
		if msg.Jobm != nil {
			job = strings.TrimSuffix(strings.TrimPrefix(*msg.Jobm, "[\""), "\"]")
		}
		if strings.TrimSpace(job) == "" {
			s.logger.Warn(
				"skipping gksphone dispatch with empty target job",
				zap.Int32("phone_dsp_id", msg.ID),
			)
			if err := s.closeGKSPhoneJobMsg(ctx, msg.ID); err != nil {
				stats.failed++
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to close gksphone dispatch. %w", err),
				)
				continue
			}
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}

		if msg.Gps == nil {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}
		gpsCoords, _ := strings.CutPrefix(*msg.Gps, "GPS: ")
		gpsSplit := strings.SplitN(gpsCoords, ", ", 2)
		if len(gpsSplit) != 2 {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}
		x, err := strconv.ParseFloat(gpsSplit[0], 32)
		if err != nil {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}
		y, err := strconv.ParseFloat(gpsSplit[1], 32)
		if err != nil {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}

		anon := false
		if msg.Anon != nil && *msg.Anon == "1" {
			anon = true
		}

		message := "N/A"
		if msg.Message != nil {
			message = utils.StringFirstNWithEllipsis(*msg.Message, 250)
		}

		dsp := &centrumdispatches.Dispatch{
			CreatedAt:  timestamp.Now(),
			Attributes: &centrumdispatches.DispatchAttributes{},
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						Name: job,
					},
				},
			},
			Message:   message,
			X:         x,
			Y:         y,
			Anon:      anon,
			CreatorId: &msg.UserId,
		}

		s.logger.Debug(
			"converted gksphone dispatch to fivenet",
			zap.String("job", job),
			zap.Int32("creator_id", msg.UserId),
			zap.Int32("phone_dsp_id", msg.ID),
		)

		if _, err := s.dispatchCreateFn(ctx, dsp); err != nil {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			errs = multierr.Append(
				errs,
				fmt.Errorf("failed to create dispatch for gksphone dispatch. %w", err),
			)
			continue
		}

		if err := s.closeGKSPhoneJobMsg(ctx, msg.ID); err != nil {
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			errs = multierr.Append(errs, fmt.Errorf("failed to close gksphone dispatch. %w", err))
			continue
		}
		stats.successful++
		s.metrics.processed.WithLabelValues(s.converterType, "successful").Inc()
	}

	if errs != nil {
		s.logger.Warn("some gksphone dispatches could not be converted", zap.Error(errs))
	}
	return stats, nil
}

func (s *Converter) closeGKSPhoneJobMsg(ctx context.Context, id int32) error {
	tGksPhoneJMsg := table.GksphoneJobMessage
	stmt := tGksPhoneJMsg.
		UPDATE(
			tGksPhoneJMsg.Owner,
		).
		SET(
			tGksPhoneJMsg.Owner.SET(mysql.Int32(1)),
		).
		WHERE(
			tGksPhoneJMsg.ID.EQ(mysql.Int32(id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Converter) convertLBPhoneJobMsgToDispatchWithCursor(
	ctx context.Context,
	lastID int32,
) (runStats, error) {
	stats := runStats{lastID: lastID}
	tPhoneServicesChannels := table.PhoneServicesChannels
	tPhoneServicesMessages := table.PhoneServicesMessages
	tPhonePhones := table.PhonePhones
	tUsers := table.FivenetUser

	targetJobsExp := make([]mysql.Expression, 0, len(s.convertJobs))
	for _, job := range s.convertJobs {
		targetJobsExp = append(targetJobsExp, mysql.String(job))
	}

	stmt := tPhoneServicesChannels.
		SELECT(
			tPhoneServicesChannels.ID.AS("id"),
			tPhoneServicesChannels.Company.AS("company"),
			tPhoneServicesChannels.PhoneNumber.AS("phone_number"),
			tPhoneServicesMessages.Message.AS("message"),
			tPhoneServicesMessages.XPos.AS("x_pos"),
			tPhoneServicesMessages.YPos.AS("y_pos"),
			tUsers.ID.AS("userid"),
		).
		FROM(
			tPhoneServicesChannels.
				INNER_JOIN(tPhoneServicesMessages,
					tPhoneServicesMessages.ChannelID.EQ(tPhoneServicesChannels.ID),
				).
				INNER_JOIN(tPhonePhones,
					tPhonePhones.PhoneNumber.EQ(tPhoneServicesChannels.PhoneNumber),
				).
				INNER_JOIN(tUsers,
					tUsers.Identifier.EQ(tPhonePhones.OwnerID),
				),
		).
		WHERE(mysql.AND(
			tPhoneServicesChannels.ID.GT(mysql.Int32(lastID)),
			// The target job is stored in the `company` field directly.
			tPhoneServicesChannels.Company.IN(targetJobsExp...),
		)).
		ORDER_BY(tPhoneServicesChannels.ID.ASC()).
		LIMIT(maxDispatchConvertCount)

	var dest []struct {
		ID          int32
		Job         string `alias:"company"`
		PhoneNumber string
		Message     string
		XPos        int32
		YPos        int32
		UserId      int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return stats, fmt.Errorf("failed to query lbphone dispatches. %w", err)
		}
	}

	s.logger.Debug("converting lbphone dispatch to fivenet", zap.Int("dispatch_count", len(dest)))
	var errs error
	for _, msg := range dest {
		if msg.ID > stats.lastID {
			stats.lastID = msg.ID
		}
		if strings.TrimSpace(msg.Job) == "" {
			s.logger.Warn(
				"skipping lbphone dispatch with empty target job",
				zap.Int32("phone_dsp_id", msg.ID),
			)
			if err := s.closeLBPhoneJobMsg(ctx, msg.ID); err != nil {
				stats.failed++
				errs = multierr.Append(
					errs,
					fmt.Errorf("failed to close lbphone dispatch (empty job). %w", err),
				)
				continue
			}
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			continue
		}

		message := "N/A"
		if msg.Message != "" {
			message = utils.StringFirstNWithEllipsis(msg.Message, 250)
		}

		dsp := &centrumdispatches.Dispatch{
			CreatedAt: timestamp.Now(),
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						Name: msg.Job,
					},
				},
			},
			Message:    message,
			X:          float64(msg.XPos),
			Y:          float64(msg.YPos),
			Anon:       false,
			Attributes: &centrumdispatches.DispatchAttributes{},
			CreatorId:  &msg.UserId,
		}

		s.logger.Debug(
			"converted lbphone dispatch to fivenet",
			zap.String("job", msg.Job),
			zap.Int32("creator_id", msg.UserId),
			zap.Int32("phone_dsp_id", msg.ID),
		)
		created := true
		if _, err := s.dispatchCreateFn(ctx, dsp); err != nil {
			created = false
			stats.failed++
			s.metrics.processed.WithLabelValues(s.converterType, "failed").Inc()
			errs = multierr.Append(
				errs,
				fmt.Errorf("failed to create dispatch for lbphone dispatch. %w", err),
			)
		}

		if err := s.closeLBPhoneJobMsg(ctx, msg.ID); err != nil {
			stats.failed++
			errs = multierr.Append(errs, fmt.Errorf("failed to close lbphone dispatch. %w", err))
			continue
		}
		if created {
			stats.successful++
			s.metrics.processed.WithLabelValues(s.converterType, "successful").Inc()
		}
	}

	if errs != nil {
		s.logger.Warn("some lbphone dispatches could not be converted", zap.Error(errs))
	}
	return stats, nil
}

func (s *Converter) closeLBPhoneJobMsg(ctx context.Context, id int32) error {
	tPhoneServicesChannels := table.PhoneServicesChannels
	stmt := tPhoneServicesChannels.
		DELETE().
		WHERE(
			tPhoneServicesChannels.ID.EQ(mysql.Int32(id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
