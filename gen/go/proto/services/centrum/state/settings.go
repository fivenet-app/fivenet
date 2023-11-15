package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"google.golang.org/protobuf/proto"
)

func (s *State) GetSettings(job string) *dispatch.Settings {
	settings, _ := s.settings.LoadOrCompute(job, func() *dispatch.Settings {
		return &dispatch.Settings{
			Job:          job,
			Enabled:      false,
			Mode:         dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
			FallbackMode: dispatch.CentrumMode_CENTRUM_MODE_MANUAL,
		}
	})

	return settings
}

func (s *State) UpdateSettings(job string, in *dispatch.Settings) error {
	current := s.GetSettings(job)
	// Simply use protobuf merge to update existing settings with incoming settings
	proto.Merge(current, in)

	return nil
}

func (s *State) ListSettings() []*dispatch.Settings {
	list := []*dispatch.Settings{}

	s.settings.Range(func(_ string, settings *dispatch.Settings) bool {
		list = append(list, settings)
		return true
	})

	return list
}
