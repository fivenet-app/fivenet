package client

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/coder/websocket"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/grpcws"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	c *websocket.Conn

	streamId atomic.Uint32
}

func New(ctx context.Context, url string, headers http.Header) (*Client, error) {
	con, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		Subprotocols: []string{"grpc-websocket-channel"},
		HTTPHeader:   headers,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		c:        con,
		streamId: atomic.Uint32{},
	}
	c.streamId.Store(0)

	go c.read(ctx)

	return c, nil
}

func (c *Client) read(ctx context.Context) {
	for {
		mType, r, err := c.c.Reader(ctx)
		if err != nil {
			return
		}

		if mType != websocket.MessageBinary {
			continue
		}

		_ = r
	}
}

func (c *Client) Close() error {
	if c.c == nil {
		return c.c.CloseNow()
	}

	return c.c.Close(websocket.StatusNormalClosure, "")
}

func (c *Client) getStreamId() uint32 {
	return c.streamId.Add(1)
}

func (c *Client) NewChannel(headers map[string]*grpcws.HeaderValue) *WSChannel {
	return &WSChannel{
		c:        c.c,
		streamId: c.getStreamId(),
		headers:  headers,
	}
}

type WSChannel struct {
	c  *websocket.Conn
	in chan grpcws.GrpcFrame

	streamId uint32

	headers map[string]*grpcws.HeaderValue
}

func (w *WSChannel) sendBody(ctx context.Context, msg proto.Message, complete bool) error {
	out, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	body := &grpcws.GrpcFrame{
		StreamId: w.streamId,
		Payload: &grpcws.GrpcFrame_Body{
			Body: &grpcws.Body{
				Data:     w.wrapBodyData(out),
				Complete: complete,
			},
		},
	}

	bodyOut, err := proto.Marshal(body)
	if err != nil {
		return err
	}

	if err := w.c.Write(ctx, websocket.MessageBinary, bodyOut); err != nil {
		return err
	}

	return nil
}

func (w *WSChannel) sendHeaders(ctx context.Context, operation string) error {
	header := &grpcws.GrpcFrame{
		StreamId: 1,
		Payload: &grpcws.GrpcFrame_Header{
			Header: &grpcws.Header{
				Operation: operation,
				Headers:   w.headers,
			},
		},
	}

	out, err := proto.Marshal(header)
	if err != nil {
		return err
	}

	if err := w.c.Write(ctx, websocket.MessageBinary, out); err != nil {
		return err
	}

	return nil
}

func (w *WSChannel) wrapBodyData(out []byte) []byte {
	msgLen := len(out)
	o := make([]byte, 5)
	o[0] = 0
	o[1] = byte(msgLen >> 24)
	o[2] = byte(msgLen >> 16)
	o[3] = byte(msgLen >> 8)
	o[4] = byte(msgLen)

	return append(o, out...)
}
