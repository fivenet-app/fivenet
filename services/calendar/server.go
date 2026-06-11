package calendar

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	discordstate "github.com/diamondburned/arikawa/v3/state"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/i18n"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	pkggrpc "github.com/fivenet-app/fivenet/v2026/pkg/grpc"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	tCalendar     = table.FivenetCalendar.AS("calendar")
	tCalendarSubs = table.FivenetCalendarSubs.AS("calendar_sub")

	tCalendarEntry          = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP           = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")
	tCalendarRSVPOccurrence = table.FivenetCalendarRsvpOccurrence.AS(
		"calendar_entry_rsvp_occurrence",
	)

	tCAccess = table.FivenetCalendarAccess.AS("calendar_access")

	tUserJobs  = table.FivenetUserJobs.AS("user_jobs")
	tUserProps = table.FivenetUserProps
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCalendar,
		JobColumn:       table.FivenetCalendar.Job,
		DeletedAtColumn: table.FivenetCalendar.DeletedAt,
		IDColumn:        table.FivenetCalendar.ID,

		MinDays: 60,

		DependantTables: []*housekeeper.Table{
			{
				Table:      table.FivenetCalendarSubs,
				ForeignKey: table.FivenetCalendarSubs.CalendarID,
			},
			{
				Table:           table.FivenetCalendarEntries,
				IDColumn:        table.FivenetCalendarEntries.ID,
				JobColumn:       table.FivenetCalendarEntries.Job,
				ForeignKey:      table.FivenetCalendarEntries.CalendarID,
				DeletedAtColumn: table.FivenetCalendarEntries.DeletedAt,

				MinDays: 60,

				DependantTables: []*housekeeper.Table{
					{
						Table:      table.FivenetCalendarRsvp,
						ForeignKey: table.FivenetCalendarRsvp.EntryID,
					},
					{
						Table:      table.FivenetCalendarRsvpOccurrence,
						ForeignKey: table.FivenetCalendarRsvpOccurrence.EntryID,
					},
					{
						Table:      table.FivenetCalendarDiscordReminderSends,
						ForeignKey: table.FivenetCalendarDiscordReminderSends.EntryID,
					},
				},
			},
		},
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetCalendarDiscordReminderSends,
		IDColumn:        table.FivenetCalendarDiscordReminderSends.ID,
		TimestampColumn: table.FivenetCalendarDiscordReminderSends.CreatedAt,
		MinDays:         30,
	})
}

type Server struct {
	pbcalendar.CalendarServiceServer
	pbcalendar.EntriesServiceServer

	logger   *zap.Logger
	tracer   trace.Tracer
	db       *sql.DB
	cfg      *config.Config
	ps       perms.Permissions
	enricher *mstlystcdata.UserAwareEnricher
	appCfg   appconfig.IConfig
	i18n     i18n.Ii18n
	notif    notifi.INotifi
	js       *events.JSWrapper
	dc       *discordstate.State

	access *access.Grouped[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, access.DummyQualificationAccess[calendaraccess.AccessLevel], *access.DummyQualificationAccess[calendaraccess.AccessLevel], calendaraccess.AccessLevel]
}

type Params struct {
	fx.In

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Config    *config.Config
	P         perms.Permissions
	Enricher  *mstlystcdata.UserAwareEnricher
	AppConfig appconfig.IConfig
	I18n      i18n.Ii18n
	Notif     notifi.INotifi
	JS        *events.JSWrapper
	Discord   *discordstate.State
}

type Result struct {
	fx.Out

	Server       *Server
	Service      pkggrpc.Service     `group:"grpcservices"`
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewServer(p Params) Result {
	s := &Server{
		logger:   p.Logger.Named("calendar"),
		tracer:   p.TP.Tracer("calendar"),
		db:       p.DB,
		cfg:      p.Config,
		ps:       p.P,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
		i18n:     p.I18n,
		notif:    p.Notif,
		js:       p.JS,
		dc:       p.Discord,

		access: access.NewGrouped[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, access.DummyQualificationAccess[calendaraccess.AccessLevel], *access.DummyQualificationAccess[calendaraccess.AccessLevel], calendaraccess.AccessLevel](
			p.DB,
			table.FivenetDocuments,
			&access.TargetTableColumns{
				ID:         table.FivenetDocuments.ID,
				DeletedAt:  table.FivenetDocuments.DeletedAt,
				CreatorJob: table.FivenetDocuments.CreatorJob,
				CreatorID:  table.FivenetDocuments.CreatorID,
			},
			access.NewJobs[calendaraccess.CalendarJobAccess, *calendaraccess.CalendarJobAccess, calendaraccess.AccessLevel](
				table.FivenetCalendarAccess,
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.ID,
						TargetID: table.FivenetCalendarAccess.TargetID,
						Access:   table.FivenetCalendarAccess.Access,
					},
					Job:          table.FivenetCalendarAccess.Job,
					MinimumGrade: table.FivenetCalendarAccess.MinimumGrade,
				},
				table.FivenetCalendarAccess.AS("calendar_job_access"),
				&access.JobAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.AS("calendar_job_access").ID,
						TargetID: table.FivenetCalendarAccess.AS("calendar_job_access").TargetID,
						Access:   table.FivenetCalendarAccess.AS("calendar_job_access").Access,
					},
					Job: table.FivenetCalendarAccess.AS("calendar_job_access").Job,
					MinimumGrade: table.FivenetCalendarAccess.AS(
						"calendar_job_access",
					).MinimumGrade,
				},
			),
			access.NewUsers[calendaraccess.CalendarUserAccess, *calendaraccess.CalendarUserAccess, calendaraccess.AccessLevel](
				table.FivenetCalendarAccess,
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.ID,
						TargetID: table.FivenetCalendarAccess.TargetID,
						Access:   table.FivenetCalendarAccess.Access,
					},
					UserID: table.FivenetCalendarAccess.UserID,
				},
				table.FivenetCalendarAccess.AS("calendar_user_access"),
				&access.UserAccessColumns{
					BaseAccessColumns: access.BaseAccessColumns{
						ID:       table.FivenetCalendarAccess.AS("calendar_user_access").ID,
						TargetID: table.FivenetCalendarAccess.AS("calendar_user_access").TargetID,
						Access:   table.FivenetCalendarAccess.AS("calendar_user_access").Access,
					},
					UserID: table.FivenetCalendarAccess.AS("calendar_user_access").UserID,
				},
			),
			nil,
		),
	}

	return Result{
		Server:       s,
		Service:      s,
		CronRegister: s,
	}
}

func (s *Server) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "calendar.rsvp.cleanup",
		Schedule: "32 */6 * * *",
		Timeout:  durationpb.New(1 * time.Minute),
	}); err != nil {
		return err
	}

	if err := registry.UnregisterCronjob(ctx, "calendar.discord_reminders"); err != nil {
		return err
	}

	return nil
}

func (s *Server) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("calendar.rsvp.cleanup", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := s.tracer.Start(ctx, "calendar.rsvp.cleanup")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if err := data.Unmarshal(dest); err != nil {
			s.logger.Warn(
				"failed to unmarshal calendar rsvp occurrences cleanup cron data",
				zap.Error(err),
			)
		}

		rowsAffected, err := s.cleanupCalendarRSVPOccurrences(ctx)
		if err != nil {
			s.logger.Error(
				"failed to generate calendar rsvp occurrences cleanup",
				zap.Error(err),
			)
			return err
		}
		dest.SetAttribute("changed_rows", strconv.Itoa(int(rowsAffected)))

		// Marshal the updated cron data
		if err := data.MarshalFrom(dest); err != nil {
			return fmt.Errorf(
				"failed to marshal updated calendar rsvp occurrences cleanup cron data. %w",
				err,
			)
		}

		return nil
	})

	return nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbcalendar.RegisterCalendarServiceServer(srv, s)
	pbcalendar.RegisterEntriesServiceServer(srv, s)
}
