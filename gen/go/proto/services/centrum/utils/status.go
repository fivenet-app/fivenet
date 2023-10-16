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

	return len(in.Units) == 0
}
