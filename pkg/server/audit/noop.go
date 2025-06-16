package audit

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
)

// Noop implements IAuditer but does nothing. Useful for disabling audit logging in certain environments.
type Noop struct{}

// Log is a no-op implementation that ignores all input and does nothing.
func (a *Noop) Log(in *audit.AuditEntry, data any, callbacks ...FilterFn) {
	// Nothing
}
