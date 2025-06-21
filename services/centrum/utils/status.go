package centrumutils

import "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"

func IsStatusDispatchComplete(in centrum.StatusDispatch) bool {
	return in == centrum.StatusDispatch_STATUS_DISPATCH_ARCHIVED ||
		in == centrum.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
		in == centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED ||
		in == centrum.StatusDispatch_STATUS_DISPATCH_DELETED
}

func IsStatusDispatchUnassigned(in centrum.StatusDispatch) bool {
	return in == centrum.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED ||
		in == centrum.StatusDispatch_STATUS_DISPATCH_NEW ||
		in == centrum.StatusDispatch_STATUS_DISPATCH_UNASSIGNED
}

func IsDispatchUnassigned(in *centrum.Dispatch) bool {
	if in == nil {
		return false
	}

	// Ignore dispatches with no status
	if in.Status == nil {
		return false
	}

	// Ignore completed dispatches
	if IsStatusDispatchComplete(in.Status.Status) {
		return false
	}

	// Dispatch is "new" or unassgined, and no units assigned to it
	return (IsStatusDispatchUnassigned(in.Status.Status) || in.Status.Status == centrum.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED) &&
		len(in.Units) == 0
}
