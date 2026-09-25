package croner

import (
	"context"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/stretchr/testify/require"
)

func TestNormalizeCronjobName(t *testing.T) {
	t.Parallel()

	require.Equal(t, "my.job_name", normalizeCronjobName("My#Job=Name"))
}

func TestHandlersNormalizeCronjobName(t *testing.T) {
	t.Parallel()

	h := &Handlers{handlers: map[string]CronjobHandlerFn{}}
	h.Add("My#Job=Name", func(context.Context, *cron.CronjobData) error { return nil })
	require.NotNil(t, h.getCronjobHandler("my.job_name"))
}
