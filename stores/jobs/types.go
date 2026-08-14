package jobsstore

import (
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

type (
	JobProps             = jobsprops.JobProps
	Group                = jobsgroups.Group
	GroupLeader          = jobsgroups.GroupLeader
	GroupManualMember    = jobsgroups.GroupManualMember
	GroupMemberExclusion = jobsgroups.GroupMemberExclusion
	GroupRule            = jobsgroups.GroupRule
	GroupResolvedMember  = jobsgroups.GroupResolvedMember
	Colleague            = jobscolleagues.Colleague
	ColleagueProps       = jobscolleagues.ColleagueProps
	ColleagueActivity    = colleaguesactivity.ColleagueActivity
	Label                = jobslabels.Label
	Labels               = jobslabels.Labels
	LabelCount           = jobslabels.LabelCount
	TimeclockEntry       = jobstimeclock.TimeclockEntry
	TimeclockStats       = jobstimeclock.TimeclockStats
	TimeclockWeeklyStats = jobstimeclock.TimeclockWeeklyStats
	ConductEntry         = jobsconduct.ConductEntry
	Timestamp            = timestamp.Timestamp
)
