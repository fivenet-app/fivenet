// Copyright 2017 Improbable. All Rights Reserved.
// See LICENSE for licensing terms.
//
//nolint:canonicalheader // The headers are not canonicalized in the trailers, etc.
package grpcws_test

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	grpcws "github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
	testproto "github.com/improbable-eng/grpc-web/integration_test/go/_proto/improbable/grpcweb/test"
	"github.com/mwitkow/go-conntrack/connhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

var (
	expectedListResponses = 3000
	expectedHeaders       = metadata.Pairs("HeaderTestKey1", "Value1", "HeaderTestKey2", "Value2")
	expectedTrailers      = metadata.Pairs("TrailerTestKey1", "Value1", "TrailerTestKey2", "Value2")
	useFlushForHeaders    = "test-internal-use-flush-for-headers"
)

const (
	grpcWebContentType     = "application/grpc-web"
	grpcWebTextContentType = "application/grpc-web-text"
)

type GrpcWebWrapperTestSuite struct {
	suite.Suite

	httpMajorVersion int
	listener         net.Listener
	grpcServer       *grpc.Server
	wrappedServer    *grpcws.WrappedGrpcServer
}

func TestHttp2GrpcWebWrapperTestSuite(t *testing.T) {
	suite.Run(t, &GrpcWebWrapperTestSuite{httpMajorVersion: 2})
}

func TestHttp1GrpcWebWrapperTestSuite(t *testing.T) {
	suite.Run(t, &GrpcWebWrapperTestSuite{httpMajorVersion: 1})
}

func TestNonRootResource(t *testing.T) {
	grpcServer := grpc.NewServer()
	testproto.RegisterTestServiceServer(grpcServer, &testServiceImpl{})
	wrappedServer := grpcws.WrapServer(grpcServer,
		grpcws.WithAllowNonRootResource(true),
		grpcws.WithOriginFunc(func(_ string) bool {
			return true
		}))

	headers := http.Header{}
	headers.Add("Access-Control-Request-Method", "POST")
	headers.Add("Access-Control-Request-Headers", "origin, x-something-custom, x-grpc-web, accept")
	req := httptest.NewRequest(
		http.MethodOptions,
		"http://host/grpc/improbable.grpcweb.test.TestService/Echo",
		nil,
	)
	req.Header = headers
	resp := httptest.NewRecorder()
	wrappedServer.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func (s *GrpcWebWrapperTestSuite) SetupTest() {
	var err error
	s.grpcServer = grpc.NewServer()
	testproto.RegisterTestServiceServer(s.grpcServer, &testServiceImpl{})
	logger, err := zap.NewDevelopment()
	s.Require().NoError(err)
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.Named("grpc")))
	s.wrappedServer = grpcws.WrapServer(s.grpcServer)

	httpServer := http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			s.Require().Equal(
				s.httpMajorVersion,
				req.ProtoMajor,
				"Requests in this test are served over the wrong protocol",
			)
			s.T().Logf("Serving over: %d", req.ProtoMajor)
			s.wrappedServer.ServeHTTP(resp, req)
		}),
	}

	s.listener, err = net.Listen("tcp", "127.0.0.1:0")
	s.Require().NoError(err, "failed to set up server socket for test")
	tlsConfig, err := connhelpers.TlsConfigForServerCerts(
		filepath.Join(basepath, "../../../internal/tests/certs/localhost.crt"),
		filepath.Join(basepath, "../../../internal/tests/certs/localhost.key"),
	)
	s.Require().NoError(err, "failed loading keys")
	if s.httpMajorVersion == 2 {
		tlsConfig, err = connhelpers.TlsConfigWithHttp2Enabled(tlsConfig)
		s.Require().NoError(err, "failed setting http2")
	}
	s.listener = tls.NewListener(s.listener, tlsConfig)
	go func() {
		httpServer.Serve(s.listener)
	}()

	// Wait for the grpcServer to start serving requests.
	time.Sleep(10 * time.Millisecond)
}

func (s *GrpcWebWrapperTestSuite) timeoutCtxForTest(t *testing.T) context.Context {
	s.T().Helper()
	ctx, _ := context.WithTimeout(t.Context(), 1*time.Second)
	return ctx
}

func (s *GrpcWebWrapperTestSuite) makeRequest(
	verb string, method string, headers http.Header, body io.Reader, isText bool,
) (*http.Response, error) {
	contentType := "application/grpc-web"
	if isText {
		// base64 encode the body
		encodedBody := &bytes.Buffer{}
		encoder := base64.NewEncoder(base64.StdEncoding, encodedBody)
		_, err := io.Copy(encoder, body)
		if err != nil {
			return nil, err
		}
		err = encoder.Close()
		if err != nil {
			return nil, err
		}
		body = encodedBody
		contentType = "application/grpc-web-text"
	}

	url := fmt.Sprintf("https://%s%s", s.listener.Addr().String(), method)
	req, err := http.NewRequest(verb, url, body)
	req = req.WithContext(s.timeoutCtxForTest(s.T()))
	s.Require().NoError(err, "failed creating a request")
	req.Header = headers

	req.Header.Set("content-type", contentType)
	client := &http.Client{
		Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	if s.httpMajorVersion < 2 {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	resp, err := client.Do(req)
	return resp, err
}

func decodeMultipleBase64Chunks(b []byte) ([]byte, error) {
	// grpc-web allows multiple base64 chunks: the implementation may send base64-encoded
	// "chunks" with potential padding whenever the runtime needs to flush a byte buffer.
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md
	output := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	outputEnd := 0

	for inputEnd := 0; inputEnd < len(b); {
		chunk := b[inputEnd:]
		paddingIndex := bytes.IndexByte(chunk, '=')
		if paddingIndex != -1 {
			// find the consecutive =
			for {
				paddingIndex++
				if paddingIndex >= len(chunk) || chunk[paddingIndex] != '=' {
					break
				}
			}
			chunk = chunk[:paddingIndex]
		}
		inputEnd += len(chunk)

		n, err := base64.StdEncoding.Decode(output[outputEnd:], chunk)
		if err != nil {
			return nil, err
		}
		outputEnd += n
	}
	return output[:outputEnd], nil
}

func (s *GrpcWebWrapperTestSuite) makeGrpcRequest(
	method string, reqHeaders http.Header, requestMessages [][]byte, isText bool,
) (http.Header, grpcws.Trailer, [][]byte, error) {
	writer := new(bytes.Buffer)
	for _, msgBytes := range requestMessages {
		grpcPreamble := []byte{0, 0, 0, 0, 0}
		binary.BigEndian.PutUint32(grpcPreamble[1:], uint32(len(msgBytes)))
		writer.Write(grpcPreamble)
		writer.Write(msgBytes)
	}
	resp, err := s.makeRequest("POST", method, reqHeaders, writer, isText)
	if err != nil {
		return nil, grpcws.Trailer{}, nil, err
	}
	defer func() {
		if err == nil {
			return
		}
		resp.Body.Close()
	}()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, grpcws.Trailer{}, nil, err
	}

	if isText {
		contents, err = decodeMultipleBase64Chunks(contents)
		if err != nil {
			return nil, grpcws.Trailer{}, nil, err
		}
	}

	var responseMessages [][]byte
	var trailers grpcws.Trailer
	reader := bytes.NewReader(contents)
	for {
		grpcPreamble := []byte{0, 0, 0, 0, 0}
		readCount, err := reader.Read(grpcPreamble)
		if errors.Is(err, io.EOF) {
			break
		}
		if readCount != 5 || err != nil {
			return nil, grpcws.Trailer{}, nil, fmt.Errorf(
				"Unexpected end of body in preamble. %w",
				err,
			)
		}
		payloadLength := binary.BigEndian.Uint32(grpcPreamble[1:])
		payloadBytes := make([]byte, payloadLength)

		readCount, err = reader.Read(payloadBytes)
		if uint32(readCount) != payloadLength || err != nil {
			return nil, grpcws.Trailer{}, nil, fmt.Errorf("Unexpected end of msg: %w", err)
		}
		if grpcPreamble[0]&(1<<7) == (1 << 7) { // MSB signifies the trailer parser
			trailers = readTrailersFromBytes(s.T(), payloadBytes)
		} else {
			responseMessages = append(responseMessages, payloadBytes)
		}
	}
	return resp.Header, trailers, responseMessages, nil
}

func (s *GrpcWebWrapperTestSuite) TestPingEmpty() {
	headers, trailers, responses, err := s.makeGrpcRequest(
		"/improbable.grpcweb.test.TestService/PingEmpty",
		headerWithFlag(),
		serializeProtoMessages([]proto.Message{&emptypb.Empty{}}),
		false)
	s.Require().NoError(err, "No error on making request")

	s.Len(responses, 1, "PingEmpty is an unary response")
	s.assertTrailerGrpcCode(trailers, codes.OK, "")
	s.assertHeadersContainMetadata(headers, expectedHeaders)
	s.assertTrailersContainMetadata(trailers, expectedTrailers)
	s.assertContentTypeSet(headers, grpcWebContentType)
}

func (s *GrpcWebWrapperTestSuite) TestPing() {
	// test both the text and binary formats
	for _, contentType := range []string{grpcWebContentType, grpcWebTextContentType} {
		headers, trailers, responses, err := s.makeGrpcRequest(
			"/improbable.grpcweb.test.TestService/Ping",
			headerWithFlag(),
			serializeProtoMessages([]proto.Message{&testproto.PingRequest{Value: "foo"}}),
			contentType == grpcWebTextContentType)
		s.Require().NoError(err, "No error on making request")

		s.Len(responses, 1, "PingEmpty is an unary response")
		s.assertTrailerGrpcCode(trailers, codes.OK, "")
		s.assertHeadersContainMetadata(headers, expectedHeaders)
		s.assertTrailersContainMetadata(trailers, expectedTrailers)
		s.assertHeadersContainCorsExpectedHeaders(headers, expectedHeaders)
		s.assertContentTypeSet(headers, contentType)
	}
}

func (s *GrpcWebWrapperTestSuite) TestPingError_WithTrailersInData() {
	// gRPC-Web spec says that if there is no payload to an answer, the trailers (including grpc-status) must be in the
	// headers and not in trailers. However, that's not true if SendHeaders are pushed before. This tests this.
	headers, trailers, responses, err := s.makeGrpcRequest(
		"/improbable.grpcweb.test.TestService/PingError",
		headerWithFlag(useFlushForHeaders),
		serializeProtoMessages([]proto.Message{&emptypb.Empty{}}),
		false)
	s.Require().NoError(err, "No error on making request")

	s.Empty(responses, "PingError is an unary response that has no payload")
	s.assertTrailerGrpcCode(trailers, codes.Unimplemented, "Not implemented PingError")
	s.assertHeadersContainMetadata(headers, expectedHeaders)
	s.assertTrailersContainMetadata(trailers, expectedTrailers)
	s.assertHeadersContainCorsExpectedHeaders(headers, expectedHeaders)
	s.assertContentTypeSet(headers, grpcWebContentType)
}

func (s *GrpcWebWrapperTestSuite) TestPingError_WithTrailersInHeaders() {
	// gRPC-Web spec says that if there is no payload to an answer, the trailers (including grpc-status) must be in the
	// headers and not in trailers.
	headers, _, responses, err := s.makeGrpcRequest(
		"/improbable.grpcweb.test.TestService/PingError",
		http.Header{},
		serializeProtoMessages([]proto.Message{&emptypb.Empty{}}),
		false)
	s.Require().NoError(err, "No error on making request")

	s.Empty(responses, "PingError is an unary response that has no payload")
	s.assertHeadersGrpcCode(headers, codes.Unimplemented, "Not implemented PingError")
	// s.assertHeadersContainMetadata(headers, expectedHeaders) // (mwitkow): There is a bug in gRPC where headers don't get added if no payload exists.
	s.assertHeadersContainMetadata(headers, expectedTrailers)
	s.assertHeadersContainCorsExpectedHeaders(headers, expectedTrailers)
	s.assertContentTypeSet(headers, grpcWebContentType)
}

func (s *GrpcWebWrapperTestSuite) TestPingList() {
	headers, trailers, responses, err := s.makeGrpcRequest(
		"/improbable.grpcweb.test.TestService/PingList",
		headerWithFlag(),
		serializeProtoMessages([]proto.Message{&testproto.PingRequest{Value: "something"}}),
		false)
	s.Require().NoError(err, "No error on making request")
	s.Len(
		responses, expectedListResponses,
		"the number of expected proto fields shouold match",
	)
	s.assertTrailerGrpcCode(trailers, codes.OK, "")
	s.assertHeadersContainMetadata(headers, expectedHeaders)
	s.assertTrailersContainMetadata(trailers, expectedTrailers)
	s.assertHeadersContainCorsExpectedHeaders(headers, expectedHeaders)
	s.assertContentTypeSet(headers, grpcWebContentType)
}

func (s *GrpcWebWrapperTestSuite) getStandardGrpcClient() *grpc.ClientConn {
	conn, err := grpc.NewClient(s.listener.Addr().String(),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})),
	)
	s.Require().NoError(err, "grpc dial must succeed")
	return conn
}

func (s *GrpcWebWrapperTestSuite) TestPingList_NormalGrpcWorks() {
	if s.httpMajorVersion < 2 {
		s.T().Skipf("Standard gRPC interop only works over HTTP2")
		return
	}
	conn := s.getStandardGrpcClient()
	client := testproto.NewTestServiceClient(conn)
	pingListClient, err := client.PingList(
		s.timeoutCtxForTest(s.T()),
		&testproto.PingRequest{Value: "foo", ResponseCount: 10},
	)
	s.Require().NoError(err, "no error during execution")
	for {
		_, err := pingListClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		s.Require().NoError(err, "no error during execution")
	}
	recvHeaders, err := pingListClient.Header()
	s.Require().NoError(err, "no error during execution")
	recvTrailers := pingListClient.Trailer()
	allExpectedHeaders := metadata.Join(
		metadata.MD{
			"content-type": []string{"application/grpc"},
			"trailer":      []string{"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin"},
		}, expectedHeaders)
	s.Equal(allExpectedHeaders, recvHeaders, "expected headers must be received")
	s.Equal(expectedTrailers, recvTrailers, "expected trailers must be received")
}

func (s *GrpcWebWrapperTestSuite) TestPingStream_NormalGrpcWorks() {
	if s.httpMajorVersion < 2 {
		s.T().Skipf("Standard gRPC interop only works over HTTP2")
		return
	}
	conn := s.getStandardGrpcClient()
	client := testproto.NewTestServiceClient(conn)
	bidiClient, err := client.PingStream(s.timeoutCtxForTest(s.T()))
	s.Require().NoError(err, "no error during execution")
	bidiClient.Send(&testproto.PingRequest{Value: "one"})
	bidiClient.Send(&testproto.PingRequest{Value: "two"})
	resp, err := bidiClient.CloseAndRecv()
	s.Require().NoError(err, "no error during execution")
	s.Equal("one,two", resp.GetValue(), "expected concatenated value must be received")
	recvHeaders, err := bidiClient.Header()
	s.Require().NoError(err, "no error during execution")
	recvTrailers := bidiClient.Trailer()
	allExpectedHeaders := metadata.Join(
		metadata.MD{
			"content-type": []string{"application/grpc"},
			"trailer":      []string{"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin"},
		}, expectedHeaders)
	s.Equal(allExpectedHeaders, recvHeaders, "expected headers must be received")
	s.Equal(expectedTrailers, recvTrailers, "expected trailers must be received")
}

func (s *GrpcWebWrapperTestSuite) TestCORSPreflight_DeniedByDefault() {
	/**
	OPTIONS /improbable.grpcweb.test.TestService/Ping
	Access-Control-Request-Method: POST
	Access-Control-Request-Headers: origin, x-requested-with, accept
	Origin: http://foo.client.com
	*/
	headers := http.Header{}
	headers.Add("Access-Control-Request-Method", "POST")
	headers.Add("Access-Control-Request-Headers", "origin, x-something-custom, x-grpc-web, accept")
	headers.Add("Origin", "https://foo.client.com")

	corsResp, err := s.makeRequest(
		"OPTIONS",
		"/improbable.grpcweb.test.TestService/PingList",
		headers,
		nil,
		false,
	)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")

	preflight := corsResp.Header
	s.Empty(
		preflight.Get("Access-Control-Allow-Origin"),
		"origin must not be in the response headers",
	)
	s.Empty(
		preflight.Get("Access-Control-Allow-Methods"),
		"allowed methods must not be in the response headers",
	)
	s.Empty(
		preflight.Get("Access-Control-Max-Age"),
		"allowed max age must not be in the response headers",
	)
	s.Empty(
		preflight.Get("Access-Control-Allow-Headers"),
		"allowed headers must not be in the response headers",
	)
}

func (s *GrpcWebWrapperTestSuite) TestCORSPreflight_AllowedByOriginFunc() {
	/**
	OPTIONS /improbable.grpcweb.test.TestService/Ping
	Access-Control-Request-Method: POST
	Access-Control-Request-Headers: origin, x-requested-with, accept
	Origin: http://foo.client.com
	*/
	headers := http.Header{}
	headers.Add("Access-Control-Request-Method", "POST")
	headers.Add("Access-Control-Request-Headers", "origin, x-something-custom, x-grpc-web, accept")
	headers.Add("Origin", "https://foo.client.com")

	// Create a new server which permits Cross-Origin Resource requests from `foo.client.com`.
	s.wrappedServer = grpcws.WrapServer(s.grpcServer,
		grpcws.WithOriginFunc(func(origin string) bool {
			return origin == "https://foo.client.com"
		}),
	)

	corsResp, err := s.makeRequest(
		"OPTIONS",
		"/improbable.grpcweb.test.TestService/PingList",
		headers,
		nil,
		false,
	)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")

	preflight := corsResp.Header
	s.Equal(
		"https://foo.client.com",
		preflight.Get("Access-Control-Allow-Origin"),
		"origin must be in the response headers",
	)
	s.Equal(
		"POST",
		preflight.Get("Access-Control-Allow-Methods"),
		"allowed methods must be in the response headers",
	)
	s.Equal(
		"600",
		preflight.Get("Access-Control-Max-Age"),
		"allowed max age must be in the response headers",
	)
	s.Equal(
		"origin, x-something-custom, x-grpc-web, accept",
		preflight.Get("Access-Control-Allow-Headers"),
		"allowed headers must be in the response headers",
	)

	corsResp, err = s.makeRequest(
		"OPTIONS",
		"/improbable.grpcweb.test.TestService/Unknown",
		headers,
		nil,
		false,
	)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")
	s.Equal(
		http.StatusMethodNotAllowed,
		corsResp.StatusCode,
		"cors should return 405 (MethodNotAllowed) as grpc server does not handle requests other than POST requests",
	)
}

func (s *GrpcWebWrapperTestSuite) TestCORSPreflight_CorsMaxAge() {
	/**
	OPTIONS /improbable.grpcweb.test.TestService/Ping
	Access-Control-Request-Method: POST
	Access-Control-Request-Headers: origin, x-requested-with, accept
	Origin: http://foo.client.com
	*/
	headers := http.Header{}
	headers.Add("Access-Control-Request-Method", "POST")
	headers.Add("Access-Control-Request-Headers", "origin, x-something-custom, x-grpc-web, accept")
	headers.Add("Origin", "https://foo.client.com")

	// Create a new server which customizes a cache time of the preflight request to a Cross-Origin Resource.
	s.wrappedServer = grpcws.WrapServer(s.grpcServer,
		grpcws.WithOriginFunc(func(string) bool {
			return true
		}),
		grpcws.WithCorsMaxAge(time.Hour),
	)

	corsResp, err := s.makeRequest(
		"OPTIONS",
		"/improbable.grpcweb.test.TestService/PingList",
		headers,
		nil,
		false,
	)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")

	preflight := corsResp.Header
	s.Equal(
		"3600",
		preflight.Get("Access-Control-Max-Age"),
		"allowed max age must be in the response headers",
	)
}

func (s *GrpcWebWrapperTestSuite) TestCORSPreflight_EndpointsOnlyTrueWithHandlerFunc() {
	/**
	OPTIONS /improbable.grpcweb.test.TestService/Ping
	Access-Control-Request-Method: POST
	Access-Control-Request-Headers: origin, x-requested-with, accept
	Origin: http://foo.client.com
	*/
	headers := http.Header{}
	headers.Add("Access-Control-Request-Method", "POST")
	headers.Add("Access-Control-Request-Headers", "origin, x-something-custom, x-grpc-web, accept")
	headers.Add("Origin", "https://foo.client.com")

	// Create a new grpc server that uses WrapHandler and the WithEndpointsFunc option
	const pingMethod = "/improbable.grpcweb.test.TestService/PingList"
	const badMethod = "/improbable.grpcweb.test.TestService/Bad"

	mux := http.NewServeMux()
	mux.Handle(pingMethod, s.grpcServer)
	mux.Handle(badMethod, http.HandlerFunc(func(resp http.ResponseWriter, _ *http.Request) {
		resp.WriteHeader(http.StatusUnauthorized)
	}))

	s.wrappedServer = grpcws.WrapHandler(mux,
		grpcws.WithEndpointsFunc(func() []string {
			return []string{pingMethod}
		}),
		grpcws.WithOriginFunc(func(origin string) bool {
			return origin == "https://foo.client.com"
		}),
	)

	corsResp, err := s.makeRequest("OPTIONS", pingMethod, headers, nil, false)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")

	preflight := corsResp.Header
	s.Equal(
		"https://foo.client.com",
		preflight.Get("Access-Control-Allow-Origin"),
		"origin must be in the response headers",
	)
	s.Equal(
		"POST",
		preflight.Get("Access-Control-Allow-Methods"),
		"allowed methods must be in the response headers",
	)
	s.Equal(
		"600",
		preflight.Get("Access-Control-Max-Age"),
		"allowed max age must be in the response headers",
	)
	s.Equal(
		"origin, x-something-custom, x-grpc-web, accept",
		preflight.Get("Access-Control-Allow-Headers"),
		"allowed headers must be in the response headers",
	)

	corsResp, err = s.makeRequest("OPTIONS", badMethod, headers, nil, false)
	defer func() {
		if err == nil {
			return
		}
		corsResp.Body.Close()
	}()
	s.Require().NoError(err, "cors preflight should not return errors")
	s.Equal(401, corsResp.StatusCode, "cors should return 403 as mocked")
}

func (s *GrpcWebWrapperTestSuite) assertHeadersContainMetadata(
	headers http.Header,
	meta metadata.MD,
) {
	for k, v := range meta {
		lowerKey := strings.ToLower(k)
		for _, vv := range v {
			s.Equal(
				headers.Get(lowerKey),
				vv,
				"Expected there to be %v=%v",
				lowerKey,
				vv,
			)
		}
	}
}

func (s *GrpcWebWrapperTestSuite) assertContentTypeSet(headers http.Header, contentType string) {
	s.Equal(
		contentType,
		headers.Get("content-type"),
		`Expected there to be content-type=%v`,
		contentType,
	)
}

func (s *GrpcWebWrapperTestSuite) assertTrailersContainMetadata(
	trailers grpcws.Trailer,
	meta metadata.MD,
) {
	for k, v := range meta {
		for _, vv := range v {
			s.Equal(trailers.Get(k), vv, "Expected there to be %v=%v", k, vv)
		}
	}
}

func (s *GrpcWebWrapperTestSuite) assertHeadersContainCorsExpectedHeaders(
	headers http.Header,
	meta metadata.MD,
) {
	value := headers.Get("Access-Control-Expose-Headers")
	s.NotEmpty(value, "cors: access control expose headers should not be empty")
	for k := range meta {
		if k == "Access-Control-Expose-Headers" {
			continue
		}
		s.Contains(
			value,
			http.CanonicalHeaderKey(k),
			"cors: exposed headers should contain metadata",
		)
	}
}

func (s *GrpcWebWrapperTestSuite) assertHeadersGrpcCode(
	headers http.Header,
	code codes.Code,
	desc string,
) {
	s.Require().NotEmpty(headers.Get("grpc-status"), "grpc-status must not be empty in trailers")
	statusCode, err := strconv.Atoi(headers.Get("grpc-status"))
	s.Require().NoError(err, "no error parsing grpc-status")
	s.EqualValues(code, statusCode, "grpc-status must match expected code")
	s.Equal(
		desc,
		headers.Get("grpc-message"),
		"grpc-message is expected to match",
	)
}

func (s *GrpcWebWrapperTestSuite) assertTrailerGrpcCode(
	trailers grpcws.Trailer,
	code codes.Code,
	desc string,
) {
	s.Require().NotEmpty(
		trailers.Get("grpc-status"),
		"grpc-status must not be empty in trailers",
	)
	statusCode, err := strconv.Atoi(trailers.Get("grpc-status"))
	s.Require().NoError(err, "no error parsing grpc-status")
	s.EqualValues(code, statusCode, "grpc-status must match expected code")
	s.Equal(
		desc,
		trailers.Get("grpc-message"),
		"grpc-message is expected to match",
	)
}

func serializeProtoMessages(messages []proto.Message) [][]byte {
	out := [][]byte{}
	for _, m := range messages {
		b, _ := proto.Marshal(m)
		out = append(out, b)
	}
	return out
}

func readTrailersFromBytes(t *testing.T, dataBytes []byte) grpcws.Trailer {
	t.Helper()

	bufferReader := bytes.NewBuffer(dataBytes)
	tp := textproto.NewReader(bufio.NewReader(bufferReader))

	// First, read bytes as MIME headers.
	// However, it normalizes header names by textproto.CanonicalMIMEHeaderKey.
	// In the next step, replace header names by raw one.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err == nil {
		return grpcws.Trailer{}
	}

	trailers := make(http.Header)
	bufferReader = bytes.NewBuffer(dataBytes)
	tp = textproto.NewReader(bufio.NewReader(bufferReader))

	// Second, replace header names because gRPC Web trailer names must be lower-case.
	for {
		line, err := tp.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err, "failed to read header line")

		i := strings.IndexByte(line, ':')
		if i == -1 {
			require.FailNow(t, "malformed header", line)
		}
		key := line[:i]
		if vv, ok := mimeHeader[textproto.CanonicalMIMEHeaderKey(key)]; ok {
			trailers[key] = vv
		}
	}
	return grpcws.HTTPTrailerToGrpcWebTrailer(trailers)
}

func headerWithFlag(flags ...string) http.Header {
	h := http.Header{}
	for _, f := range flags {
		h.Set(f, "true")
	}
	return h
}

type testServiceImpl struct{}

func (s *testServiceImpl) PingEmpty(
	ctx context.Context,
	_ *emptypb.Empty,
) (*testproto.PingResponse, error) {
	grpc.SendHeader(ctx, expectedHeaders)
	grpclog.Info("handling PingEmpty")
	grpc.SetTrailer(ctx, expectedTrailers)
	return &testproto.PingResponse{Value: "foobar"}, nil
}

func (s *testServiceImpl) Ping(
	ctx context.Context,
	ping *testproto.PingRequest,
) (*testproto.PingResponse, error) {
	grpc.SendHeader(ctx, expectedHeaders)
	grpclog.Info("handling Ping")
	grpc.SetTrailer(ctx, expectedTrailers)
	return &testproto.PingResponse{Value: ping.GetValue()}, nil
}

func (s *testServiceImpl) PingError(
	ctx context.Context,
	ping *testproto.PingRequest,
) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	if _, exists := md[useFlushForHeaders]; exists {
		grpc.SendHeader(ctx, expectedHeaders)
		grpclog.Info("handling PingError with flushed headers")
	} else {
		grpc.SetHeader(ctx, expectedHeaders)
		grpclog.Info("handling PingError without flushing")
	}
	grpc.SetTrailer(ctx, expectedTrailers)
	return nil, status.Errorf(codes.Unimplemented, "Not implemented PingError")
}

func (s *testServiceImpl) PingList(
	ping *testproto.PingRequest,
	stream testproto.TestService_PingListServer,
) error {
	stream.SendHeader(expectedHeaders)
	stream.SetTrailer(expectedTrailers)
	grpclog.Info("handling PingList")
	for i := range int32(expectedListResponses) {
		stream.Send(
			&testproto.PingResponse{Value: fmt.Sprintf("%s %d", ping.GetValue(), i), Counter: i},
		)
	}
	return nil
}

func (s *testServiceImpl) PingStream(stream testproto.TestService_PingStreamServer) error {
	stream.SendHeader(expectedHeaders)
	stream.SetTrailer(expectedTrailers)
	grpclog.Info("handling PingStream")
	allValues := ""
	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			stream.SendAndClose(&testproto.PingResponse{
				Value: allValues,
			})
			return nil
		}
		if err != nil {
			return err
		}
		if allValues == "" {
			allValues = in.GetValue()
		} else {
			allValues = allValues + "," + in.GetValue()
		}
		if in.GetFailureType() == testproto.PingRequest_CODE {
			if in.GetErrorCodeReturned() == 0 {
				stream.SendAndClose(&testproto.PingResponse{
					Value: allValues,
				})
				return nil
			}
			return status.Errorf(
				codes.Code(in.GetErrorCodeReturned()),
				"Intentionally returning status code: %d",
				in.GetErrorCodeReturned(),
			)
		}
	}
}

func (s *testServiceImpl) Echo(
	ctx context.Context,
	text *testproto.TextMessage,
) (*testproto.TextMessage, error) {
	grpc.SendHeader(ctx, expectedHeaders)
	grpclog.Info("handling Echo")
	grpc.SetTrailer(ctx, expectedTrailers)
	return text, nil
}

func (s *testServiceImpl) PingPongBidi(stream testproto.TestService_PingPongBidiServer) error {
	stream.SendHeader(expectedHeaders)
	stream.SetTrailer(expectedTrailers)
	grpclog.Info("handling PingPongBidi")
	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if in.GetFailureType() == testproto.PingRequest_CODE {
			if in.GetErrorCodeReturned() == 0 {
				stream.Send(&testproto.PingResponse{
					Value: in.GetValue(),
				})
				return nil
			}
			return status.Errorf(
				codes.Code(in.GetErrorCodeReturned()),
				"Intentionally returning status code: %d",
				in.GetErrorCodeReturned(),
			)
		}
	}
}
