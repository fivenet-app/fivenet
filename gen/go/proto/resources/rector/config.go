package rector

import (
	"database/sql/driver"
	"time"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	DefaultUserTrackerRefreshTime   = 3*time.Second + 250*time.Millisecond
	DefaultUserTrackerDbRefreshTime = 1 * time.Second

	DefaultDiscordSyncInterval = 15 * time.Minute
)

func (x *AppConfig) Default() {
	if x.Auth == nil {
		x.Auth = &Auth{
			SignupEnabled: true,
			LastCharLock:  false,
		}
	}

	if x.Perms == nil {
		x.Perms = &Perms{
			Default: []*Perm{
				{
					Category: "AuthService",
					Name:     "ChooseCharacter",
				},
				{
					Category: "CompletorService",
					Name:     "CompleteJobs",
				},
				{
					Category: "DocStoreService",
					Name:     "ListDocuments",
				},
				{
					Category: "DocStoreService",
					Name:     "GetDocument",
				},
				{
					Category: "DocStoreService",
					Name:     "PostComment",
				},
			},
		}
	}

	if x.Website == nil {
		x.Website = &Website{
			Links:     &Links{},
			StatsPage: false,
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
	if x.Discord.SyncInterval == nil {
		x.Discord.SyncInterval = durationpb.New(DefaultDiscordSyncInterval)
	}
}

// Scan implements driver.Valuer for protobuf AppConfig.
func (x *AppConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protoutils.UnmarshalPartial([]byte(t), x)
	case []byte:
		return protoutils.UnmarshalPartial(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *AppConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
