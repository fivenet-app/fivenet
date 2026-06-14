package jobs

import (
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/go-jet/jet/v2/mysql"
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
	Days   int32
	Sort   *database.Sort
	Offset int64
	Limit  int64
	Job    string
}

type ConductQuery struct {
	Sort        *database.Sort
	Offset      int64
	Limit       int64
	Job         string
	Types       []jobsconduct.ConductType
	ShowExpired bool
	ShowDrafts  bool
	UserIDs     []int32
	IDs         []int64
	CreatorID   int32
	OwnOnly     bool
	AllAccess   bool
}
