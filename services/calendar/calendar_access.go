package calendar

import (
	calendaraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/calendar/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
)

var calendarSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(calendaraccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_SHARE),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_EDIT),
		int32(calendaraccess.AccessLevel_ACCESS_LEVEL_MANAGE),
	},
}
