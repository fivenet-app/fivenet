package jobs

import (
	"context"
	"database/sql"

	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	GetMOTD(ctx context.Context, db qrm.DB, job string) (string, error)
	SetMOTD(ctx context.Context, db qrm.DB, job string, motd string) error
	GetJobProps(ctx context.Context, db qrm.DB, job string) (*jobsprops.JobProps, error)

	CountColleagues(ctx context.Context, db qrm.DB, q ListColleaguesQuery) (int64, error)
	ListColleagues(
		ctx context.Context,
		db qrm.DB,
		q ListColleaguesQuery,
	) ([]*jobscolleagues.Colleague, error)
	GetColleague(
		ctx context.Context,
		db qrm.DB,
		job string,
		userId int32,
		withColumns mysql.ProjectionList,
	) (*jobscolleagues.Colleague, error)
	GetColleagueProps(
		ctx context.Context,
		db qrm.DB,
		job string,
		userId int32,
		fields []string,
	) (*jobscolleagues.ColleagueProps, error)
	HandleColleaguePropsChanges(
		ctx context.Context,
		db qrm.DB,
		x *jobscolleagues.ColleagueProps,
		in *jobscolleagues.ColleagueProps,
		job string,
		sourceUserId *int32,
		reason string,
	) ([]*colleaguesactivity.ColleagueActivity, error)
	CreateColleagueActivity(
		ctx context.Context,
		db qrm.DB,
		activities ...*colleaguesactivity.ColleagueActivity,
	) error
	ValidateLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		labels []*jobslabels.Label,
	) (bool, error)
	GetUserLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		userId int32,
	) (*jobslabels.Labels, error)
	GetUsersLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		userIds []int32,
	) ([]*UserLabels, error)
	CountColleagueActivity(ctx context.Context, db qrm.DB, q ListQuery) (int64, error)
	ListColleagueActivity(
		ctx context.Context,
		db qrm.DB,
		q ListQuery,
	) ([]*colleaguesactivity.ColleagueActivity, error)

	GetColleagueLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		search string,
	) ([]*jobslabels.Label, error)
	ManageLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		labels []*jobslabels.Label,
	) ([]*jobslabels.Label, error)
	GetColleagueLabelsStats(
		ctx context.Context,
		db qrm.DB,
		job string,
	) ([]*jobslabels.LabelCount, error)

	CountTimeclock(ctx context.Context, db qrm.DB, q TimeclockQuery) (int64, error)
	ListTimeclock(
		ctx context.Context,
		db qrm.DB,
		q TimeclockQuery,
	) ([]*jobstimeclock.TimeclockEntry, error)
	ListTimeclockTimeline(
		ctx context.Context,
		db qrm.DB,
		q TimeclockQuery,
	) ([]*jobstimeclock.TimeclockEntry, error)
	GetTimeclockStats(
		ctx context.Context,
		db qrm.DB,
		q TimeclockQuery,
	) (*jobstimeclock.TimeclockStats, error)
	GetTimeclockWeeklyStats(
		ctx context.Context,
		db qrm.DB,
		q TimeclockQuery,
	) ([]*jobstimeclock.TimeclockWeeklyStats, error)
	CountInactiveEmployees(ctx context.Context, db qrm.DB, q InactiveEmployeesQuery) (int64, error)
	ListInactiveEmployees(
		ctx context.Context,
		db qrm.DB,
		q InactiveEmployeesQuery,
	) ([]*jobscolleagues.Colleague, error)
	CleanupTimeclock(ctx context.Context, db qrm.DB) error

	CountConductEntries(ctx context.Context, db qrm.DB, q ConductQuery) (int64, error)
	ListConductEntries(
		ctx context.Context,
		db qrm.DB,
		q ConductQuery,
	) ([]*jobsconduct.ConductEntry, error)
	GetConductEntry(ctx context.Context, db qrm.DB, id int64) (*jobsconduct.ConductEntry, error)
	CreateConductEntry(
		ctx context.Context,
		db qrm.DB,
		entry *jobsconduct.ConductEntry,
	) (int64, error)
	UpdateConductEntry(ctx context.Context, db qrm.DB, entry *jobsconduct.ConductEntry) error
	DeleteConductEntry(
		ctx context.Context,
		db qrm.DB,
		job string,
		id int64,
		deletedAt *timestamp.Timestamp,
	) error
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{db: db, customDB: customDB}
}
