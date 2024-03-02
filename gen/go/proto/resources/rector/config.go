package rector

import (
	"database/sql/driver"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultUserTrackerRefreshTime   = 3*time.Second + 350*time.Millisecond
	DefaultUserTrackerDbRefreshTime = 1 * time.Second
)

func (x *AppConfig) Default() {
	if x.Auth == nil {
		x.Auth = &Auth{
			SignupEnabled: true,
		}
	}

	if x.Website == nil {
		x.Website = &Website{
			Links: &Links{},
		}
	}
	if x.Website.Links == nil {
		x.Website.Links = &Links{}
	}

	if x.JobInfo == nil {
		x.JobInfo = &JobInfo{}
	}

	if x.UserTracker == nil {
		x.UserTracker = &UserTracker{
			LivemapJobs:   []string{},
			TimeclockJobs: []string{},
		}
	}
	if x.UserTracker.RefreshTime == nil {
		x.UserTracker.RefreshTime = durationpb.New(DefaultUserTrackerRefreshTime)
	}
	if x.UserTracker.DbRefreshTime == nil {
		x.UserTracker.DbRefreshTime = durationpb.New(DefaultUserTrackerDbRefreshTime)
	}

	if x.Oauth2 == nil {
		x.Oauth2 = &OAuth2{
			Providers: []*OAuth2Provider{},
		}
	}

	if x.Discord == nil {
		x.Discord = &Discord{
			Enabled: false,
		}
	}
}

func (x *PluginConfig) Default() {
	// TODO
}

func (x *AppConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf AppConfig.
func (x *AppConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return string(out), err
}

func (x *PluginConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf PluginConfig.
func (x *PluginConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return string(out), err
}
