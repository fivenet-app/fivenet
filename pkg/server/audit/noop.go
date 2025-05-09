package audit

import (
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

type Noop struct{}

func (a *Noop) Log(in *model.FivenetAuditLog, data any, callbacks ...FilterFn) {
	// Nothing
}
