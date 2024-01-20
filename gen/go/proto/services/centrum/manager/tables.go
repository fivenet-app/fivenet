package manager

import "github.com/galexrt/fivenet/query/fivenet/table"

var (
	tUnits           = table.FivenetCentrumUnits.AS("unit")
	tUnitStatus      = table.FivenetCentrumUnitsStatus.AS("unitstatus")
	tUnitUser        = table.FivenetCentrumUnitsUsers.AS("unitassignment")
	tUsers           = table.Users.AS("usershort")
	tCentrumSettings = table.FivenetCentrumSettings
	tCentrumUsers    = table.FivenetCentrumUsers
	tDispatch        = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus  = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchUnit    = table.FivenetCentrumDispatchesAsgmts.AS("dispatchassignment")

	// Converter
	tGksPhoneJMsg     = table.GksphoneJobMessage
	tGksPhoneSettings = table.GksphoneSettings
)
