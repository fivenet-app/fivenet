package access

import (
	"context"
	"sync"

	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
)

type GroupedAccess interface {
	CanUserAccessTarget(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access int32) (bool, error)
}

type GroupedAccessAdapter struct {
	CanUserAccessTargetFn func(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access int32) (bool, error)
}

func (a *GroupedAccessAdapter) CanUserAccessTarget(ctx context.Context, targetId uint64, userInfo *userinfo.UserInfo, access int32) (bool, error) {
	return a.CanUserAccessTargetFn(ctx, targetId, userInfo, access)
}

var (
	groupedAccessesMu = &sync.RWMutex{}
	groupedAccesses   = map[string]GroupedAccess{}
)

func RegisterAccess(name string, access GroupedAccess) {
	groupedAccessesMu.Lock()
	defer groupedAccessesMu.Unlock()

	if _, exists := groupedAccesses[name]; exists {
		panic("access already registered: " + name)
	}

	groupedAccesses[name] = access
}

func GetAccess(accessLevel string) GroupedAccess {
	groupedAccessesMu.RLock()
	defer groupedAccessesMu.RUnlock()

	access, ok := groupedAccesses[accessLevel]
	if !ok {
		return nil
	}

	return access
}
