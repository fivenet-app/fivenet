package collab

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client represents a collaborative editing client with its state and communication channels.
type Client struct {
	// logger is the zap logger for this client instance.
	logger *zap.Logger
	// Id is the unique client identifier.
	Id uint64
	// UserId is the user ID associated with this client.
	UserId int32
	// Role is the role of the client in the collaboration session.
	Role collab.ClientRole
	// RoomId is the ID of the room the client is connected to.
	RoomId uint64
	// Stream is the bidirectional gRPC stream for client-server communication.
	Stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]
	// SendCh is a buffered channel for outgoing server packets.
	SendCh chan *collab.ServerPacket
}

// NewClient creates and returns a new Client instance with the provided parameters and a buffered send channel.
func NewClient(logger *zap.Logger, clientId uint64, roomId uint64, UserId int32, role collab.ClientRole, stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]) *Client {
	return &Client{
		logger: logger,
		Id:     clientId,
		Role:   role,
		RoomId: roomId,
		Stream: stream,
		// Buffered channel
		SendCh: make(chan *collab.ServerPacket, 32),
	}
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
