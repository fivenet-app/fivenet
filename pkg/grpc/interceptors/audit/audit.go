package grpc_audit

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	audit "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logger matches your s.aud.Log signature (adapt as needed).
type Logger interface {
	Log(entry *audit.AuditEntry, req any)
}

// UserExtractor pulls user info (id/job) from ctx (auth middleware, JWT, etc.).
type UserExtractor func(ctx context.Context) (userID uint64, userJob string, ok bool)

// Options configure the interceptors.
type Options struct {
	Logger Logger

	// Override time source (tests)
	Now func() time.Time
	// Last chance to tweak entry
	OnFinalize func(ae *audit.AuditEntry)
	// Record first inbound msg (streams)
	RecordFirst bool
}

// NewUnary returns a unary server interceptor with the given options.
func NewUnary(opts Options) grpc.UnaryServerInterceptor {
	validate(&opts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// Only log authenticated requests
		userInfo, ok := auth.GetUserInfoFromContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		svc, method := splitFullMethod(info.FullMethod)
		start := opts.Now()

		ae := &audit.AuditEntry{
			Service: svc,
			Method:  method,
			Action:  audit.EventAction_EVENT_ACTION_UNSPECIFIED,
			UserId:  userInfo.GetUserId(),
			UserJob: userInfo.GetJob(),
			Meta:    &audit.AuditEntryMeta{},
		}

		// Store Entry handle in context so handlers can mutate it.
		handle := &Entry{entry: ae}
		ctx = withEntry(ctx, handle)

		// Panic safety and finalize logging.
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
				handle.set(func(a *audit.AuditEntry) {
					a.Result = audit.EventResult_EVENT_RESULT_ERRORED
				})
			}

			code := status.Code(err)
			duration := opts.Now().Sub(start)
			handle.set(func(a *audit.AuditEntry) {
				a.Meta.Set("grpc_code", code.String())
				a.Meta.Set("duration_ms", strconv.FormatInt(int64(duration/time.Millisecond), 10))

				if err != nil && a.Result != audit.EventResult_EVENT_RESULT_SUCCEEDED {
					a.Result = audit.EventResult_EVENT_RESULT_ERRORED
				}
				if a.Result == audit.EventResult_EVENT_RESULT_UNSPECIFIED {
					a.Result = audit.EventResult_EVENT_RESULT_SUCCEEDED
				}
				if a.Action == audit.EventAction_EVENT_ACTION_UNSPECIFIED {
					a.Action = audit.EventAction_EVENT_ACTION_VIEWED
				}
			})

			if opts.OnFinalize != nil {
				opts.OnFinalize(handle.entry)
			}

			if handle.IsSkipped() {
				return
			}

			opts.Logger.Log(handle.entry, any(req))
		}()

		return handler(ctx, req)
	}
}

// NewStream returns a stream server interceptor with the given options.
func NewStream(opts Options) grpc.StreamServerInterceptor {
	validate(&opts)
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := ss.Context()
		userInfo, ok := auth.GetUserInfoFromContext(ctx)
		if !ok {
			return handler(srv, ss)
		}

		svc, method := splitFullMethod(info.FullMethod)
		start := opts.Now()

		ae := &audit.AuditEntry{
			Service: svc,
			Method:  method,
			Action:  audit.EventAction_EVENT_ACTION_UNSPECIFIED,
			UserId:  userInfo.GetUserId(),
			UserJob: userInfo.GetJob(),
			Meta:    &audit.AuditEntryMeta{},
		}

		handle := &Entry{entry: ae}
		wrapped := &auditStream{
			ServerStream: ss,
			ctx:          withEntry(ss.Context(), handle),
			opts:         opts,
			firstInOnce:  sync.Once{},
			firstIn:      nil,
		}

		if userInfo, ok := auth.GetUserInfoFromContext(wrapped.ctx); ok {
			ae.UserId = userInfo.GetUserId()
			ae.UserJob = userInfo.GetJob()
		}

		defer func() {
			if r := recover(); r != nil {
				handle.set(func(a *audit.AuditEntry) {
					a.Result = audit.EventResult_EVENT_RESULT_ERRORED
				})
			}

			code := status.Code(err)
			duration := opts.Now().Sub(start)
			handle.set(func(a *audit.AuditEntry) {
				if a.Meta == nil {
					a.Meta = &audit.AuditEntryMeta{
						Meta: make(map[string]string),
					}
				}
				a.Meta.Set("grpc_code", code.String())
				a.Meta.Set("duration_ms", strconv.FormatInt(int64(duration/time.Millisecond), 10))
				a.Meta.Set("stream_recv_count", strconv.FormatInt(wrapped.recvCount, 10))
				a.Meta.Set("stream_send_count", strconv.FormatInt(wrapped.sendCount, 10))

				if err != nil && a.Result != audit.EventResult_EVENT_RESULT_SUCCEEDED {
					a.Result = audit.EventResult_EVENT_RESULT_ERRORED
				}
				if a.Result == audit.EventResult_EVENT_RESULT_UNSPECIFIED {
					a.Result = audit.EventResult_EVENT_RESULT_SUCCEEDED
				}
			})

			if opts.OnFinalize != nil {
				opts.OnFinalize(handle.entry)
			}

			if handle.IsSkipped() {
				return
			}

			// For streams we often log the first inbound message (if enabled),
			// otherwise nil to avoid giant logs.
			logReq := any(nil)
			if opts.RecordFirst {
				logReq = wrapped.firstIn
			}
			opts.Logger.Log(handle.entry, logReq)
		}()

		return handler(srv, wrapped)
	}
}

// Entry is a threadsafe handle stored in context for handlers to mutate.
type Entry struct {
	mu    sync.Mutex
	entry *audit.AuditEntry
	skip  bool
}

type ctxKey struct{}

func withEntry(ctx context.Context, e *Entry) context.Context {
	return context.WithValue(ctx, ctxKey{}, e)
}

// FromContext returns the audit Entry handle (nil if none).
func FromContext(ctx context.Context) *Entry {
	v, _ := ctx.Value(ctxKey{}).(*Entry)
	return v
}

// set applies a mutation to the underlying AuditEntry with locking.
func (e *Entry) set(fn func(*audit.AuditEntry)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	fn(e.entry)
}

func (e *Entry) setE(fn func(*Entry)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	fn(e)
}

func (e *Entry) IsSkipped() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.skip
}

func Skip(ctx context.Context) {
	if e := FromContext(ctx); e != nil {
		e.setE(func(a *Entry) { a.skip = true })
	}
}

func IsSkip(ctx context.Context) bool {
	if e := FromContext(ctx); e != nil {
		return e.IsSkipped()
	}
	return false
}

func SetAction(ctx context.Context, ac audit.EventAction) {
	if e := FromContext(ctx); e != nil {
		e.set(func(a *audit.AuditEntry) { a.Action = ac })
	}
}

func SetResult(ctx context.Context, st audit.EventResult) {
	if e := FromContext(ctx); e != nil {
		e.set(func(a *audit.AuditEntry) { a.Result = st })
	}
}

func SetUser(ctx context.Context, userId int32, job string) {
	if e := FromContext(ctx); e != nil {
		e.set(func(a *audit.AuditEntry) {
			a.UserId = userId
			a.UserJob = job
		})
	}
}

func SetTargetUser(ctx context.Context, userId int32, job string) {
	if e := FromContext(ctx); e != nil {
		e.set(func(a *audit.AuditEntry) {
			a.TargetUserId = &userId
			if job != "" {
				a.TargetUserJob = &job
			}
		})
	}
}

func AddMeta(ctx context.Context, key, val string) {
	if e := FromContext(ctx); e != nil {
		e.set(func(a *audit.AuditEntry) {
			// ensure a.Meta is a map<string,string> in your proto
			if a.Meta == nil {
				a.Meta = &audit.AuditEntryMeta{
					Meta: make(map[string]string),
				}
			}
			a.Meta.Meta[key] = val
		})
	}
}

func splitFullMethod(full string) (service, method string) {
	i := strings.LastIndex(full, "/")
	if i < 0 || i+1 >= len(full) {
		return "", strings.TrimPrefix(full, "/")
	}
	return full[1:i], full[i+1:]
}

func validate(o *Options) {
	if o.Logger == nil {
		panic("auditrpc: Options.Logger is required")
	}
	if o.Now == nil {
		o.Now = time.Now
	}
}

// auditStream wraps ServerStream to count messages and capture the first inbound.
type auditStream struct {
	grpc.ServerStream

	ctx         context.Context
	opts        Options
	recvCount   int64
	sendCount   int64
	firstInOnce sync.Once
	firstIn     any
}

func (w *auditStream) Context() context.Context { return w.ctx }

func (w *auditStream) RecvMsg(m any) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		w.recvCount++
		if w.opts.RecordFirst {
			w.firstInOnce.Do(func() {
				// shallow capture
				w.firstIn = m
			})
		}
	}
	return err
}

func (w *auditStream) SendMsg(m any) error {
	err := w.ServerStream.SendMsg(m)
	if err == nil {
		w.sendCount++
	}
	return err
}
