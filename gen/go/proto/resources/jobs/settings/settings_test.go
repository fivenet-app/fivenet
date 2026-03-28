package jobssettings

import (
	"strconv"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestDiscordSyncChanges(t *testing.T) {
	t.Parallel()
	c := &DiscordSyncChanges{}

	ts := timestamp.Now()
	for i := 1; i <= 12; i++ {
		c.Add(&DiscordSyncChange{
			Plan: strconv.Itoa(i),
			Time: ts,
		})
	}

	assert.Len(t, c.GetChanges(), 12)

	c.Add(&DiscordSyncChange{
		Plan: "13",
		Time: ts,
	})

	assert.Len(t, c.GetChanges(), 12)
}
