package centrumutils

import "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"

func IsStatusDispatchComplete(in dispatch.StatusDispatch) bool {
	return in == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED
}

func IsStatusDispatchUnassigned(in dispatch.StatusDispatch) bool {
	return in == dispatch.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_NEW ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED
}

func IsDispatchUnassigned(in *dispatch.Dispatch) bool {
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

	// Dispatch is "new" or unassgined
	if IsStatusDispatchUnassigned(in.Status.Status) {
		return true
	}

	return len(in.Units) == 0
}
