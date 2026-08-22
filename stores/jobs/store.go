package jobsstore

import (
	"context"
	"database/sql"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
)

type ListColleaguesQuery struct {
	Job    string
	Search string
	Where  mysql.BoolExpression

	UserIDs    []int32
	UserOnly   bool
	Absent     bool
	LabelIDs   []int64
	NamePrefix string
	NameSuffix string

	Sort   *database.Sort
	Offset int64
	Limit  int64
}

type ListColleaguesByUserIDsQuery struct {
	Job         string
	UserIDs     []int32
	WithColumns mysql.ProjectionList
}

type ListQuery struct {
	Job string

	Where  mysql.BoolExpression
	Sort   *database.Sort
	Offset int64
	Limit  int64
}

type TimeclockQuery struct {
	UserMode jobstimeclock.TimeclockViewMode
	Mode     jobstimeclock.TimeclockMode
	Date     *database.DateRange
	PerDay   bool
	UserIDs  []int32
	Sort     *database.Sort
	Offset   int64
	Limit    int64
	Job      string
	UserID   int32
}

type InactiveEmployeesQuery struct {
	Days    int32
	UserIDs []int32
	Sort    *database.Sort
	Offset  int64
	Limit   int64
	Job     string
}

type ConductQuery struct {
	Sort           *database.Sort
	Offset         int64
	Limit          int64
	Job            string
	Types          []jobsconduct.ConductType
	ShowExpired    bool
	ShowDrafts     bool
	UserIDs        []int32
	IDs            []int64
	CreatorID      int32
	OwnOnly        bool
	AllAccess      bool
	IncludeDeleted bool
}

type GroupsQuery struct {
	Job             string
	Kind            jobsgroups.GroupType
	States          []jobsgroups.GroupState
	Search          string
	IncludeCounts   bool
	IncludeInactive bool
	IncludeArchived bool
	IDs             []int64
	Sort            *database.Sort
	Offset          int64
	Limit           int64
}

type GroupQuery struct {
	Job             string
	IncludeArchived bool
}

type GroupItemsQuery struct {
	GroupID int64
	Search  string
	Offset  int64
	Limit   int64
}

type GroupRuleMemberMatch struct {
	GroupID int64
	UserID  int32
	RuleID  int64
	Label   string
}

type IGroupsQuery interface {
	GetGroup(ctx context.Context, db qrm.DB, q GroupQuery, id int64) (*jobsgroups.Group, error)
	ListGroupManualMembers(
		ctx context.Context,
		db qrm.DB,
		q GroupItemsQuery,
	) ([]*jobsgroups.GroupManualMember, error)
	ListGroupRuleMemberMatches(
		ctx context.Context,
		db qrm.DB,
		group *jobsgroups.Group,
		search string,
	) ([]*GroupRuleMemberMatch, error)
	ListGroupMemberExclusions(
		ctx context.Context,
		db qrm.DB,
		q GroupItemsQuery,
	) ([]*jobsgroups.GroupMemberExclusion, error)
	ListGroupLeaders(
		ctx context.Context,
		db qrm.DB,
		q GroupItemsQuery,
	) ([]*jobsgroups.GroupLeader, error)
}

type IStore interface {
	IGroupsQuery

	GetMOTD(ctx context.Context, db qrm.DB, job string) (string, error)
	SetMOTD(ctx context.Context, db qrm.DB, job string, motd string) error
	GetJobProps(ctx context.Context, db qrm.DB, job string) (*jobsprops.JobProps, error)

	CountGroups(
		ctx context.Context,
		db qrm.DB,
		q GroupsQuery,
		userInfo *userinfo.UserInfo,
	) (int64, error)
	ListGroups(
		ctx context.Context,
		db qrm.DB,
		q GroupsQuery,
		userInfo *userinfo.UserInfo,
	) ([]*jobsgroups.Group, error)
	CreateGroup(ctx context.Context, db qrm.DB, group *jobsgroups.Group) (int64, error)
	UpdateGroup(ctx context.Context, db qrm.DB, group *jobsgroups.Group) error
	ArchiveGroup(ctx context.Context, db qrm.DB, job string, id int64, updatedByUserID int32) error
	RestoreGroup(ctx context.Context, db qrm.DB, job string, id int64, updatedByUserID int32) error
	RecountGroupStats(ctx context.Context, db qrm.DB, groupID int64) error
	CountGroupManualMembers(ctx context.Context, db qrm.DB, q GroupItemsQuery) (int64, error)
	AddGroupManualMember(
		ctx context.Context,
		db qrm.DB,
		groupID int64,
		userID int32,
		createdByUserID int32,
		reason *string,
	) (*jobsgroups.GroupManualMember, bool, error)
	RemoveGroupManualMember(ctx context.Context, db qrm.DB, groupID int64, userID int32) error
	CountGroupMemberExclusions(ctx context.Context, db qrm.DB, q GroupItemsQuery) (int64, error)
	AddGroupMemberExclusion(
		ctx context.Context,
		db qrm.DB,
		groupID int64,
		userID int32,
		reasonType jobsgroups.GroupExclusionReason,
		createdByUserID int32,
		reason *string,
	) (*jobsgroups.GroupMemberExclusion, bool, error)
	RemoveGroupMemberExclusion(ctx context.Context, db qrm.DB, groupID int64, userID int32) error
	CountGroupLeaders(ctx context.Context, db qrm.DB, q GroupItemsQuery) (int64, error)
	AddGroupLeader(
		ctx context.Context,
		db qrm.DB,
		groupID int64,
		userID int32,
		createdByUserID int32,
	) (*jobsgroups.GroupLeader, bool, error)
	RemoveGroupLeader(ctx context.Context, db qrm.DB, groupID int64, userID int32) error
	CountGroupRules(ctx context.Context, db qrm.DB, q GroupItemsQuery) (int64, error)
	ListGroupRules(
		ctx context.Context,
		db qrm.DB,
		q GroupItemsQuery,
	) ([]*jobsgroups.GroupRule, error)
	GetGroupRule(
		ctx context.Context,
		db qrm.DB,
		groupID int64,
		ruleID int64,
	) (*jobsgroups.GroupRule, error)
	CreateGroupRule(
		ctx context.Context,
		db qrm.DB,
		rule *jobsgroups.GroupRule,
	) (*jobsgroups.GroupRule, error)
	UpdateGroupRule(
		ctx context.Context,
		db qrm.DB,
		rule *jobsgroups.GroupRule,
		updatedByUserID int32,
	) (*jobsgroups.GroupRule, error)
	DeleteGroupRule(ctx context.Context, db qrm.DB, groupID int64, ruleID int64) error
	CreateGroupActivity(
		ctx context.Context,
		db qrm.DB,
		activities ...*jobsgroups.GroupActivity,
	) error
	CountGroupActivity(ctx context.Context, db qrm.DB, q ListQuery) (int64, error)
	ListGroupActivity(
		ctx context.Context,
		db qrm.DB,
		q ListQuery,
	) ([]*jobsgroups.GroupActivity, error)

	UserInJob(ctx context.Context, db qrm.DB, job string, userID int32) (bool, error)
	CountColleagues(ctx context.Context, db qrm.DB, q ListColleaguesQuery) (int64, error)
	ListColleagues(
		ctx context.Context,
		db qrm.DB,
		q ListColleaguesQuery,
	) ([]*jobscolleagues.Colleague, error)
	ListColleaguesByUserIDs(
		ctx context.Context,
		db qrm.DB,
		q ListColleaguesByUserIDsQuery,
	) ([]*jobscolleagues.Colleague, error)
	GetColleague(
		ctx context.Context,
		db qrm.DB,
		job string,
		userId int32,
		withColumns mysql.ProjectionList,
		includeDeleted bool,
	) (*jobscolleagues.Colleague, error)
	GetColleagueProps(
		ctx context.Context,
		db qrm.DB,
		job string,
		userId int32,
		fields []string,
		includeDeleted bool,
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
		includeDeleted bool,
	) (*jobslabels.Labels, error)
	GetUsersLabels(
		ctx context.Context,
		db qrm.DB,
		job string,
		userIds []int32,
		includeDeleted bool,
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
		includeDeleted bool,
	) ([]*jobslabels.Label, error)
	GetLabel(
		ctx context.Context,
		db qrm.DB,
		job string,
		labelId int64,
		includeDeleted bool,
	) (*jobslabels.Label, error)
	NextLabelSortOrder(
		ctx context.Context,
		db qrm.Queryable,
		job string,
	) (int32, error)
	UpdateLabel(ctx context.Context, db qrm.DB, label *jobslabels.Label, job string) error
	InsertLabel(ctx context.Context, db qrm.DB, label *jobslabels.Label) (int64, error)
	DeleteLabel(
		ctx context.Context,
		db qrm.DB,
		job string,
		labelId int64,
		deletedAt *timestamp.Timestamp,
	) error
	ReorderLabels(
		ctx context.Context,
		job string,
		labelIds []int64,
	) error
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
	GetConductEntry(
		ctx context.Context,
		db qrm.DB,
		id int64,
		includeDeleted bool,
	) (*jobsconduct.ConductEntry, error)
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
	db          *sql.DB
	customDB    *config.CustomDB
	groupAccess access.JobGroupsAccess
}

type Result struct {
	fx.Out

	Store       IStore
	GroupsQuery IGroupsQuery
}

func New(db *sql.DB, customDB *config.CustomDB, groupAccess *access.JobGroupsObjectAccess) Result {
	s := &Store{
		db:          db,
		customDB:    customDB,
		groupAccess: groupAccess,
	}
	return Result{
		Store:       s,
		GroupsQuery: s,
	}
}
