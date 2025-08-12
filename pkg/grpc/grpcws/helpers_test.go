// Copyright 2017 Improbable. All Rights Reserved.
// See LICENSE for licensing terms.

package grpcws_test

import (
	"sort"
	"testing"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	testproto "github.com/improbable-eng/grpc-web/integration_test/go/_proto/improbable/grpcweb/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestListGRPCResources(t *testing.T) {
	server := grpc.NewServer()
	testproto.RegisterTestServiceServer(server, &testServiceImpl{})
	expected := []string{
		"/improbable.grpcweb.test.TestService/PingEmpty",
		"/improbable.grpcweb.test.TestService/Ping",
		"/improbable.grpcweb.test.TestService/PingError",
		"/improbable.grpcweb.test.TestService/PingList",
		"/improbable.grpcweb.test.TestService/Echo",
		"/improbable.grpcweb.test.TestService/PingPongBidi",
		"/improbable.grpcweb.test.TestService/PingStream",
	}
	actual := grpcweb.ListGRPCResources(server)
	sort.Strings(expected)
	sort.Strings(actual)
	assert.Equal(t,
		expected,
		actual,
		"list grpc resources must provide an exhaustive list of all registered handlers")
}
