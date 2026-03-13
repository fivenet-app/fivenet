package stats

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
)

type DocumentMetricExtractor interface {
	SourceKey() string
	Supports(doc *documents.Document) bool
	Extract(ctx context.Context, doc *documents.Document) ([]*DocumentMetric, error)
}
