package centrummanager

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
)

var (
	tUnits           = table.FivenetCentrumUnits.AS("unit")
	tUnitStatus      = table.FivenetCentrumUnitsStatus.AS("unitstatus")
	tUnitUser        = table.FivenetCentrumUnitsUsers.AS("unitassignment")
	tUserProps       = table.FivenetUserProps
	tJobsUserProps   = table.FivenetJobsUserProps.AS("jobs_user_props")
	tCentrumSettings = table.FivenetCentrumSettings
	tCentrumUsers    = table.FivenetCentrumUsers
	tDispatch        = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus  = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchUnit    = table.FivenetCentrumDispatchesAsgmts.AS("dispatchassignment")
)
