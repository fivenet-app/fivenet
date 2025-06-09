package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
)

func (s *State) GetUserUnitMapping(ctx context.Context, userId int32) (*centrum.UserUnitMapping, bool) {
	mapping, err := s.userIDToUnitID.Load(ctx, userIdKey(userId))
	if mapping == nil || err != nil {
		return nil, false
	}

	return mapping, true
}

func (s *State) SetUnitForUser(ctx context.Context, job string, userId int32, unitId uint64) error {
	mapping := &centrum.UserUnitMapping{
		UnitId:    unitId,
		UserId:    userId,
		Job:       job,
		CreatedAt: timestamp.Now(),
	}

	if err := s.userIDToUnitID.Put(ctx, userIdKey(userId), mapping); err != nil {
		return err
	}

	return nil
}

func (s *State) UnsetUnitIDForUser(ctx context.Context, userId int32) error {
	return s.userIDToUnitID.Delete(ctx, userIdKey(userId))
}

func (s *State) ListUserUnitMappings(ctx context.Context) (map[int32]*centrum.UserUnitMapping, error) {
	mappings := s.userIDToUnitID.List()
	ids := map[int32]*centrum.UserUnitMapping{}
	for _, m := range mappings {
		ids[m.UserId] = m
	}

	return ids, nil
}
