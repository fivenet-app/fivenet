package centrumutils

import "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"

func IsStatusDispatchComplete(in dispatch.StatusDispatch) bool {
	return in == dispatch.StatusDispatch_STATUS_DISPATCH_ARCHIVED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_CANCELLED ||
		in == dispatch.StatusDispatch_STATUS_DISPATCH_COMPLETED
}

func IsDispatchUnassigned(in *dispatch.Dispatch) bool {
	if in == nil {
		return false
	}

	if in.Status != nil &&
		(in.Status.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED ||
			in.Status.Status == dispatch.StatusDispatch_STATUS_DISPATCH_NEW ||
			in.Status.Status == dispatch.StatusDispatch_STATUS_DISPATCH_UNASSIGNED) {
		return true
	}

	return len(in.Units) == 0
}
