package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	"google.golang.org/protobuf/proto"
)

func (s *State) GetSettings(ctx context.Context, job string) *centrum.Settings {
	settings, _ := s.settings.LoadOrCompute(job, func() *centrum.Settings {
		return &centrum.Settings{
			Job:              job,
			Enabled:          false,
			Mode:             centrum.CentrumMode_CENTRUM_MODE_MANUAL,
			FallbackMode:     centrum.CentrumMode_CENTRUM_MODE_MANUAL,
			PredefinedStatus: &centrum.PredefinedStatus{},
			Timings: &centrum.Timings{
				DispatchMaxWait:            900,
				RequireUnit:                false,
				RequireUnitReminderSeconds: 180,
			},
		}
	})

	return proto.Clone(settings).(*centrum.Settings)
}

func (s *State) UpdateSettings(ctx context.Context, job string, in *centrum.Settings) error {
	current := s.GetSettings(ctx, job)
	// Simply use protobuf merge to update existing settings with incoming settings
	proto.Merge(current, in)

	s.settings.Store(job, current)

	return nil
}

func (s *State) ListSettings(ctx context.Context) []*centrum.Settings {
	list := []*centrum.Settings{}

	s.settings.Range(func(_ string, settings *centrum.Settings) bool {
		list = append(list, settings)
		return true
	})

	return list
}
