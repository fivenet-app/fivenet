package stats

import (
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/require"
)

func TestPenaltyCalculatorExtractor_WithTotals(t *testing.T) {
	extractor := NewPenaltyCalculatorExtractor()
	createdAt := time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)

	doc := &documents.Document{
		Id:         10,
		CreatorJob: "police",
		CreatedAt:  timestamp.New(createdAt),
		Data: &documentsdata.DocumentData{
			PenaltyCalculator: &documentsdata.PenaltyCalculatorData{
				Reduction: 15,
				Selected: []*documentsdata.SelectedPenalty{
					{LawId: 5, Count: 2},
					{LawId: 7, Count: 1},
				},
				Total: &documentsdata.PenaltyCalculatorTotal{
					Count:         ptrUint32(9),
					Fine:          ptrUint32(2000),
					DetentionTime: ptrUint32(60),
					StvoPoints:    ptrUint32(4),
				},
			},
		},
	}

	metrics, err := extractor.Extract(t.Context(), doc)
	require.NoError(t, err)

	actual := toMetricMap(metrics)
	require.Equal(t, int64(1), actual["case_count|"])
	require.Equal(t, int64(15), actual["reduction_percent|"])
	require.Equal(t, int64(2), actual["selected_law_distinct_count|"])
	require.Equal(t, int64(9), actual["selected_law_total_count|"])
	require.Equal(t, int64(2000), actual["fine_total|"])
	require.Equal(t, int64(60), actual["detention_time_total|"])
	require.Equal(t, int64(4), actual["stvo_points_total|"])
	require.Equal(t, int64(2), actual["law_count|5"])
	require.Equal(t, int64(1), actual["law_count|7"])

	for _, m := range metrics {
		require.Equal(t, createdAt, m.OccurredAt)
		require.Equal(t, PenaltyCalculatorSourceKey, m.SourceKey)
		require.Equal(t, int64(10), m.DocumentID)
		require.Equal(t, "police", m.Job)
	}
}

func TestPenaltyCalculatorExtractor_FallbackTotalCount(t *testing.T) {
	extractor := NewPenaltyCalculatorExtractor()
	doc := &documents.Document{
		Id:         11,
		CreatorJob: "police",
		CreatedAt:  timestamp.New(time.Date(2026, 3, 12, 10, 0, 0, 0, time.UTC)),
		Data: &documentsdata.DocumentData{
			PenaltyCalculator: &documentsdata.PenaltyCalculatorData{
				Reduction: 20,
				Selected: []*documentsdata.SelectedPenalty{
					{LawId: 10, Count: 2},
					{LawId: 11, Count: 3},
				},
				Total: &documentsdata.PenaltyCalculatorTotal{},
			},
		},
	}

	metrics, err := extractor.Extract(t.Context(), doc)
	require.NoError(t, err)

	actual := toMetricMap(metrics)
	require.Equal(t, int64(5), actual["selected_law_total_count|"])
	_, hasFine := actual["fine_total|"]
	require.False(t, hasFine)
}

func TestPenaltyCalculatorExtractor_Supports(t *testing.T) {
	extractor := NewPenaltyCalculatorExtractor()
	require.False(t, extractor.Supports(&documents.Document{}))
	require.True(
		t,
		extractor.Supports(
			&documents.Document{
				Data: &documentsdata.DocumentData{
					PenaltyCalculator: &documentsdata.PenaltyCalculatorData{},
				},
			},
		),
	)
}

func ptrUint32(v uint32) *uint32 {
	return &v
}

func toMetricMap(items []*DocumentMetric) map[string]int64 {
	out := map[string]int64{}
	for _, metric := range items {
		key := metric.MetricKey + "|"
		if metric.Dimension1 != nil {
			key += *metric.Dimension1
		}
		out[key] = metric.Value
	}
	return out
}
