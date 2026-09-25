package croner

import (
	"strings"

	"github.com/fivenet-app/fivenet/v2026/pkg/events"
)

func normalizeCronjobName(name string) string {
	return strings.ToLower(events.SanitizeKey(name))
}
