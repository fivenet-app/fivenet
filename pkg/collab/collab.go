package collab

import (
	"context"
	"errors"
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

// metricTotalCollabRooms tracks the number of active collaborative rooms for Prometheus monitoring.
var metricTotalCollabRooms = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "collab",
	Name:      "room_count",
	Help:      "Number of active collaborative rooms.",
}, []string{"category"})

const (
	kvBucket = "collab_state"
	// key expiration TTL if not touched.
	keyTTL = 3 * time.Second
	// heartbeat period (< keyTTL).
	hbEvery = time.Duration(1.5 * float64(time.Second))
)

// CollabServer manages collaborative editing rooms and client connections.
type CollabServer struct {
	// ctx is the base context for the server and rooms.
	ctx context.Context
	// logger is the zap logger for this server instance.
	logger *zap.Logger
	// js is the JetStream wrapper for event streaming.
	js *events.JSWrapper

	// category is the logical category for this server (used for stream subjects and metrics).
	category string

	// mu protects the rooms map.
	mu sync.Mutex
	// rooms maps target IDs to active CollabRoom instances.
	rooms map[int64]*CollabRoom

	// Key-Value store for state management
	stateKV jetstream.KeyValue
}

// New creates and returns a new CollabServer for the given category, logger, and JetStream wrapper.
// Panics if the category is empty.
func New(
	ctx context.Context,
	logger *zap.Logger,
	js *events.JSWrapper,
	category string,
) *CollabServer {
	if category == "" {
		panic("collab category must not be empty")
	}

	return &CollabServer{
		ctx:    ctx,
		logger: logger.Named("collab.server").With(zap.String("category", category)),
		js:     js,

		category: category,

		mu:    sync.Mutex{},
		rooms: make(map[int64]*CollabRoom),
	}
}

// Start creates or updates the JetStream stream for collaborative editing events for this server's category.
func (s *CollabServer) Start(ctx context.Context) error {
	stateKV, err := s.js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:         kvBucket,
		Description:    "Collab state Store",
		TTL:            keyTTL,
		History:        1,
		LimitMarkerTTL: 2 * keyTTL,
		Storage:        jetstream.MemoryStorage,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update state key-value store. %w", err)
	}
	s.stateKV = stateKV

	cfg := jetstream.StreamConfig{
		Name:        "COLLAB",
		Description: "Collaborative editing events stream",
		Subjects:    []string{fmt.Sprintf("collab.%s.*", s.category)},
		Storage:     jetstream.MemoryStorage,
		Retention:   jetstream.InterestPolicy,
		MaxAge:      5 * time.Minute,
		Discard:     jetstream.DiscardOld,
		Duplicates:  time.Minute,
	}

	if _, err := s.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return fmt.Errorf("failed to create/update stream. %w", err)
	}

	if err := s.js.DeleteKeyValue(ctx, "COLLAB_STATE"); err != nil &&
		!errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to delete old collab state key-value store. %w", err)
	}

	return nil
}

// HandleFirstMsg waits for the first message from a client to determine the target ID.
// Returns the target ID or an error if the message is invalid.
func (s *CollabServer) HandleFirstMsg(
	_ context.Context,
	_ uint64,
	stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket],
) (int64, error) {
	// Wait for the first message to determine client/target id
	firstMsg, err := stream.Recv()
	if err != nil {
		return 0, err
	}
	hello := firstMsg.GetHello()
	if hello == nil {
		return 0, status.Error(codes.InvalidArgument, "first message must be CollabInit")
	}
	if hello.GetTargetId() == 0 {
		return 0, status.Error(codes.InvalidArgument, "zero target id provided in first message")
	}

	return hello.GetTargetId(), nil
}

// sendHelloResponse sends a handshake response to the client with its client ID and whether it is the first in the room.
func (s *CollabServer) sendHelloResponse(
	clientId uint64,
	first bool,
	stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket],
) error {
	if err := stream.Send(&collab.ServerPacket{
		SenderId: clientId,
		Msg: &collab.ServerPacket_Handshake{
			Handshake: &collab.CollabHandshake{
				ClientId: clientId,
				First:    first,
			},
		},
	}); err != nil {
		return status.Error(codes.Internal, fmt.Sprintf("failed to send client ID. %v", err))
	}

	return nil
}

// getOrCreateRoom retrieves an existing CollabRoom for the target or creates a new one if needed.
// Returns the room, a boolean indicating if it was created, and any error.
func (s *CollabServer) getOrCreateRoom(targetId int64) (*CollabRoom, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	// Get or create the document room
	room, exists := s.rooms[targetId]
	if !exists {
		room, err = NewCollabRoom(s.ctx, s.logger, s.stateKV, targetId, s.js.JetStream, s.category)
		if err != nil {
			return nil, false, err
		}
		s.rooms[targetId] = room
		metricTotalCollabRooms.WithLabelValues(s.category).Inc()

		return room, true, err
	}

	return room, false, nil
}

// HandleClient manages the lifecycle and message handling for a single collaborative client connection.
// It joins the client to the room, handles sending and receiving, and cleans up when the client leaves.
func (s *CollabServer) HandleClient(
	ctx context.Context,
	targetId int64,
	userId int32,
	clientId uint64,
	role collab.ClientRole,
	stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket],
) error {
	room, created, err := s.getOrCreateRoom(targetId)
	if err != nil {
		return fmt.Errorf("get or create room. %w", err)
	}

	if err := s.sendHelloResponse(clientId, created, stream); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	client := NewClient(s.logger.Named("client"), clientId, room, userId, role, stream)
	if err := room.Join(ctx, client); err != nil {
		return fmt.Errorf("failed to join room for client %d. %w", clientId, err)
	}

	defer func() {
		// If the room is empty after the client leaves, remove it
		if room.Leave(s.ctx, clientId) {
			// Room now has zero clients
			s.mu.Lock()
			delete(s.rooms, targetId)
			s.mu.Unlock()
			metricTotalCollabRooms.WithLabelValues(s.category).Dec()
		}
	}()

	// Handle sending
	go client.SendLoop()

	// Handle receiving
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			return err
		}

		switch m := msg.GetMsg().(type) {
		case *collab.ClientPacket_SyncStep:
			switch m.SyncStep.GetStep() {
			case 1:
				if m.SyncStep.ReceiverId != nil {
					return status.Error(codes.InvalidArgument, "sync step 1 must not have a receiver ID")
				}
				room.BroadcastSyncStep1(clientId, m.SyncStep.GetData())

			case 2:
				if m.SyncStep.ReceiverId == nil {
					return status.Error(codes.InvalidArgument, "sync step 2 must have a receiver ID")
				}

				room.ForwardSyncStep2ToClient(clientId, m.SyncStep.GetReceiverId(), m.SyncStep.GetData())

			default:
				return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid sync step: %d", m.SyncStep.GetStep()))
			}

		case *collab.ClientPacket_YjsUpdate:
			if client.Role < collab.ClientRole_CLIENT_ROLE_WRITER {
				// Reader client updates are ignored
				continue
			}

			room.BroadcastYjs(clientId, m.YjsUpdate.GetData())

		case *collab.ClientPacket_Awareness:
			room.BroadcastAwareness(clientId, m.Awareness.GetData())

		default:
			// ignore unknown packet types
		}
	}
}

// SendTargetSaved notifies the room for the given targetId that the target has been saved.
func (s *CollabServer) SendTargetSaved(_ context.Context, targetId int64) {
	s.mu.Lock()
	room, exists := s.rooms[targetId]
	s.mu.Unlock()

	if !exists {
		return // No room exists for this target, nothing to do
	}

	room.SendTargetSaved()
}
