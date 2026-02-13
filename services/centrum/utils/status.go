package centrumutils

import (
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
)

func IsStatusDispatchComplete(in centrumdispatches.StatusDispatch) bool {
	return in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_ARCHIVED ||
		in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
		in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED ||
		in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_DELETED
}

func IsStatusDispatchUnassigned(in centrumdispatches.StatusDispatch) bool {
	return in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED ||
		in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW ||
		in == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNASSIGNED
}

func IsDispatchUnassigned(in *centrumdispatches.Dispatch) bool {
	if in == nil {
		return false
	}

	// Ignore dispatches with no status
	if in.GetStatus() == nil {
		return false
	}

	// Ignore completed dispatches
	if IsStatusDispatchComplete(in.GetStatus().GetStatus()) {
		return false
	}

	// Dispatch is "new" or unassgined, and no units assigned to it
	return (IsStatusDispatchUnassigned(in.GetStatus().GetStatus()) || in.GetStatus().GetStatus() == centrumdispatches.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) &&
		len(in.GetUnits()) == 0
}
