package collab

import (
	"context"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client represents a collaborative editing client with its state and communication channels.
type Client struct {
	// logger is the zap logger for this client instance.
	logger *zap.Logger
	// Id is the unique client identifier.
	Id uint64
	// Room is a reference to the CollabRoom this client is connected to.
	room *CollabRoom
	// UserId is the user ID associated with this client.
	UserId int32
	// Role is the role of the client in the collaboration session.
	Role collab.ClientRole
	// Stream is the bidirectional gRPC stream for client-server communication.
	Stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]
	// SendCh is a buffered channel for outgoing server packets.
	SendCh chan *collab.ServerPacket

	presenceKey string
	firstKey    string
	hbCancel    context.CancelFunc
}

// NewClient creates and returns a new Client instance with the provided parameters and a buffered send channel.
func NewClient(logger *zap.Logger, clientId uint64, room *CollabRoom, UserId int32, role collab.ClientRole, stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]) *Client {
	return &Client{
		logger: logger,
		Id:     clientId,
		Role:   role,
		room:   room,
		Stream: stream,
		// Buffered channel
		SendCh: make(chan *collab.ServerPacket, 32),
	}
}

func (c *Client) StartPresence(ctx context.Context) {
	stateKV := c.room.stateKV

	cid := strconv.FormatUint(c.Id, 10)
	roomId := strconv.FormatUint(c.room.Id, 10)
	c.presenceKey = "presence." + roomId + "." + cid
	c.firstKey = "first." + roomId

	// Announce this client
	stateKV.Put(ctx, c.presenceKey, nil)

	// Try to become FIRST (atomic Create)
	if _, err := stateKV.Create(ctx, c.firstKey, []byte(cid), jetstream.KeyTTL(keyTTL)); err == nil {
		c.room.notifyFirst(c.Id) // tell browser it must seed the doc
	}

	// Launch heartbeat + watcher
	hbCtx, cancel := context.WithCancel(ctx)
	c.hbCancel = cancel
	go hbLoop(hbCtx, c.room.stateKV, c.presenceKey, c.firstKey, cid)
	go firstWatch(hbCtx, c.room.stateKV, c.room, c.firstKey, cid, c.Id)
}

func (c *Client) StopPresence() {
	if c.hbCancel != nil {
		c.hbCancel()
	}
	kv := c.room.stateKV
	kv.Delete(context.Background(), c.presenceKey)
	kv.Delete(context.Background(), c.firstKey)
}

// Send attempts to enqueue a message for the client. If the channel is full, the message is dropped and a debug log is emitted.
func (c *Client) Send(msg *collab.ServerPacket) {
	select {
	case c.SendCh <- msg:
		// Message enqueued successfully
	default:
		c.logger.Debug("dropping message for client", zap.Int32("user_id", c.UserId))
	}
}

// SendLoop continuously sends messages from the SendCh channel to the client stream.
// Returns an error if sending fails or the channel is closed.
func (c *Client) SendLoop() error {
	for msg := range c.SendCh {
		if err := c.Stream.Send(msg); err != nil {
			c.logger.Error("error sending to client", zap.Int32("user_id", c.UserId), zap.Error(err))
			return err
		}
	}

	return nil
}

// hbLoop – resets TTL every 2 s
func hbLoop(ctx context.Context, kv jetstream.KeyValue, pKey, fKey, cid string) {
	t := time.NewTicker(hbEvery)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Presence
			kv.Put(ctx, pKey, nil)

			// Safe refresh of “first” (compare-and-swap)
			if e, _ := kv.Get(ctx, fKey); e != nil && string(e.Value()) == cid {
				if _, err := kv.Update(ctx, fKey, []byte(cid), e.Revision()); err != nil {
					// If the update fails, it means another client has taken over as "first"
					continue
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// firstWatch – auto-promotion when “first” key disappears
func firstWatch(ctx context.Context, kv jetstream.KeyValue, room *CollabRoom,
	fKey, cid string, myID uint64,
) {
	sub, _ := kv.Watch(ctx, fKey, jetstream.UpdatesOnly())
	for {
		select {
		case upd := <-sub.Updates():
			if upd == nil {
				continue
			} // catch-up barrier

			if upd.Operation() == jetstream.KeyValueDelete {
				if _, err := kv.Create(ctx, fKey, []byte(cid), jetstream.KeyTTL(keyTTL)); err == nil {
					room.notifyFirst(myID)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
