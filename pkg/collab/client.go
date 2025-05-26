package collab

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Client struct {
	logger *zap.Logger
	Id     uint64
	UserId int32
	Role   collab.ClientRole
	RoomId uint64
	Stream grpc.BidiStreamingServer[collab.ClientPacket, collab.ServerPacket]
	SendCh chan *collab.ServerPacket
}

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

func (c *Client) Send(msg *collab.ServerPacket) {
	select {
	case c.SendCh <- msg:

	default:
		c.logger.Debug("dropping message for client", zap.Int32("user_id", c.UserId))
	}
}

func (c *Client) SendLoop() error {
	for msg := range c.SendCh {
		if err := c.Stream.Send(msg); err != nil {
			c.logger.Error("error sending to client", zap.Int32("user_id", c.UserId), zap.Error(err))
			return err
		}
	}

	return nil
}
