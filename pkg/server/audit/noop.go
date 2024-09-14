package audit

import (
	"github.com/fivenet-app/fivenet/query/fivenet/model"
)

type Noop struct{}

func (a *Noop) Log(in *model.FivenetAuditLog, data any, callbacks ...FilterFn) {
	// Nothing
}
