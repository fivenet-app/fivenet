package audit

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

type Noop struct {
}

func (a *Noop) Log(service string, method string, state rector.EVENT_TYPE, targetUserId int32, data any) {
	// Nothing
}

func (a *Noop) AddEntry(in *model.FivenetAuditLog) {
	// Nothing
}

func (a *Noop) AddEntryWithData(in *model.FivenetAuditLog, data any) {
	// Nothing
}
