package collab

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var metricTotalCollabRooms = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "collab",
	Name:      "room_count",
	Help:      "Number of active collaborative rooms.",
}, []string{"category"})

type CollabServer struct {
	ctx    context.Context
	logger *zap.Logger
	js     *events.JSWrapper

	category string

	mu    sync.Mutex
	rooms map[uint64]*CollabRoom
}

func New(ctx context.Context, logger *zap.Logger, js *events.JSWrapper, category string) *CollabServer {
	if category == "" {
		panic("collab category must not be empty")
	}

	return &CollabServer{
		ctx:    ctx,
		logger: logger.Named("collab_server").With(zap.String("category", category)),
		js:     js,

		category: category,

		mu:    sync.Mutex{},
		rooms: make(map[uint64]*CollabRoom),
	}
}

func (s *CollabServer) Start(ctx context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:       "COLLAB",
		Subjects:   []string{fmt.Sprintf("collab.%s.*", s.category)},
		Storage:    jetstream.MemoryStorage,
		Retention:  jetstream.LimitsPolicy,
		MaxAge:     24 * time.Hour,
		Discard:    jetstream.DiscardOld,
		Duplicates: time.Minute,
	}

	if _, err := s.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create/update stream. %w", err)
	}

	return nil
}

func (s *CollabServer) HandleFirstMsg(ctx context.Context, clientId uint64, stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]) (uint64, error) {
	// Wait for the first message to determine client/target id
	firstMsg, err := stream.Recv()
	if err != nil {
		return 0, err
	}
	hello := firstMsg.GetHello()
	if hello == nil {
		return 0, status.Error(codes.InvalidArgument, "first message must be CollabHello")
	}
	if hello.TargetId == 0 {
		return 0, status.Error(codes.InvalidArgument, "zero target Id provided in first message")
	}

	if err := s.SendHelloResponse(clientId, stream); err != nil {
		return 0, status.Error(codes.Internal, err.Error())
	}

	return hello.TargetId, nil
}

func (s *CollabServer) SendHelloResponse(clientId uint64, stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]) error {
	if err := stream.Send(&collab.ServerPacket{
		SenderId: clientId,
		Msg: &collab.ServerPacket_ClientId{
			ClientId: clientId,
		},
	}); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to send client ID. %v", err))
	}

	return nil
}

func (s *CollabServer) getOrCreateRoom(targetId uint64) (*CollabRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	// Get or create the document room
	room, exists := s.rooms[targetId]
	if !exists {
		room, err = NewCollabRoom(s.ctx, s.logger, targetId, s.js.JetStream, s.category)
		if err != nil {
			return nil, err
		}
		s.rooms[targetId] = room
		metricTotalCollabRooms.WithLabelValues(s.category).Inc()
	}

	return room, nil
}

func (s *CollabServer) HandleClient(ctx context.Context, targetId uint64, userId int32, clientId uint64, role collab.ClientRole, stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]) error {
	room, err := s.getOrCreateRoom(targetId)
	if err != nil {
		return fmt.Errorf("get or create room: %w", err)
	}

	client := NewClient(s.logger.Named("client"), clientId, targetId, userId, role, stream)
	room.Join(client)
	defer func() {
		// If the room is empty after the client leaves, remove it
		if room.Leave(clientId) {
			// Room now has zero clients
			s.mu.Lock()
			delete(s.rooms, targetId)
			s.mu.Unlock()
			metricTotalCollabRooms.WithLabelValues(s.category).Dec()
		} else if aw := encodeAwarenessRemove(clientId); len(aw) > 0 {
			// If not, send (valid) leave message to all clients left in the room
			leave := &collab.ServerPacket{
				SenderId: clientId,
				Msg: &collab.ServerPacket_Awareness{
					Awareness: &collab.AwarenessPing{
						Data: aw,
					},
				},
			}
			// TODO in frontend yjs code: `Uncaught (in promise) RangeError: attempting to construct out-of-bounds Uint8Array on ArrayBuffer``
			room.Broadcast(clientId, leave)
		}
	}()

	// Handle sending
	go client.SendLoop()

	// Handle receiving
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch m := msg.Msg.(type) {
		case *collab.ClientPacket_YjsUpdate:
			if client.Role < collab.ClientRole_CLIENT_ROLE_WRITER {
				// Reader client updates are ignored
				continue
			}
			room.BroadcastYjs(clientId, m.YjsUpdate.Data)

		case *collab.ClientPacket_Awareness:
			room.BroadcastAwareness(clientId, m.Awareness.Data)

		default:
			// ignore unknown packet types
		}
	}
}
