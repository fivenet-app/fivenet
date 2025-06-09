package centrumstate

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
)

func (s *State) GetSettings(ctx context.Context, job string) (*centrum.Settings, error) {
	settings, err := s.settings.GetOrLoad(ctx, job)
	if err != nil {
		return nil, err
	}
	settings.Default(job)

	return settings, nil
}

func (s *State) UpdateSettings(ctx context.Context, job string, in *centrum.Settings) error {
	current, err := s.GetSettings(ctx, job)
	if err != nil {
		return err
	}
	current.Merge(in)

	if err := s.settings.Put(ctx, job, current); err != nil {
		return err
	}

	return nil
}

func (s *State) ListSettings(ctx context.Context) []*centrum.Settings {
	list := []*centrum.Settings{}

	s.settings.Range(ctx, func(_ string, settings *centrum.Settings) bool {
		list = append(list, settings)
		return true
	})

	return list
}

func (s *State) ListSettingsFunc(ctx context.Context, fn func(*centrum.Settings) bool) []*centrum.Settings {
	list := []*centrum.Settings{}

	s.settings.Range(ctx, func(_ string, settings *centrum.Settings) bool {
		if fn(settings) {
			list = append(list, settings)
		}
		return true
	})

	return list
}
