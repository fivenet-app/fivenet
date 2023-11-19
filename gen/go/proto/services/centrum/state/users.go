package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
)

func (s *State) GetUserUnitID(userId int32) (uint64, bool) {
	mapping, err := s.userIDToUnitID.Load(UserIdKey(userId))
	if mapping == nil || err != nil {
		return 0, false
	}

	return mapping.UnitId, true
}

func (s *State) SetUnitForUser(userId int32, unitId uint64) error {
	mapping := &centrum.UserUnitMapping{
		UnitId:    unitId,
		UserId:    userId,
		CreatedAt: timestamp.Now(),
	}

	if err := s.userIDToUnitID.Put(UserIdKey(userId), mapping); err != nil {
		return err
	}

	return nil
}

func (s *State) UnsetUnitIDForUser(userId int32) error {
	return s.userIDToUnitID.Delete(UserIdKey(userId))
}
