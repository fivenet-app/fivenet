package auth

import (
	"context"

	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
)

const (
	auditFailurePrefix = "fivenet.auth.failure"
)

func auditAuthFailure(ctx context.Context, operation, reason string, attrs map[string]string) {
	grpc_audit.AddMeta(ctx, auditFailurePrefix+".operation", operation)
	grpc_audit.AddMeta(ctx, auditFailurePrefix+".reason", reason)

	for k, v := range attrs {
		if v == "" {
			continue
		}
		grpc_audit.AddMeta(ctx, auditFailurePrefix+"."+k, v)
	}
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
