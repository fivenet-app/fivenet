package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
		pgs.SupportedFeatures(&optional),
	).RegisterModule(
		Permify(),
		ListSvcMethods(),
	).Render()
}
