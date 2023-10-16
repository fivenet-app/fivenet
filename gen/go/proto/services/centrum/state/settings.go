package state

import "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"

func (s *State) GetSettings(job string) *dispatch.Settings {
	settings, ok := s.Settings.Load(job)
	if !ok {
		// Return default settings
		return &dispatch.Settings{
			Job:          job,
			Enabled:      false,
			Mode:         dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
			FallbackMode: dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
		}
	}

	return settings
}
