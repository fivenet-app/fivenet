package grpcws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/desertbit/timer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/grpcws"
	"google.golang.org/protobuf/proto"
)

type WebsocketChannel struct {
	mu             sync.Mutex
	ctx            context.Context
	wsConn         *websocket.Conn
	grpcHandler    func(resp http.ResponseWriter, req *http.Request)
	maxStreamCount int
	req            *http.Request

	activeStreams   map[uint32]*GrpcStream
	timeoutInterval time.Duration
	timer           *timer.Timer
}

var framePingResponse = &grpcws.GrpcFrame{
	StreamId: 0,
	Payload: &grpcws.GrpcFrame_Ping{
		Ping: &grpcws.Ping{
			Pong: true,
		},
	},
}

func NewWebsocketChannel(
	ctx context.Context,
	websocket *websocket.Conn,
	grpcHandler func(resp http.ResponseWriter, req *http.Request),
	maxStreamCount int,
	req *http.Request,
) *WebsocketChannel {
	return &WebsocketChannel{
		mu:             sync.Mutex{},
		ctx:            ctx,
		wsConn:         websocket,
		grpcHandler:    grpcHandler,
		maxStreamCount: maxStreamCount,
		req:            req,

		activeStreams:   make(map[uint32]*GrpcStream, maxStreamCount),
		timeoutInterval: 12 * time.Second,
		timer:           nil,
	}
}

func (ws *WebsocketChannel) Start() {
	for {
		if err := ws.poll(); err != nil {
			return
		}
	}
}

func (ws *WebsocketChannel) writeError(streamId uint32, message string) error {
	return ws.write(
		&grpcws.GrpcFrame{
			StreamId: streamId,
			Payload: &grpcws.GrpcFrame_Failure{
				Failure: &grpcws.Failure{
					ErrorMessage: message,
				},
			},
		},
	)
}

func (ws *WebsocketChannel) getStream(streamId uint32) *GrpcStream {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.activeStreams[streamId]
}

func (ws *WebsocketChannel) deleteStream(streamId uint32) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	delete(ws.activeStreams, streamId)
}

func (ws *WebsocketChannel) poll() error {
	frame, err := ws.readFrame()
	if errors.Is(err, io.EOF) {
		ws.Close()
	}

	if err != nil {
		return err
	}

	stream := ws.getStream(frame.GetStreamId())

	switch payload := frame.GetPayload().(type) {
	case *grpcws.GrpcFrame_Ping:
		return ws.write(&grpcws.GrpcFrame{
			StreamId: frame.GetStreamId(),
			Payload: &grpcws.GrpcFrame_Ping{
				Ping: &grpcws.Ping{
					Pong: true,
				},
			},
		})

	case *grpcws.GrpcFrame_Header:
		if stream != nil {
			return ws.writeError(frame.GetStreamId(), "stream already exists")
		}

		stream, err := func() (*GrpcStream, error) {
			ws.mu.Lock()
			defer ws.mu.Unlock()

			if ws.maxStreamCount > 0 && len(ws.activeStreams) > ws.maxStreamCount {
				return nil, ws.writeError(frame.GetStreamId(), "rejecting max number of streams reached for this channel")
			}

			stream := newGrpcStream(frame.GetStreamId(), ws, ws.maxStreamCount)
			ws.activeStreams[frame.GetStreamId()] = stream

			return stream, nil
		}()
		if err != nil {
			return err
		}

		url, err := url.Parse("http://localhost/")
		if err != nil {
			return ws.writeError(frame.GetStreamId(), err.Error())
		}
		url.Scheme = ws.req.URL.Scheme
		url.Host = ws.req.URL.Host
		url.Path = "/" + frame.GetHeader().GetOperation()

		req := &http.Request{
			Method:     http.MethodPost,
			URL:        url,
			Header:     ws.req.Header.Clone(),
			Body:       stream,
			RemoteAddr: ws.req.RemoteAddr,
		}
		for key, element := range frame.GetHeader().GetHeaders() {
			req.Header[key] = element.GetValue()
		}

		interceptedReq := makeGrpcRequest(req.WithContext(stream.ctx))
		// Forward the request to the grpcHandler
		go func() {
			defer ws.deleteStream(stream.id)

			ws.grpcHandler(stream, interceptedReq)
		}()

	case *grpcws.GrpcFrame_Body:
		if stream == nil || stream.inputClosed {
			return ws.writeError(frame.GetStreamId(), "stream does not exist")
		}

		stream.inputFrames <- frame

		// Close channel if body frame says so
		if payload.Body.GetComplete() && !stream.inputClosed {
			stream.inputClosed = true
			close(stream.inputFrames)
		}

	case *grpcws.GrpcFrame_Cancel:
		if stream == nil {
			// If a stream is being cancelled and it's not there anymore, don't return an error
			return nil
		}

		stream.cancel()
		if !stream.inputClosed {
			stream.inputClosed = true
			close(stream.inputFrames)
		}
		ws.deleteStream(frame.GetStreamId())

	case *grpcws.GrpcFrame_Complete:
		if stream == nil {
			return ws.writeError(frame.GetStreamId(), "stream does not exist")
		}

		if !stream.inputClosed {
			stream.inputClosed = true
			close(stream.inputFrames)
		}

	case *grpcws.GrpcFrame_Failure:
		if stream == nil {
			return ws.writeError(frame.GetStreamId(), "stream does not exist")
		}

		stream.inputFrames <- frame
		if !stream.inputClosed {
			stream.inputClosed = true
			close(stream.inputFrames)
		}
		ws.deleteStream(frame.GetStreamId())

	default:
		return ws.writeError(frame.GetStreamId(), "unknown frame type")
	}

	return nil
}

func (ws *WebsocketChannel) readFrame() (*grpcws.GrpcFrame, error) {
	// we assume a large limit is set for the websocket to avoid handling multiple frames.
	msgType, bytesValue, err := ws.wsConn.Read(ws.ctx)
	if err != nil {
		return nil, err
	}

	if msgType != websocket.MessageBinary {
		return nil, errors.New("websocket channel only supports binary messages")
	}

	request := &grpcws.GrpcFrame{}
	if err := proto.Unmarshal(bytesValue, request); err != nil {
		return nil, fmt.Errorf("frame unmarshal error. %w", err)
	}

	return request, nil
}

func (ws *WebsocketChannel) write(frame *grpcws.GrpcFrame) error {
	binaryFrame, err := proto.Marshal(frame)
	if err != nil {
		return err
	}
	return ws.wsConn.Write(ws.ctx, websocket.MessageBinary, binaryFrame)
}

func (w *WebsocketChannel) enablePing(timeOutInterval time.Duration) {
	w.timeoutInterval = timeOutInterval
	w.timer = timer.NewTimer(w.timeoutInterval)
	go w.ping()
}

func (w *WebsocketChannel) ping() {
	defer w.timer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-w.timer.C:
			w.timer.Reset(w.timeoutInterval)

			stream := w.getStream(0)
			if stream == nil {
				return
			}

			stream.channel.write(framePingResponse)
		}
	}
}

func (ws *WebsocketChannel) Close() {}
