// Copyright 2017 Improbable. All Rights Reserved.
// See LICENSE for licensing terms.

package grpcws

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var internalRequestHeadersWhitelist = []string{
	"U-A", // for gRPC-Web User Agent indicator.
}

// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md#protocol-differences-vs-grpc-over-http2
const (
	grpcContentType        = "application/grpc"
	grpcWebContentType     = "application/grpc-web"
	grpcWebTextContentType = "application/grpc-web-text"
)

type WrappedGrpcServer struct {
	handler             http.Handler
	opts                *options
	corsWrapperHandler  http.Handler
	originFunc          func(origin string) bool
	websocketOriginFunc func(req *http.Request) bool
	websocketReadLimit  int64
	allowedHeaders      []string
	endpointFunc        func(req *http.Request) string
	registeredEndpoints []string
}

// WrapServer takes a gRPC Server in Go and returns a *WrappedGrpcServer that provides gRPC-Web Compatibility.
//
// The internal implementation fakes out a http.Request that carries standard gRPC, and performs the remapping inside
// http.ResponseWriter, i.e. mostly the re-encoding of Trailers (that carry gRPC status).
//
// You can control the behaviour of the wrapper (e.g. modifying CORS behaviour) using `With*` options.
func WrapServer(server *grpc.Server, options ...Option) *WrappedGrpcServer {
	return wrapGrpc(options, server, func() []string {
		return ListGRPCResources(server)
	})
}

// WrapHandler takes a http.Handler (such as a http.Mux) and returns a *WrappedGrpcServer that provides gRPC-Web
// Compatibility.
//
// This behaves nearly identically to WrapServer except when the WithCorsForRegisteredEndpointsOnly setting is true.
// Then a WithEndpointsFunc option must be provided or all CORS requests will NOT be handled.
func WrapHandler(handler http.Handler, options ...Option) *WrappedGrpcServer {
	return wrapGrpc(options, handler, func() []string {
		return []string{}
	})
}

func wrapGrpc(options []Option, handler http.Handler, endpointsFunc func() []string) *WrappedGrpcServer {
	opts := evaluateOptions(options)
	allowedHeaders := append(opts.allowedRequestHeaders, internalRequestHeadersWhitelist...)
	corsWrapper := cors.New(cors.Options{
		AllowOriginFunc:  opts.originFunc,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   nil,  // make sure that this is *nil*, otherwise the WebResponse overwrite will not work.
		AllowCredentials: true, // always allow credentials, otherwise :authorization headers won't work
		MaxAge:           int(opts.corsMaxAge.Seconds()),
	})
	websocketOriginFunc := opts.websocketOriginFunc
	if websocketOriginFunc == nil {
		websocketOriginFunc = defaultWebsocketOriginFunc
	}

	endpointFunc := func(req *http.Request) string {
		return req.URL.Path
	}

	if opts.allowNonRootResources {
		endpointFunc = getGRPCEndpoint
	}

	if opts.endpointsFunc != nil {
		endpointsFunc = *opts.endpointsFunc
	}

	w := &WrappedGrpcServer{
		handler:             handler,
		opts:                opts,
		originFunc:          opts.originFunc,
		websocketOriginFunc: websocketOriginFunc,
		websocketReadLimit:  opts.websocketReadLimit,
		allowedHeaders:      allowedHeaders,
		endpointFunc:        endpointFunc,
		registeredEndpoints: endpointsFunc(),
	}
	w.corsWrapperHandler = corsWrapper.Handler(http.HandlerFunc(w.HandleGrpcWebRequest))

	return w
}

// ServeHTTP takes a HTTP request and if it is a gRPC-Web request wraps it with a compatibility layer to transform it to
// a standard gRPC request for the wrapped gRPC server and transforms the response to comply with the gRPC-Web protocol.
//
// The gRPC-Web compatibility is only invoked if the request is a gRPC-Web request as determined by IsGrpcWebRequest or
// the request is a pre-flight (CORS) request as determined by IsAcceptableGrpcCorsRequest.
//
// You can control the CORS behaviour using `With*` options in the WrapServer function.
func (w *WrappedGrpcServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if w.IsGrpcWebSocketChannelRequest(req) {
		if w.websocketOriginFunc(req) {
			w.HandleGrpcWebsocketChannelRequest(resp, req)
			return
		}

		resp.WriteHeader(http.StatusForbidden)
		_, _ = resp.Write(make([]byte, 0))
		return
	}

	if w.IsAcceptableGrpcCorsRequest(req) || w.IsGrpcWebRequest(req) {
		w.corsWrapperHandler.ServeHTTP(resp, req)
		return
	}

	w.handler.ServeHTTP(resp, req)
}

// HandleGrpcWebRequest takes a HTTP request that is assumed to be a gRPC-Web request and wraps it with a compatibility
// layer to transform it to a standard gRPC request for the wrapped gRPC server and transforms the response to comply
// with the gRPC-Web protocol.
func (w *WrappedGrpcServer) HandleGrpcWebRequest(resp http.ResponseWriter, req *http.Request) {
	intReq, isTextFormat := hackIntoNormalGrpcRequest(req)
	intResp := newGrpcWebResponse(resp, isTextFormat)
	intReq.URL.Path = w.endpointFunc(intReq)
	w.handler.ServeHTTP(intResp, intReq)
	intResp.finishRequest(req)
}

// IsGrpcWebSocketRequest determines if a request is a gRPC-Web request by checking that the "Upgrade" header is set and
// "Sec-Websocket-Protocol" header value is "grpc-websocket-channel" and that the "root" path is requested
func (w *WrappedGrpcServer) IsGrpcWebSocketChannelRequest(req *http.Request) bool {
	if strings.ToLower(req.Header.Get("Upgrade")) != "websocket" {
		return false
	}

	for _, subproto := range req.Header.Values("Sec-Websocket-Protocol") {
		for token := range strings.SplitSeq(subproto, ",") {
			token = strings.TrimSpace(token)
			if strings.EqualFold(token, "grpc-websocket-channel") {
				return true
			}
		}
	}

	return false
}

// IsGrpcWebRequest determines if a request is a gRPC-Web request by checking that the "content-type" is
// "application/grpc-web" and that the method is POST.
func (w *WrappedGrpcServer) IsGrpcWebRequest(req *http.Request) bool {
	return req.Method == http.MethodPost && strings.HasPrefix(req.Header.Get("content-type"), grpcWebContentType)
}

// HandleGrpcWebsocketChannelRequest takes a HTTP request that is assumed to be a gRPC-Websocket-channel request and starts a
// duplexed grpc-websocket-channel which will create multiple virtual streams over a single websocket.
func (w *WrappedGrpcServer) HandleGrpcWebsocketChannelRequest(resp http.ResponseWriter, req *http.Request) {
	grpclog.Infof("handle grpc channel request %s", req.Host)

	wsConn, err := websocket.Accept(resp, req, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // managed by ServeHTTP
		Subprotocols:       []string{"grpc-websocket-channel"},
	})
	if err != nil {
		grpclog.Errorf("unable to upgrade websocket request: %v", err)
		return
	}
	defer wsConn.CloseNow()

	wsConn.SetReadLimit(w.websocketReadLimit)

	headers := make(http.Header)
	for _, name := range w.allowedHeaders {
		if values, exist := req.Header[name]; exist {
			headers[name] = values
		}
	}
	req.Header = headers

	ctx, cancelFunc := context.WithCancel(req.Context())
	defer cancelFunc()

	websocketChannel := NewWebsocketChannel(wsConn, w.handler, ctx, w.opts.websocketChannelMaxStreamCount, req)
	if w.opts.websocketPingInterval >= time.Second {
		websocketChannel.enablePing(w.opts.websocketPingInterval)
	}

	websocketChannel.Start()
}

// IsAcceptableGrpcCorsRequest determines if a request is a CORS pre-flight request for a gRPC-Web request and that this
// request is acceptable for CORS.
//
// You can control the CORS behaviour using `With*` options in the WrapServer function.
func (w *WrappedGrpcServer) IsAcceptableGrpcCorsRequest(req *http.Request) bool {
	accessControlHeaders := strings.ToLower(req.Header.Get("Access-Control-Request-Headers"))
	if req.Method == http.MethodOptions && strings.Contains(accessControlHeaders, "x-grpc-web") {
		if w.opts.corsForRegisteredEndpointsOnly {
			return w.isRequestForRegisteredEndpoint(req)
		}
		return true
	}
	return false
}

func (w *WrappedGrpcServer) isRequestForRegisteredEndpoint(req *http.Request) bool {
	requestedEndpoint := w.endpointFunc(req)
	return slices.Contains(w.registeredEndpoints, requestedEndpoint)
}

// readerCloser combines an io.Reader and an io.Closer into an io.ReadCloser.
type readerCloser struct {
	reader io.Reader
	closer io.Closer
}

func (r *readerCloser) Read(dest []byte) (int, error) {
	return r.reader.Read(dest)
}

func (r *readerCloser) Close() error {
	return r.closer.Close()
}

func hackIntoNormalGrpcRequest(req *http.Request) (*http.Request, bool) {
	// Hack, this should be a shallow copy, but let's see if this works
	req.ProtoMajor = 2
	req.ProtoMinor = 0

	contentType := req.Header.Get("content-type")
	incomingContentType := grpcWebContentType
	isTextFormat := strings.HasPrefix(contentType, grpcWebTextContentType)
	if isTextFormat {
		// body is base64-encoded: decode it; Wrap it in readerCloser so Body is still closed
		decoder := base64.NewDecoder(base64.StdEncoding, req.Body)
		req.Body = &readerCloser{reader: decoder, closer: req.Body}
		incomingContentType = grpcWebTextContentType
	}
	req.Header.Set("content-type", strings.Replace(contentType, incomingContentType, grpcContentType, 1))

	// Remove content-length header since it represents http1.1 payload size, not the sum of the h2
	// DATA frame payload lengths. https://http2.github.io/http2-spec/#malformed This effectively
	// switches to chunked encoding which is the default for h2
	req.Header.Del("content-length")

	return req, isTextFormat
}

func defaultWebsocketOriginFunc(req *http.Request) bool {
	origin, err := WebsocketRequestOrigin(req)
	if err != nil {
		grpclog.Warning(err)
		return false
	}
	return origin == req.Host
}
