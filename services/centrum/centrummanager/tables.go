package centrummanager

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
)

var (
	tUnits             = table.FivenetCentrumUnits.AS("unit")
	tUnitStatus        = table.FivenetCentrumUnitsStatus.AS("unit_status")
	tUnitUser          = table.FivenetCentrumUnitsUsers.AS("unit_assignment")
	tUserProps         = table.FivenetUserProps
	tColleagueProps    = table.FivenetJobColleagueProps.AS("jobs_user_props")
	tCentrumSettings   = table.FivenetCentrumSettings
	tCentrumDisponents = table.FivenetCentrumDisponents
	tDispatch          = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus    = table.FivenetCentrumDispatchesStatus.AS("dispatch_status")
	tDispatchUnit      = table.FivenetCentrumDispatchesAsgmts.AS("dispatch_assignment")
)
