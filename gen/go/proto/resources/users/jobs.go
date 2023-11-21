package users

import (
	"database/sql/driver"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const DefaultLivemapMarkerColor = "5C7AFF"

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}

func (x *JobProps) Default(job string) {
	if x.Job == "" {
		x.Job = job
	}

	if x.Theme == "" {
		x.Theme = "default"
	}

	if x.LivemapMarkerColor == "" {
		x.LivemapMarkerColor = DefaultLivemapMarkerColor
	}

	if x.DiscordSyncSettings == nil {
		x.DiscordSyncSettings = &DiscordSyncSettings{
			UserInfoSync: true,
		}
	}
}

func (x *QuickButtons) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf QuickButtons.
func (x *QuickButtons) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}

func (x *DiscordSyncSettings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf DiscordSyncSettings.
func (x *DiscordSyncSettings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}
