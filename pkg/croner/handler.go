package croner

import (
	"context"
	"sync"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CronjobHandlerFn func(ctx context.Context, data *cron.CronjobData) error

var HandlerModule = fx.Module("cron_handlers",
	fx.Provide(
		NewHandlers,
	),
)

type Handlers struct {
	logger *zap.Logger

	mu       sync.Mutex
	handlers map[string]CronjobHandlerFn
}

func NewHandlers(logger *zap.Logger) *Handlers {
	return &Handlers{
		logger: logger.Named("cron_handlers"),

		mu:       sync.Mutex{},
		handlers: map[string]CronjobHandlerFn{},
	}
}

func (h *Handlers) Add(name string, fn CronjobHandlerFn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	name = events.SanitizeKey(name)

	if _, ok := h.handlers[name]; ok {
		// Getting the stacktrace is expensive but should help tracking down any duplicate cron handlers in no time
		h.logger.Warn("duplicate cron handler override detected", zap.String("name", name), zap.Stack("trace"))
	}

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
