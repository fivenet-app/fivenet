package grpcws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"
	"github.com/desertbit/timer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/grpcws"
	"google.golang.org/protobuf/proto"
)

type WebsocketChannel struct {
	ctx            context.Context
	wsConn         *websocket.Conn
	handler        http.Handler
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

func NewWebsocketChannel(ctx context.Context, websocket *websocket.Conn, handler http.Handler, maxStreamCount int, req *http.Request) *WebsocketChannel {
	return &WebsocketChannel{
		ctx:            ctx,
		wsConn:         websocket,
		handler:        handler,
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
	return ws.activeStreams[streamId]
}

func (ws *WebsocketChannel) deleteStream(streamId uint32) {
	delete(ws.activeStreams, streamId)
}

func (ws *WebsocketChannel) poll() error {
	frame, err := ws.readFrame()
	if err == io.EOF {
		ws.Close()
	}

	if err != nil {
		return err
	}

	stream := ws.getStream(frame.StreamId)

	switch frame.Payload.(type) {
	case *grpcws.GrpcFrame_Header:
		// grpclog.Infof("received Header for stream %v %v", frame.StreamId, frame.GetHeader().Operation)
		if stream != nil {
			ws.writeError(frame.StreamId, "stream already exists")
		}

		if ws.maxStreamCount > 0 && len(ws.activeStreams) > ws.maxStreamCount {
			return ws.writeError(frame.StreamId, "rejecting max number of streams reached for this channel")
		}

		stream := newGrpcStream(frame.StreamId, ws, ws.maxStreamCount)
		ws.activeStreams[frame.StreamId] = stream

		url, err := url.Parse("http://localhost/")
		if err != nil {
			return ws.writeError(frame.StreamId, err.Error())
		}
		url.Scheme = ws.req.URL.Scheme
		url.Host = ws.req.URL.Host
		url.Path = "/" + frame.GetHeader().Operation

		req := &http.Request{
			Method:     http.MethodPost,
			URL:        url,
			Header:     ws.req.Header.Clone(),
			Body:       stream,
			RemoteAddr: ws.req.RemoteAddr,
		}
		for key, element := range frame.GetHeader().Headers {
			req.Header[key] = element.Value
		}
		// grpclog.Infof("starting grpc request to %v", url)

		// TODO add handler to the websocket channel and then forward it to this.
		interceptedRequest := makeGrpcRequest(req.WithContext(stream.ctx))
		// grpclog.Infof("starting call to http server %q", interceptedRequest.Method)
		go ws.handler.ServeHTTP(stream, interceptedRequest)

	case *grpcws.GrpcFrame_Body:
		// grpclog.Infof("received Body for stream %v", frame.StreamId)
		if stream == nil {
			return ws.writeError(frame.StreamId, "stream does not exist")
		}

		// grpclog.Infof("received body %v", frame)
		stream.inputFrames <- frame

		// Close channel if body frame says so
		body := frame.Payload.(*grpcws.GrpcFrame_Body)
		if body.Body.Complete {
			stream.close()
		}

	case *grpcws.GrpcFrame_Cancel:
		// grpclog.Infof("received cancel for stream %v", frame.StreamId)
		if stream == nil {
			// If a stream is being cancelled and it's not there anymore, don't return an error
			return nil
		}

		// grpclog.Infof("stream %v is canceled", frame.StreamId)
		stream.cancel()
		stream.close()
		ws.deleteStream(frame.StreamId)

	case *grpcws.GrpcFrame_Complete:
		// grpclog.Infof("received complete for stream %v", frame.StreamId)
		if stream == nil {
			return ws.writeError(frame.StreamId, "stream does not exist")
		}

		// grpclog.Infof("completing input stream %v", frame.StreamId)
		stream.close()

	case *grpcws.GrpcFrame_Failure:
		// grpclog.Infof("received Failure for stream %v", frame.StreamId)
		if stream == nil {
			return ws.writeError(frame.StreamId, "stream does not exist")
		}

		// grpclog.Infof("error on stream %v: %v", frame.StreamId, frame.GetFailure().ErrorMessage)
		stream.inputFrames <- frame
		ws.deleteStream(frame.StreamId)

	default:
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
