package croner

import (
	"context"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CronjobHandlerFn func(ctx context.Context, data *cron.CronjobData) error

var HandlersModule = fx.Module("cron_handlers",
	fx.Provide(
		NewHandlers,
	),
)

type HandlersParams struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	Handlers []CronHandlersRegister `group:"cronjobhandlers"`
}

type Handlers struct {
	logger *zap.Logger

	mu       sync.Mutex
	handlers map[string]CronjobHandlerFn
}

func NewHandlers(p HandlersParams) (*Handlers, error) {
	h := &Handlers{
		logger: p.Logger.Named("cron.handlers"),

		mu:       sync.Mutex{},
		handlers: map[string]CronjobHandlerFn{},
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		for _, reg := range p.Handlers {
			if err := reg.RegisterCronjobHandlers(h); err != nil {
				return err
			}
		}

		return nil
	}))

	return h, nil
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

type CronHandlersRegister interface {
	RegisterCronjobHandlers(h *Handlers) error
}

// AsCronjobHandlers annotates the given constructor to state that
// it provides a GRPC service to the "cronjobhandlers" group.
func AsCronjobHandlers(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(CronHandlersRegister)),
		fx.ResultTags(`group:"cronjobhandlers"`),
	)
}
