package grpcws

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/desertbit/timer"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/grpcws"
	"google.golang.org/protobuf/proto"
)

type WebsocketChannel struct {
	mutex sync.RWMutex

	wsConn          *websocket.Conn
	activeStreams   map[uint32]*GrpcStream
	timeOutInterval time.Duration
	timer           *timer.Timer
	context         context.Context
	handler         http.Handler
	maxStreamCount  int
	req             *http.Request
}

type GrpcStream struct {
	mutex sync.Mutex

	id                uint32
	hasWrittenHeaders bool
	responseHeaders   http.Header
	inputFrames       chan *grpcws.GrpcFrame
	channel           *WebsocketChannel
	ctx               context.Context
	cancel            context.CancelFunc
	remainingBuffer   []byte
	remainingError    error
	readHeader        bool
	bytesToWrite      uint32
	writeBuffer       []byte
	closed            bool
	// TODO add a context to return to close the connection
}

var framePingResponse = &grpcws.GrpcFrame{
	StreamId: 0,
	Payload: &grpcws.GrpcFrame_Ping{
		Ping: &grpcws.Ping{
			Pong: true,
		},
	},
}

func NewWebsocketChannel(websocket *websocket.Conn, handler http.Handler, ctx context.Context, maxStreamCount int, req *http.Request) *WebsocketChannel {
	wsCh := &WebsocketChannel{
		mutex: sync.RWMutex{},

		wsConn: websocket,
		activeStreams: make(
			map[uint32]*GrpcStream),
		timeOutInterval: 12 * time.Second,
		timer:           nil,
		context:         ctx,
		handler:         handler,
		maxStreamCount:  maxStreamCount,
		req:             req,
	}
	wsCh.activeStreams[0] = newGrpcStream(0, wsCh, maxStreamCount)

	return wsCh
}

func newGrpcStream(streamId uint32, channel *WebsocketChannel, streamBufferSize int) *GrpcStream {
	ctx, cancel := context.WithCancel(channel.context)
	return &GrpcStream{
		mutex: sync.Mutex{},

		id:              streamId,
		inputFrames:     make(chan *grpcws.GrpcFrame, streamBufferSize),
		channel:         channel,
		ctx:             ctx,
		cancel:          cancel,
		responseHeaders: make(http.Header),
	}
}

func (stream *GrpcStream) Flush() {}

func (stream *GrpcStream) Header() http.Header {
	return stream.responseHeaders
}

func (stream *GrpcStream) close() {
	if stream.closed {
		return
	}

	stream.closed = true
	close(stream.inputFrames)
}

func (stream *GrpcStream) Read(p []byte) (int, error) {
	// grpclog.Infof("reading from channel %v", stream.id)
	if stream.remainingBuffer != nil {
		// If the remaining buffer fits completely inside the argument slice then read all of it and return any error
		// that was retained from the original call
		if len(stream.remainingBuffer) <= len(p) {
			copy(p, stream.remainingBuffer)

			remainingLength := len(stream.remainingBuffer)
			err := stream.remainingError

			// Clear the remaining buffer and error so that the next read will be a read from the websocket frame,
			// unless the error terminates the stream
			stream.remainingBuffer = nil
			stream.remainingError = nil
			return remainingLength, err
		}

		// The remaining buffer doesn't fit inside the argument slice, so copy the bytes that will fit and retain the
		// bytes that don't fit - don't return the remainingError as there are still bytes to be read from the frame
		copy(p, stream.remainingBuffer[:len(p)])
		stream.remainingBuffer = stream.remainingBuffer[len(p):]

		// Return the length of the argument slice as that was the length of the written bytes
		return len(p), nil
	}

	frame, more := <-stream.inputFrames
	// grpclog.Infof("received message %v more: %v", frame, more)
	if more {
		switch op := frame.Payload.(type) {
		case *grpcws.GrpcFrame_Body:
			stream.remainingBuffer = op.Body.Data
			return stream.Read(p)
		case *grpcws.GrpcFrame_Failure:
			// TODO how to propagate this to the server?
			return 0, errors.New("grpc client error")
		}
	}
	return 0, io.EOF
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
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	return ws.activeStreams[streamId]
}

func (ws *WebsocketChannel) deleteStream(streamId uint32) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

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

	switch op := frame.Payload; op.(type) {
	case *grpcws.GrpcFrame_Header:
		// grpclog.Infof("received Header for stream %v %v", frame.StreamId, frame.GetHeader().Operation)
		if stream != nil {
			ws.writeError(frame.StreamId, "stream already exists")
		}

		if ws.maxStreamCount > 0 && len(ws.activeStreams) > ws.maxStreamCount {
			return ws.writeError(frame.StreamId, "rejecting max number of streams reached for this channel")
		}

		stream := newGrpcStream(frame.StreamId, ws, ws.maxStreamCount)
		func() {
			ws.mutex.Lock()
			defer ws.mutex.Unlock()

			ws.activeStreams[frame.StreamId] = stream
		}()

		url, err := url.Parse("http://localhost/")
		url.Scheme = ws.req.URL.Scheme
		url.Host = ws.req.URL.Host
		url.Path = "/" + frame.GetHeader().Operation
		if err != nil {
			return ws.writeError(frame.StreamId, err.Error())
		}

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

		body := frame.Payload.(*grpcws.GrpcFrame_Body)
		if body.Body.Complete {
			stream.close()
		}

	case *grpcws.GrpcFrame_Cancel:
		// grpclog.Infof("received cancel for stream %v", frame.StreamId)
		if stream == nil {
			// If a stream is being cancelled and it's not there anymore, no error here
			return nil
		}

		// grpclog.Infof("stream %v is canceled", frame.StreamId)
		stream.cancel()
		stream.close()
		ws.deleteStream(frame.StreamId)

	case *grpcws.GrpcFrame_Complete:
		// grpclog.Infof("received Complete for stream %v", frame.StreamId)
		if stream == nil {
			return ws.writeError(frame.StreamId, "stream does not exist")
		}

		// grpclog.Infof("completing stream %v", frame.StreamId)
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

func makeGrpcRequest(req *http.Request) *http.Request {
	// Hack, this should be a shallow copy, but let's see if this works
	req.ProtoMajor = 2
	req.ProtoMinor = 0

	req.Header.Set("content-type", "application/grpc+proto")

	// Remove content-length header since it represents http1.1 payload size, not the sum of the h2
	// DATA frame payload lengths. https://http2.github.io/http2-spec/#malformed This effectively
	// switches to chunked encoding which is the default for h2
	req.Header.Del("content-length")
	return req
}

func (ws *WebsocketChannel) readFrame() (*grpcws.GrpcFrame, error) {
	// we assume a large limit is set for the websocket to avoid handling multiple frames.
	typ, bytesValue, err := ws.wsConn.Read(ws.context)
	if err != nil {
		return nil, err
	}

	if typ != websocket.MessageBinary {
		return nil, errors.New("websocket channel only supports binary messages")
	}

	request := &grpcws.GrpcFrame{}
	if err := proto.Unmarshal(bytesValue, request); err != nil {
		return nil, fmt.Errorf("fram unmarshal error. %w", err)
	}
	return request, nil
}

func (ws *WebsocketChannel) write(frame *grpcws.GrpcFrame) error {
	binaryFrame, err := proto.Marshal(frame)
	if err != nil {
		return err
	}
	return ws.wsConn.Write(ws.context, websocket.MessageBinary, binaryFrame)
}

func (w *WebsocketChannel) enablePing(timeOutInterval time.Duration) {
	w.timeOutInterval = timeOutInterval
	w.timer = timer.NewTimer(w.timeOutInterval)
	go w.ping()
}

func (w *WebsocketChannel) ping() {
	defer w.timer.Stop()

	for {
		select {
		case <-w.context.Done():
			return
		case <-w.timer.C:
			w.timer.Reset(w.timeOutInterval)

			stream := func() *GrpcStream {
				w.mutex.Lock()
				defer w.mutex.Unlock()

				stream, ok := w.activeStreams[0]
				if !ok {
					return nil
				}
				return stream
			}()

			if stream == nil {
				return
			}

			stream.channel.write(framePingResponse)
		}
	}
}

func (ws *WebsocketChannel) Close() {}

func (stream *GrpcStream) Close() error {
	stream.mutex.Lock()
	defer stream.mutex.Unlock()

	if st, ok := stream.responseHeaders["Grpc-Status"]; ok && (len(st) == 0 || st[0] != "0") {
		statusCode := st[0]
		statusMessage, ok := stream.responseHeaders["Grpc-Message"]
		if !ok {
			statusMessage = []string{"Unknown"}
		}

		headers := map[string]*grpcws.HeaderValue{}
		for k, v := range stream.responseHeaders {
			headers[k] = &grpcws.HeaderValue{
				Value: v,
			}
		}

		stream.channel.write(
			&grpcws.GrpcFrame{
				StreamId: stream.id,
				Payload: &grpcws.GrpcFrame_Failure{
					Failure: &grpcws.Failure{
						ErrorMessage: strings.Join(statusMessage, ";"),
						ErrorStatus:  statusCode,
						Headers:      headers,
					},
				},
			},
		)
	} else {
		stream.WriteHeader(http.StatusOK)
	}

	defer stream.cancel()
	stream.channel.write(&grpcws.GrpcFrame{StreamId: stream.id, Payload: &grpcws.GrpcFrame_Complete{Complete: &grpcws.Complete{}}})
	stream.channel.deleteStream(stream.id)
	return nil
}

func (stream *GrpcStream) Write(data []byte) (int, error) {
	stream.WriteHeader(http.StatusOK)
	// grpclog.Infof("write body %v", len(data))

	// Not sure if it is enough to check the writeBuffer length
	if stream.bytesToWrite == 0 && !stream.readHeader {
		stream.bytesToWrite += binary.BigEndian.Uint32(data[1:])
		stream.writeBuffer = data[5:]
		stream.readHeader = true
		return len(data), nil
	} else {
		stream.bytesToWrite -= uint32(len(data))
		stream.writeBuffer = append(stream.writeBuffer, data...)

		if stream.bytesToWrite != 0 {
			return len(data), nil
		}

		err := stream.channel.write(&grpcws.GrpcFrame{
			StreamId: stream.id,
			Payload: &grpcws.GrpcFrame_Body{
				Body: &grpcws.Body{
					Data: data,
				},
			},
		})
		stream.readHeader = false
		return len(data), err
	}
}

func (stream *GrpcStream) WriteHeader(statusCode int) {
	if !stream.hasWrittenHeaders {
		headerResponse := make(map[string]*grpcws.HeaderValue)
		for key, element := range stream.responseHeaders {
			headerResponse[key] = &grpcws.HeaderValue{
				Value: element,
			}
		}
		stream.hasWrittenHeaders = true
		stream.channel.write(
			&grpcws.GrpcFrame{
				StreamId: stream.id,
				Payload: &grpcws.GrpcFrame_Header{
					Header: &grpcws.Header{
						Operation: "",
						Headers:   headerResponse,
						Status:    int32(statusCode),
					},
				},
			})
	}
}
