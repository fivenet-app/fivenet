package grpc_audit

import (
	"context"
	"testing"
	"time"

	audit "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	filepb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type testLogger struct {
	entry *audit.AuditEntry
	req   any
}

func TestAddFileMeta(t *testing.T) {
	t.Parallel()

	entry := &audit.AuditEntry{Meta: &audit.AuditEntryMeta{}}
	ctx := withEntry(t.Context(), newEntry(t.Context(), entry))
	AddFileMeta(ctx, &filepb.File{
		Id:          42,
		FilePath:    "citizen_files/42.png",
		ContentType: "image/png",
		ByteSize:    1234,
	})
	AddPreviousFileMeta(ctx, 17, "citizen_files/old.png")

	require.Equal(t, map[string]string{
		"file.id":                   "42",
		"file.storage_key":          "citizen_files/42.png",
		"file.content_type":         "image/png",
		"file.byte_size":            "1234",
		"file.previous_id":          "17",
		"file.previous_storage_key": "citizen_files/old.png",
	}, entry.GetMeta().GetMeta())
}

func (l *testLogger) Log(entry *audit.AuditEntry, req any) {
	l.entry = entry
	l.req = req
}

func TestUnaryLogsAllowedCallWithoutUserInfo(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	interceptor := NewUnary(Options{
		Logger: logger,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
	})

	req := &struct {
		Username string
	}{
		Username: "new-user",
	}

	resp, err := interceptor(
		t.Context(),
		req,
		&grpc.UnaryServerInfo{FullMethod: "/services.auth.AuthService/CreateAccount"},
		func(ctx context.Context, req any) (any, error) {
			SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
			SetAccountID(ctx, 42)
			return "ok", nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)

	require.NotNil(t, logger.entry)
	require.Equal(t, "services.auth.AuthService", logger.entry.GetService())
	require.Equal(t, "CreateAccount", logger.entry.GetMethod())
	require.False(t, logger.entry.HasUserId())
	require.Equal(t, int64(42), logger.entry.GetAccountId())
	require.Equal(t, audit.EventAction_EVENT_ACTION_CREATED, logger.entry.GetAction())
	require.Equal(t, audit.EventResult_EVENT_RESULT_SUCCEEDED, logger.entry.GetResult())
	require.Same(t, req, logger.req)
}

func TestUnarySkipsAPITokenCallWithoutUserInfo(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	interceptor := NewUnary(Options{
		Logger: logger,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
	})

	resp, err := interceptor(
		auth.ContextWithAuthKind(t.Context(), auth.AuthKindAPIToken),
		&struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/services.example.ExampleService/Sync"},
		func(ctx context.Context, req any) (any, error) {
			require.True(t, IsSkip(ctx))
			return "ok", nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.Nil(t, logger.entry)
}

func TestUnaryLogsAPITokenCallWithAccountUserInfo(t *testing.T) {
	t.Parallel()

	logger := &testLogger{}
	interceptor := NewUnary(Options{
		Logger: logger,
		Now: func() time.Time {
			return time.Unix(0, 0)
		},
	})

	ctx := auth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		AccountId: 42,
	})
	ctx = auth.ContextWithAuthKind(ctx, auth.AuthKindAPIToken)

	resp, err := interceptor(
		ctx,
		&struct{}{},
		&grpc.UnaryServerInfo{FullMethod: "/services.example.ExampleService/Sync"},
		func(ctx context.Context, req any) (any, error) {
			require.False(t, IsSkip(ctx))
			return "ok", nil
		},
	)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)

	require.NotNil(t, logger.entry)
	require.Equal(t, "services.example.ExampleService", logger.entry.GetService())
	require.Equal(t, "Sync", logger.entry.GetMethod())
	require.Equal(t, int64(42), logger.entry.GetAccountId())
}
