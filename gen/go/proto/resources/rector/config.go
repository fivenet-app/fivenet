package rector

import (
	"database/sql/driver"
	"time"

	"github.com/galexrt/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultUserTrackerRefreshTime   = 3*time.Second + 350*time.Millisecond
	DefaultUserTrackerDbRefreshTime = 1 * time.Second
)

func (x *AppConfig) Default() {
	if x.Auth == nil {
		x.Auth = &Auth{
			SignupEnabled: false,
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
		x.JobInfo = &JobInfo{
			PublicJobs: []string{},
			HiddenJobs: []string{},
		}
	}
	if x.JobInfo.UnemployedJob == nil {
		x.JobInfo.UnemployedJob = &UnemployedJob{
			Name:  "unemployed",
			Grade: 1,
		}
	}

	if x.UserTracker == nil {
		x.UserTracker = &UserTracker{
			LivemapJobs: []string{},
		}
	}
	if x.UserTracker.RefreshTime == nil {
		x.UserTracker.RefreshTime = durationpb.New(DefaultUserTrackerRefreshTime)
	}
	if x.UserTracker.DbRefreshTime == nil {
		x.UserTracker.DbRefreshTime = durationpb.New(DefaultUserTrackerDbRefreshTime)
	}

	if x.Discord == nil {
		x.Discord = &Discord{
			Enabled: false,
		}
	}
}

func (x *AppConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protoutils.UnmarshalPartial([]byte(t), x)
	case []byte:
		return protoutils.UnmarshalPartial(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf AppConfig.
func (x *AppConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
