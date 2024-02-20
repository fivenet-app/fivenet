package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
)

func (s *State) GetUserUnitID(userId int32) (uint64, bool) {
	mapping, err := s.userIDToUnitID.Load(userIdKey(userId))
	if mapping == nil || err != nil {
		return 0, false
	}

	return mapping.UnitId, true
}

func (s *State) SetUnitForUser(job string, userId int32, unitId uint64) error {
	mapping := &centrum.UserUnitMapping{
		UnitId:    unitId,
		UserId:    userId,
		Job:       job,
		CreatedAt: timestamp.Now(),
	}

	if err := s.userIDToUnitID.Put(userIdKey(userId), mapping); err != nil {
		return err
	}

	return nil
}

func (s *State) UnsetUnitIDForUser(userId int32) error {
	return s.userIDToUnitID.Delete(userIdKey(userId))
}
