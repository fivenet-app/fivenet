package audit

import (
	"github.com/galexrt/fivenet/query/fivenet/model"
)

type Noop struct {
}

func (a *Noop) Log(in *model.FivenetAuditLog, data any) {
	// Nothing
}
