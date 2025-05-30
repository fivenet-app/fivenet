// Copyright 2017 Improbable. All Rights Reserved.
// See LICENSE for licensing terms.

package grpcws

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"google.golang.org/grpc"
)

var pathMatcher = regexp.MustCompile(`/[^/]*/[^/]*$`)

// ListGRPCResources is a helper function that lists all URLs that are registered on gRPC server.
//
// This makes it easy to register all the relevant routes in your HTTP router of choice.
func ListGRPCResources(server *grpc.Server) []string {
	ret := []string{}
	for serviceName, serviceInfo := range server.GetServiceInfo() {
		for _, methodInfo := range serviceInfo.Methods {
			fullResource := fmt.Sprintf("/%s/%s", serviceName, methodInfo.Name)
			ret = append(ret, fullResource)
		}
	}
	return ret
}

// WebsocketRequestOrigin returns the host from which a websocket request made by a web browser
// originated.
func WebsocketRequestOrigin(req *http.Request) (string, error) {
	origin := req.Header.Get("Origin")
	parsed, err := url.ParseRequestURI(origin)
	if err != nil {
		return "", fmt.Errorf("failed to parse url for grpc-websocket origin check: %q. error: %v", origin, err)
	}
	return parsed.Host, nil
}

func getGRPCEndpoint(req *http.Request) string {
	endpoint := pathMatcher.FindString(strings.TrimRight(req.URL.Path, "/"))
	if len(endpoint) == 0 {
		return req.URL.Path
	}

	return endpoint
}

// IsGrpcWebSocketRequest determines if a request is a gRPC-Web request by checking that the "Upgrade" header is set and
// "Sec-Websocket-Protocol" header value is "grpc-websocket-channel" and that the "root" path is requested
func IsGrpcWebSocketChannelRequest(req *http.Request) bool {
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
func IsGrpcWebRequest(req *http.Request) bool {
	return req.Method == http.MethodPost && strings.HasPrefix(req.Header.Get("content-type"), grpcWebContentType)
}
