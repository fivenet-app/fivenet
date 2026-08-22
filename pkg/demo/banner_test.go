package demo

import (
	"fmt"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDemoBannerMessageIncludesCredentialsAndIsIndefinite(t *testing.T) {
	t.Parallel()

	d := newTestDemo(7)
	msg := d.buildDemoBannerMessage()

	require.NotNil(t, msg)
	assert.Equal(
		t,
		fmt.Sprintf(
			"<p><strong>Demo credentials</strong>: username <code>%s</code>, password <code>%s</code>.</p>",
			demoAccountUsername,
			demoAccountPassword,
		),
		msg.GetTitle(),
	)
	assert.Equal(
		t,
		utils.GetSHA256HashFromString(msg.GetTitle()+"-"+time.Time{}.String()),
		msg.GetId(),
	)
	assert.Nil(t, msg.GetExpiresAt(), "expected demo banner to remain visible indefinitely")
}
