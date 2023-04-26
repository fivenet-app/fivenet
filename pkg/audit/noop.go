package audit

import (
	"context"

	"github.com/galexrt/fivenet/proto/resources/rector"
)

type Noop struct {
}

func (a *Noop) Log(ctx context.Context, service string, method string, state rector.EVENT_TYPE, targetUserId int32, data interface{}) {
	// Nothing
}
