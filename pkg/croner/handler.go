package croner

import (
	"context"
	"sync"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/events"
	"go.uber.org/fx"
)

type CronjobHandlerFn func(ctx context.Context, data *cron.CronjobData) error

var HandlerModule = fx.Module("cron_handlers",
	fx.Provide(
		NewHandlers,
	),
)

type Handlers struct {
	mu       sync.Mutex
	handlers map[string]CronjobHandlerFn
}

func NewHandlers() *Handlers {
	return &Handlers{
		mu:       sync.Mutex{},
		handlers: map[string]CronjobHandlerFn{},
	}
}

func (h *Handlers) Add(name string, fn CronjobHandlerFn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	name = events.SanitizeKey(name)

	h.handlers[name] = fn
}

func (h *Handlers) Remove(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	name = events.SanitizeKey(name)

	delete(h.handlers, name)
}

func (h *Handlers) getCronjobHandler(name string) CronjobHandlerFn {
	h.mu.Lock()
	defer h.mu.Unlock()

	name = events.SanitizeKey(name)

	return h.handlers[name]
}
