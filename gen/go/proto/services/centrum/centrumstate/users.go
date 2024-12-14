package centrumstate

import (
	"context"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
)

func (s *State) GetUserUnitID(ctx context.Context, userId int32) (uint64, bool) {
	mapping, err := s.userIDToUnitID.Load(ctx, userIdKey(userId))
	if mapping == nil || err != nil {
		return 0, false
	}

	return mapping.UnitId, true
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

func (s *State) ListUserIdsFromUserIdUnitIds(ctx context.Context) ([]int32, error) {
	keys, err := s.userIDToUnitID.Keys(ctx, "")
	if err != nil {
		return nil, err
	}

	ids := make([]int32, len(keys))
	for idx := range keys {
		id, err := strconv.Atoi(keys[idx])
		if err != nil {
			continue
		}

		ids[idx] = int32(id)
	}

	return ids, nil
}
