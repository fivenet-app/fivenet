package state

func (s *State) GetUserUnitID(userId int32) (uint64, bool) {
	return s.userIDToUnitID.Load(userId)
}

func (s *State) SetUnitForUser(userId int32, unitId uint64) {
	s.userIDToUnitID.Store(userId, unitId)
}

func (s *State) UnsetUnitForUser(userId int32) {
	s.userIDToUnitID.Delete(userId)
}
