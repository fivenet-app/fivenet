package collab

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var metricTotalConnectedClients = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "collab",
	Name:      "client_count",
	Help:      "Number of connected clients by category.",
}, []string{"category"})

// CollabRoom – one collaborative doc, served by this process
type CollabRoom struct {
	logger   *zap.Logger
	category string

	Id uint64

	// > runtime
	mu      sync.RWMutex
	clients map[uint64]*Client

	// > JetStream handles
	js       jetstream.JetStream
	subject  string
	consumer jetstream.Consumer

	// > lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCollabRoom wires the room to NATS JetStream using the *modern* API.
func NewCollabRoom(ctx context.Context, logger *zap.Logger, roomId uint64, js jetstream.JetStream, category string) (*CollabRoom, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Create consumer
	subject := fmt.Sprintf("collab.%s.%d", category, roomId)

	consumer, err := js.CreateOrUpdateConsumer(ctx, "COLLAB", jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create consumer for room. %w", err)
	}

	room := &CollabRoom{
		logger:   logger.Named("room").With(zap.Uint64("room_id", roomId)),
		category: category,

		Id:       roomId,
		clients:  make(map[uint64]*Client),
		js:       js,
		consumer: consumer,
		subject:  subject,
		ctx:      ctx,
		cancel:   cancel,
	}

	go room.consumeLoop()

	return room, nil
}

// Client management
func (r *CollabRoom) Join(c *Client) {
	r.mu.Lock()
	r.clients[c.Id] = c
	clientCount := len(r.clients)
	r.mu.Unlock()
	metricTotalConnectedClients.WithLabelValues(r.category).Inc()
	r.logger.Debug("client joined", zap.Uint64("client_id", c.Id), zap.Int("clients", clientCount))
}

func (r *CollabRoom) Leave(clientId uint64) bool {
	r.mu.Lock()
	if c, ok := r.clients[clientId]; ok {
		close(c.SendCh)
		delete(r.clients, clientId)
	}
	clientCount := len(r.clients)
	empty := clientCount == 0
	r.mu.Unlock()

	if empty {
		r.shutdown()
	}

	r.logger.Debug("client left", zap.Uint64("client_id", clientId), zap.Int("clients", clientCount))

	metricTotalConnectedClients.WithLabelValues(r.category).Dec()

	return empty
}

// Broadcast – publish to JetStream (everyone, every instance)
func (r *CollabRoom) Broadcast(fromId uint64, msg *collab.ServerPacket) {
	// Ignore "hello" packets that carry no useful data
	if (msg.GetYjsUpdate() != nil && len(msg.GetYjsUpdate().Data) == 0) || (msg.GetAwareness() != nil && len(msg.GetAwareness().Data) == 0) {
		return
	}

	// Encode packet to protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		r.logger.Error("failed to marshal collab message", zap.Error(err))
		return
	}

	// Handle error only for logging
	if _, err := r.js.PublishAsync(r.subject, data); err != nil {
		r.logger.Error("publish error", zap.Error(err))
	}

	// Forward to *local* peers immediately
	r.forwardToLocal(fromId, msg)
}

func (r *CollabRoom) SendToClient(fromId uint64, toId uint64, msg *collab.ServerPacket) {
	// Ignore "hello" packets that carry no useful data
	if (msg.GetYjsUpdate() != nil && len(msg.GetYjsUpdate().Data) == 0) || (msg.GetAwareness() != nil && len(msg.GetAwareness().Data) == 0) {
		return
	}

	r.mu.RLock()
	client, ok := r.clients[toId]
	r.mu.RUnlock()

	if !ok {
		r.logger.Warn("client not found", zap.Uint64("to_id", toId))
		return
	}

	msg.SenderId = fromId // Set sender ID
	client.Send(msg)
}

// JetStream pull loop
func (r *CollabRoom) consumeLoop() {
	for {
		batch, err := r.consumer.Fetch(32,
			jetstream.FetchMaxWait(2*time.Second),
		)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, jetstream.ErrNoMessages) {
				continue // idle
			}
			return
		}

		select {
		case <-r.ctx.Done():
			return

		case msg := <-batch.Messages():
			if err := batch.Error(); err != nil {
				return // Error in batch, stop processing
			}
			if msg == nil {
				continue
			}

			var sp collab.ServerPacket
			if err := proto.Unmarshal(msg.Data(), &sp); err == nil {
				r.forwardToLocal(sp.GetSenderId(), &sp)
			}
			msg.Ack()
		}
	}
}

// forwardToLocal delivers to all *local* clients except the sender.
func (r *CollabRoom) forwardToLocal(fromId uint64, cm *collab.ServerPacket) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for id, c := range r.clients {
		if id == fromId {
			continue // Skip sender
		}

		c.Send(cm)
	}
}

func (r *CollabRoom) BroadcastSyncStep1(fromId uint64, data []byte) {
	pkt := &collab.ServerPacket{
		SenderId: fromId,
		Msg: &collab.ServerPacket_SyncStep{
			SyncStep: &collab.SyncStep{
				Step: 1,
				Data: data,
			},
		},
	}
	r.Broadcast(fromId, pkt)
}

func (r *CollabRoom) ForwardSyncStep2ToClient(fromId uint64, toId uint64, data []byte) {
	pkt := &collab.ServerPacket{
		SenderId: fromId,
		Msg: &collab.ServerPacket_SyncStep{
			SyncStep: &collab.SyncStep{
				Step: 2,
				Data: data,
			},
		},
	}

	r.SendToClient(fromId, toId, pkt)
}

func (r *CollabRoom) BroadcastYjs(fromId uint64, data []byte) {
	pkt := &collab.ServerPacket{
		SenderId: fromId,
		Msg: &collab.ServerPacket_YjsUpdate{
			YjsUpdate: &collab.YjsUpdate{Data: data},
		},
	}
	r.Broadcast(fromId, pkt)
}

func (r *CollabRoom) BroadcastAwareness(fromId uint64, data []byte) {
	pkt := &collab.ServerPacket{
		SenderId: fromId,
		Msg: &collab.ServerPacket_Awareness{
			Awareness: &collab.AwarenessPing{Data: data},
		},
	}
	r.Broadcast(fromId, pkt)
}

func (r *CollabRoom) SendTargetSaved() {
	r.Broadcast(0, &collab.ServerPacket{
		Msg: &collab.ServerPacket_TargetSaved{
			TargetSaved: &collab.TargetSaved{
				TargetId: r.Id,
			},
		},
	})
}

// Shutdown graceful teardown of room
func (r *CollabRoom) shutdown() {
	r.logger.Debug("shutting down")
	r.cancel() // Stop consumeLoop
}
