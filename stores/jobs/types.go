package jobsstore

import (
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	jobsprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/props"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

type (
	JobProps             = jobsprops.JobProps
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
