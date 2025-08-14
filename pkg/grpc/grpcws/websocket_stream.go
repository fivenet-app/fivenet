package grpcws

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/grpcws"
)

var ErrClient = errors.New("grpc client error")

type GrpcStream struct {
	id                uint32
	hasWrittenHeaders bool
	responseHeaders   http.Header
	inputFrames       chan *grpcws.GrpcFrame
	channel           *WebsocketChannel
	ctx               context.Context
	cancel            context.CancelFunc
	remainingBuffer   []byte
	remainingError    error
	bytesToWrite      uint32
	writeBuffer       []byte
	inputClosed       bool
}

func newGrpcStream(streamId uint32, channel *WebsocketChannel, streamBufferSize int) *GrpcStream {
	ctx, cancel := context.WithCancel(channel.ctx)
	return &GrpcStream{
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

func (stream *GrpcStream) Read(p []byte) (int, error) {
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
	if more {
		switch op := frame.GetPayload().(type) {
		case *grpcws.GrpcFrame_Body:
			stream.remainingBuffer = op.Body.GetData()
			return stream.Read(p)

		case *grpcws.GrpcFrame_Failure:
			return 0, ErrClient
		}
	}
	return 0, io.EOF
}

func (stream *GrpcStream) Close() error {
	if st, ok := stream.responseHeaders["grpc-status"]; ok && (len(st) == 0 || st[0] != "0") {
		statusCode := st[0]
		statusMessage, ok := stream.responseHeaders["grpc-message"]
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

	stream.cancel()
	stream.channel.write(
		&grpcws.GrpcFrame{
			StreamId: stream.id,
			Payload:  &grpcws.GrpcFrame_Complete{Complete: &grpcws.Complete{}},
		},
	)
	return nil
}

func (stream *GrpcStream) Write(data []byte) (int, error) {
	stream.WriteHeader(http.StatusOK)

	// Not sure if it is enough to check the writeBuffer length
	if stream.bytesToWrite == 0 && len(data) != 0 {
		stream.bytesToWrite += binary.BigEndian.Uint32(data[1:])
		stream.writeBuffer = data[5:]
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
					Data: stream.writeBuffer,
				},
			},
		})
		stream.writeBuffer = nil
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
