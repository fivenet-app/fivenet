package centrum

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	centrum "github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/centrummanager"
	eventscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/events"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/coords/postals"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/tracker"
	"github.com/fivenet-app/fivenet/pkg/utils/broker"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	CentrumServiceServer

	logger *zap.Logger
	wg     sync.WaitGroup
	jsCons jetstream.ConsumeContext

	tracer   trace.Tracer
	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	js       *events.JSWrapper
	tracker  tracker.ITracker
	postals  postals.Postals
	appCfg   appconfig.IConfig
	enricher *mstlystcdata.UserAwareEnricher

	brokersMutex sync.RWMutex
	brokers      map[string]*broker.Broker[*StreamResponse]

	state *centrummanager.Manager
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	TP        *tracesdk.TracerProvider
	DB        *sql.DB
	Perms     perms.Permissions
	Audit     audit.IAuditer
	JS        *events.JSWrapper
	Config    *config.Config
	AppConfig appconfig.IConfig
	Tracker   tracker.ITracker
	Postals   postals.Postals
	Manager   *centrummanager.Manager
	Enricher  *mstlystcdata.UserAwareEnricher
}

func NewServer(p Params) (*Server, error) {
	ctxCancel, cancel := context.WithCancel(context.Background())

	brokers := map[string]*broker.Broker[*StreamResponse]{}

	s := &Server{
		logger: p.Logger.Named("centrum"),
		wg:     sync.WaitGroup{},

		tracer: p.TP.Tracer("centrum-cache"),

		db:       p.DB,
		ps:       p.Perms,
		aud:      p.Audit,
		js:       p.JS,
		tracker:  p.Tracker,
		postals:  p.Postals,
		appCfg:   p.AppConfig,
		enricher: p.Enricher,

		brokersMutex: sync.RWMutex{},
		brokers:      brokers,

		state: p.Manager,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		s.handleAppConfigUpdate(ctxCancel, p.AppConfig.Get())

		if err := s.registerSubscriptions(ctxStartup, ctxCancel); err != nil {
			return fmt.Errorf("failed to subscribe to events: %w", err)
		}

		// Handle app config updates
		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-ctxCancel.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					s.handleAppConfigUpdate(ctxCancel, cfg)
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		s.wg.Wait()

		if s.jsCons != nil {
			s.jsCons.Stop()
			s.jsCons = nil
		}

		return nil
	}))

	return s, nil
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterCentrumServiceServer(srv, s)
}

func (s *Server) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	streamCfg, err := eventscentrum.RegisterStream(ctxStartup, s.js)
	if err != nil {
		return fmt.Errorf("failed to register events: %w", err)
	}

	consumer, err := s.js.CreateConsumer(ctxStartup, streamCfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubject: fmt.Sprintf("%s.>", eventscentrum.BaseSubject),
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
		s.jsCons = nil
	}

	s.jsCons, err = consumer.Consume(s.watchForChanges,
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger,
			s.registerSubscriptions,
		))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleAppConfigUpdate(ctx context.Context, cfg *appconfig.Cfg) {
	s.brokersMutex.Lock()
	defer s.brokersMutex.Unlock()

	for _, job := range cfg.UserTracker.LivemapJobs {
		if _, ok := s.brokers[job]; !ok {
			s.brokers[job] = broker.New[*StreamResponse]()
			go s.brokers[job].Start(ctx)
		}
	}
}

func (s *Server) watchForChanges(msg jetstream.Msg) {
	remoteCtx, _ := events.GetJetstreamMsgContext(msg)
	_, span := s.tracer.Start(trace.ContextWithRemoteSpanContext(context.Background(), remoteCtx), msg.Subject())
	defer span.End()

	startTime := time.Now()

	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack message", zap.Error(err))
	}

	job, topic, tType := eventscentrum.SplitSubject(msg.Subject())
	if job == "" || topic == "" || tType == "" {
		if err := msg.TermWithReason("invalid centrum subject"); err != nil {
			s.logger.Error("invalid centrum subject", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	broker, ok := s.getJobBroker(job)
	if !ok {
		s.logger.Debug("no broker found for job", zap.String("job", job))
		return
	}

	resp := &StreamResponse{}
	switch topic {
	case eventscentrum.TopicGeneral:
		switch tType {
		case eventscentrum.TypeGeneralDisponents:
			dest := &centrum.Disponents{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_Disponents{
				Disponents: dest,
			}

		case eventscentrum.TypeGeneralSettings:
			dest := &centrum.Settings{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_Settings{
				Settings: dest,
			}
		}

	case eventscentrum.TopicUnit:
		switch tType {
		case eventscentrum.TypeUnitCreated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitCreated{
				UnitCreated: dest,
			}

		case eventscentrum.TypeUnitDeleted:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitDeleted{
				UnitDeleted: dest,
			}

		case eventscentrum.TypeUnitUpdated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitUpdated{
				UnitUpdated: dest,
			}

		case eventscentrum.TypeUnitStatus:
			dest := &centrum.UnitStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_UnitStatus{
				UnitStatus: dest,
			}
		}

	case eventscentrum.TopicDispatch:
		switch tType {
		case eventscentrum.TypeDispatchCreated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchCreated{
				DispatchCreated: dest,
			}

		case eventscentrum.TypeDispatchDeleted:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchDeleted{
				DispatchDeleted: dest,
			}

		case eventscentrum.TypeDispatchUpdated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchUpdated{
				DispatchUpdated: dest,
			}

		case eventscentrum.TypeDispatchStatus:
			dest := &centrum.DispatchStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &StreamResponse_DispatchStatus{
				DispatchStatus: dest,
			}
		}
	}

	broker.Publish(resp)

	meta, err := msg.Metadata()
	if err != nil {
		s.logger.Error("sent centrum message broker, but failed to get msg metadata ", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
			zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)), zap.Error(err))
		return
	}
	s.logger.Debug("sent centrum message broker", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
		zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)),
		zap.Duration("duration", time.Since(startTime)))
}
