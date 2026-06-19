package calendarstore

import (
	"context"
	"database/sql"

	calendarresource "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar"
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	calendarentries "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/entries"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbcalendar "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/calendar"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tCalendar     = table.FivenetCalendar.AS("calendar")
	tCalendarSubs = table.FivenetCalendarSubs.AS("calendar_sub")

	tUserProps = table.FivenetUserProps

	tCalendarEntry          = table.FivenetCalendarEntries.AS("calendar_entry")
	tCalendarRSVP           = table.FivenetCalendarRsvp.AS("calendar_entry_rsvp")
	tCalendarRSVPOccurrence = table.FivenetCalendarRsvpOccurrence.AS(
		"calendar_entry_rsvp_occurrence",
	)
)

type ListQuery struct {
	UserInfo       *userinfo.UserInfo
	MinAccessLevel *calendaraccess.AccessLevel
	OnlyPublic     bool
	After          *timestamp.Timestamp
	CalendarIDs    []int64
}

type ListCalendarEntriesOptions struct {
	ShowHidden  bool
	After       *timestamp.Timestamp
	Year        int32
	Month       int32
	CalendarIDs []int64
}

type GetUpcomingEntriesOptions struct {
	Seconds int64
}

type ListCalendarEntryRSVPOptions struct {
	EntryID    int64
	Pagination *database.PaginationRequest
}

type IStore interface {
	CountCalendars(ctx context.Context, q ListQuery) (int64, error)
	ListCalendars(
		ctx context.Context,
		q ListQuery,
		offset, limit int64,
	) ([]*calendarresource.Calendar, error)
	GetCalendar(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
		condition mysql.BoolExpression,
	) (*calendarresource.Calendar, error)
	GetAccessibleCalendar(
		ctx context.Context,
		calendarID int64,
		userInfo *userinfo.UserInfo,
		accessLevel calendaraccess.AccessLevel,
		publicOk bool,
	) (*calendarresource.Calendar, error)
	ListTargetAccess(
		ctx context.Context,
		calendarID int64,
		options access.SubjectAccessOptions,
	) (*calendaraccess.CalendarAccess, error)
	CheckIfUserHasAccessToCalendar(
		ctx context.Context,
		calendarID int64,
		userInfo *userinfo.UserInfo,
		accessLevel calendaraccess.AccessLevel,
		publicOk bool,
	) (bool, error)
	CheckIfUserHasAccessToCalendarEntry(
		ctx context.Context,
		calendarID int64,
		entryID int64,
		userInfo *userinfo.UserInfo,
		accessLevel calendaraccess.AccessLevel,
		publicOk bool,
	) (bool, error)
	CreateCalendar(
		ctx context.Context,
		tx qrm.DB,
		cal *calendarresource.Calendar,
		userInfo *userinfo.UserInfo,
		discordSettingsJSON *string,
	) (int64, error)
	UpdateCalendar(
		ctx context.Context,
		tx qrm.DB,
		cal *calendarresource.Calendar,
		discordSettingsJSON *string,
	) error
	DeleteCalendar(
		ctx context.Context,
		tx qrm.DB,
		calendarID int64,
		deletedAt *timestamp.Timestamp,
	) error
	SetSubscription(
		ctx context.Context,
		calendarID int64,
		userID int32,
		subscribe bool,
		confirmed bool,
		muted bool,
	) error
	GetCalendarSub(
		ctx context.Context,
		userID int32,
		calendarID int64,
	) (*calendarresource.CalendarSub, error)
	EnsureBirthdayCalendarAccess(
		ctx context.Context,
		q qrm.DB,
		calendarID int64,
		job string,
		jobInfo *jobs.Job,
	) error
	ListCalendarEntries(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
		opts ListCalendarEntriesOptions,
	) ([]*calendarentries.CalendarEntry, error)
	GetUpcomingEntries(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
		opts GetUpcomingEntriesOptions,
	) ([]*calendarentries.CalendarEntry, error)
	FilterUpcomingCalendarEntries(
		entries []*calendarentries.CalendarEntry,
		userInfo *userinfo.UserInfo,
	) []*calendarentries.CalendarEntry
	GetEntry(
		ctx context.Context,
		userInfo *userinfo.UserInfo,
		condition mysql.BoolExpression,
	) (*calendarentries.CalendarEntry, error)
	UpsertCalendarEntry(
		ctx context.Context,
		tx qrm.DB,
		entry *calendarentries.CalendarEntry,
		oldEntry *calendarentries.CalendarEntry,
		userInfo *userinfo.UserInfo,
	) (int64, error)
	DeleteCalendarEntry(
		ctx context.Context,
		tx qrm.DB,
		entryID int64,
		calendarID int64,
		deletedAt *timestamp.Timestamp,
	) error
	ShareCalendarEntry(
		ctx context.Context,
		tx qrm.DB,
		entryID int64,
		inUserIds []int32,
	) ([]int32, error)
	GetUserShortByID(ctx context.Context, userID int32) (*usershort.UserShort, error)
	ListCalendarEntryRSVP(
		ctx context.Context,
		opts ListCalendarEntryRSVPOptions,
		userInfo *userinfo.UserInfo,
	) (*pbcalendar.ListCalendarEntryRSVPResponse, error)
	GetRSVPCalendarEntry(
		ctx context.Context,
		entryID int64,
		userID int32,
		occurrenceKey string,
	) (*calendarentries.CalendarEntryRSVP, error)
	SetCalendarEntryRSVP(
		ctx context.Context,
		entry *calendarentries.CalendarEntryRSVP,
		userInfo *userinfo.UserInfo,
		occurrenceKey string,
		remove bool,
	) error
	ValidateRecurringOccurrenceKey(entry *calendarentries.CalendarEntry, occurrenceKey string) error
	GetCalendarReminderGuildID(ctx context.Context, job string) (string, error)
	CleanupCalendarRSVPOccurrences(ctx context.Context) (int64, error)

	LoadBirthdayColleagues(
		ctx context.Context,
		tx qrm.DB,
		job string,
	) ([]*BirthdayColleague, error)
	InsertBirthdayEntry(
		ctx context.Context,
		tx qrm.DB,
		calendarID int64,
		job string,
		colleague *BirthdayColleague,
	) error
}

type Store struct {
	db             *sql.DB
	access         *access.SubjectObjectAccess
	accessResolver *access.SubjectResolver
}

func New(db *sql.DB) IStore {
	return &Store{
		db:             db,
		access:         access.NewCalendarSubjectObjectAccess(db),
		accessResolver: access.NewSubjectResolver(db),
	}
}

func (s *Store) ListTargetAccess(
	ctx context.Context,
	calendarID int64,
	options access.SubjectAccessOptions,
) (*calendaraccess.CalendarAccess, error) {
	return s.access.ListTargetAccess(ctx, s.db, calendarID, options)
}
