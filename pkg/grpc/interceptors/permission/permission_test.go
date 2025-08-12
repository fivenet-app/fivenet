package grpc_permission

import (
	"context"
	"strconv"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func failFromMD(ctx context.Context) (bool, error) {
	val := metadata.ExtractIncoming(ctx).Get("fail")
	if val == "" {
		return false, status.Error(codes.Internal, "Failed to get fail info from metadata")
	}

	return strconv.ParseBool(val)
}

func buildDummyUnaryPermsFunction(
	t *testing.T,
	hasRemap bool,
) func(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
	t.Helper()
	return func(ctx context.Context, info *grpc.UnaryServerInfo) (context.Context, error) {
		fail, err := failFromMD(ctx)
		if err != nil {
			return nil, err
		}

		_, ok := info.Server.(GetPermsRemapFunc)
		assert.Equal(t, hasRemap, ok, "expected grpc server have or have not a perms remap func")

		if fail {
			return nil, status.Error(
				codes.PermissionDenied,
				"buildDummyUnaryPermsFunction fail set in context",
			)
		}
		return ctx, nil
	}
}

func buildDummyStreamPermsFunction(
	t *testing.T,
	hasRemap bool,
) func(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error) {
	t.Helper()
	return func(ctx context.Context, srv any, info *grpc.StreamServerInfo) (context.Context, error) {
		fail, err := failFromMD(ctx)
		if err != nil {
			return nil, err
		}

		_, ok := srv.(GetPermsRemapFunc)
		assert.Equal(t, hasRemap, ok)

		if fail {
			return nil, status.Error(
				codes.PermissionDenied,
				"buildDummyStreamPermsFunction fail set in context",
			)
		}
		return ctx, nil
	}
}

func ctxWithFail(ctx context.Context, fail bool) context.Context {
	md := grpcMetadata.Pairs("fail", strconv.FormatBool(fail))
	return metadata.MD(md).ToOutgoing(ctx)
}

type assertingPingService struct {
	testpb.TestServiceServer

	T *testing.T
}

func (s *assertingPingService) PingError(
	ctx context.Context,
	ping *testpb.PingErrorRequest,
) (*testpb.PingErrorResponse, error) {
	return s.TestServiceServer.PingError(ctx, ping)
}

func (s *assertingPingService) PingList(
	ping *testpb.PingListRequest,
	stream testpb.TestService_PingListServer,
) error {
	return s.TestServiceServer.PingList(ping, stream)
}

func TestPermsTestSuite(t *testing.T) {
	s := &PermissionTestSuite{
		InterceptorTestSuite: &testpb.InterceptorTestSuite{
			TestService: &assertingPingService{&testpb.TestPingService{}, t},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(
					StreamServerInterceptor(buildDummyStreamPermsFunction(t, false)),
				),
				grpc.UnaryInterceptor(
					UnaryServerInterceptor(buildDummyUnaryPermsFunction(t, false)),
				),
			},
		},
	}
	suite.Run(t, s)
}

type PermissionTestSuite struct {
	*testpb.InterceptorTestSuite
}

func (s *PermissionTestSuite) TestUnary_NoFail() {
	ctx := s.SimpleCtx()
	_, err := s.Client.Ping(ctx, testpb.GoodPing)
	s.Require().Error(err, "there must be an error")
	s.Equal(codes.Internal, status.Code(err), "must error with internal")
}

func (s *PermissionTestSuite) TestUnary_NoPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), true)
	_, err := s.Client.Ping(ctx, testpb.GoodPing)
	s.Require().Error(err, "there must be an error")
	s.Equal(
		codes.PermissionDenied,
		status.Code(err),
		"must error with permission denied",
	)
}

func (s *PermissionTestSuite) TestUnary_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	_, err := s.Client.Ping(ctx, testpb.GoodPing)
	s.Require().NoError(err, "no error must occur")
}

func (s *PermissionTestSuite) TestStream_NoFail() {
	ctx := s.SimpleCtx()
	stream, err := s.Client.PingList(ctx, testpb.GoodPingList)
	s.Require().NoError(err, "should not fail on establishing the stream")
	_, err = stream.Recv()
	s.Require().Error(err, "there must be an error")
	s.Equal(codes.Internal, status.Code(err), "must error with internal")
}

func (s *PermissionTestSuite) TestStream_NoPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), true)
	stream, err := s.Client.PingList(ctx, testpb.GoodPingList)
	s.Require().NoError(err, "should not fail on establishing the stream")
	_, err = stream.Recv()
	s.Require().Error(err, "there must be an error")
	s.Equal(
		codes.PermissionDenied,
		status.Code(err),
		"must error with permission denied",
	)
}

func (s *PermissionTestSuite) TestStream_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	stream, err := s.Client.PingList(ctx, testpb.GoodPingList)
	s.Require().NoError(err, "should not fail on establishing the stream")
	pong, err := stream.Recv()
	s.Require().NoError(err, "no error must occur")
	s.Require().NotNil(pong, "pong must not be nil")
}

type permsOverrideTestService struct {
	testpb.TestServiceServer

	T *testing.T
}

func (s *permsOverrideTestService) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	assert.NotEmpty(s.T, fullMethodName, "method name of caller is passed around")
	return buildDummyUnaryPermsFunction(s.T, false)(ctx, &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: fullMethodName,
	})
}

func TestPermsOverrideTestSuite(t *testing.T) {
	s := &PermsOverrideTestSuite{
		InterceptorTestSuite: &testpb.InterceptorTestSuite{
			TestService: &permsOverrideTestService{
				&assertingPingService{&testpb.TestPingService{}, t},
				t,
			},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(
					StreamServerInterceptor(buildDummyStreamPermsFunction(t, false)),
				),
				grpc.UnaryInterceptor(
					UnaryServerInterceptor(buildDummyUnaryPermsFunction(t, false)),
				),
			},
		},
	}
	suite.Run(t, s)
}

type PermsOverrideTestSuite struct {
	*testpb.InterceptorTestSuite
}

func (s *PermsOverrideTestSuite) TestUnary_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	_, err := s.Client.Ping(ctx, testpb.GoodPing)
	s.Require().NoError(err, "no error must occur")
}

func (s *PermsOverrideTestSuite) TestStream_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	stream, err := s.Client.PingList(ctx, testpb.GoodPingList)
	s.Require().NoError(err, "should not fail on establishing the stream")
	pong, err := stream.Recv()
	s.Require().NoError(err, "no error must occur")
	s.Require().NotNil(pong, "pong must not be nil")
}

type permsRemapTestService struct {
	testpb.TestServiceServer

	T *testing.T
}

func (s *permsRemapTestService) GetPermsRemap() map[string]string {
	return map[string]string{
		"": "",
	}
}

func TestPermsRemapTestSuite(t *testing.T) {
	s := &PermsRemapTestSuite{
		InterceptorTestSuite: &testpb.InterceptorTestSuite{
			TestService: &permsRemapTestService{
				&assertingPingService{&testpb.TestPingService{}, t},
				t,
			},
			ServerOpts: []grpc.ServerOption{
				grpc.StreamInterceptor(
					StreamServerInterceptor(buildDummyStreamPermsFunction(t, true)),
				),
				grpc.UnaryInterceptor(
					UnaryServerInterceptor(buildDummyUnaryPermsFunction(t, true)),
				),
			},
		},
	}
	suite.Run(t, s)
}

type PermsRemapTestSuite struct {
	*testpb.InterceptorTestSuite
}

func (s *PermsRemapTestSuite) TestUnary_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	_, err := s.Client.Ping(ctx, testpb.GoodPing)
	s.Require().NoError(err, "no error must occur")
}

func (s *PermsRemapTestSuite) TestStream_HasPerms() {
	ctx := ctxWithFail(s.SimpleCtx(), false)
	stream, err := s.Client.PingList(ctx, testpb.GoodPingList)
	s.Require().NoError(err, "should not fail on establishing the stream")
	pong, err := stream.Recv()
	s.Require().NoError(err, "no error must occur")
	s.Require().NotNil(pong, "pong must not be nil")
}
