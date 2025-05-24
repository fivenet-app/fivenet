package audit

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
)

type Noop struct{}

func (a *Noop) Log(in *audit.AuditEntry, data any, callbacks ...FilterFn) {
	// Nothing
}
