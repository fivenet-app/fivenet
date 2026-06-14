package jobs

import "github.com/fivenet-app/fivenet/v2026/query/fivenet/table"

var (
	tUserProps         = table.FivenetUserProps
	tUserJobs          = table.FivenetUserJobs
	tJobProps          = table.FivenetJobProps
	tJobLabels         = table.FivenetJobLabels.AS("label")
	tColleagueLabels   = table.FivenetJobColleagueLabels
	tColleagueProps    = table.FivenetJobColleagueProps.AS("colleague_props")
	tColleagueActivity = table.FivenetJobColleagueActivity
	tAvatar            = table.FivenetFiles.AS("profile_picture")
	tConduct           = table.FivenetJobConduct.AS("conduct_entry")
	tTimeClock         = table.FivenetJobTimeclock.AS("timeclock_entry")
)

const (
	nameColumn = "name"
	rankColumn = "rank"
	timeColumn = "time"
)
