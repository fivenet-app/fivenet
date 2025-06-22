package settings

import "time"

const DispatchExpirationTime = 30 * time.Second

func (s *SettingsDB) DispatchAssignmentExpirationTime() time.Time {
	return time.Now().Add(DispatchExpirationTime)
}
