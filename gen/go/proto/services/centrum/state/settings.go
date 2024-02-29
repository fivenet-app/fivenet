package state

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	"google.golang.org/protobuf/proto"
)

func (s *State) GetSettings(job string) *centrum.Settings {
	settings, _ := s.settings.LoadOrCompute(job, func() *centrum.Settings {
		return &centrum.Settings{
			Job:              job,
			Enabled:          false,
			Mode:             centrum.CentrumMode_CENTRUM_MODE_MANUAL,
			FallbackMode:     centrum.CentrumMode_CENTRUM_MODE_MANUAL,
			PredefinedStatus: &centrum.PredefinedStatus{},
		}
	})

	return settings
}

func (s *State) UpdateSettings(job string, in *centrum.Settings) error {
	current := s.GetSettings(job)
	// Simply use protobuf merge to update existing settings with incoming settings
	proto.Merge(current, in)

	return nil
}

func (s *State) ListSettings() []*centrum.Settings {
	list := []*centrum.Settings{}

	s.settings.Range(func(_ string, settings *centrum.Settings) bool {
		list = append(list, settings)
		return true
	})

	return list
}
