package stats

import (
	"context"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
)

const PenaltyCalculatorSourceKey = "penalty_calculator"

type PenaltyCalculatorExtractor struct{}

func NewPenaltyCalculatorExtractor() *PenaltyCalculatorExtractor {
	return &PenaltyCalculatorExtractor{}
}

func (e *PenaltyCalculatorExtractor) SourceKey() string {
	return PenaltyCalculatorSourceKey
}

func (e *PenaltyCalculatorExtractor) Supports(doc *documents.Document) bool {
	return doc != nil && doc.GetData() != nil && doc.GetData().GetPenaltyCalculator() != nil
}

func (e *PenaltyCalculatorExtractor) Extract(
	_ context.Context,
	doc *documents.Document,
) ([]*DocumentMetric, error) {
	if !e.Supports(doc) {
		return nil, nil
	}

	data := doc.GetData().GetPenaltyCalculator()
	selected := data.GetSelected()
	metrics := make([]*DocumentMetric, 0, 8+len(selected))
	occurredAt := time.Now().UTC()
	if doc.GetCreatedAt() != nil {
		occurredAt = doc.GetCreatedAt().AsTime()
	}

	addMetric := func(metricKey string, value int64, dimension1 *string) {
		metrics = append(metrics, &DocumentMetric{
			DocumentID: doc.GetId(),
			Job:        doc.GetCreatorJob(),
			SourceKey:  PenaltyCalculatorSourceKey,
			MetricKey:  metricKey,
			Dimension1: dimension1,
			Value:      value,
			OccurredAt: occurredAt,
		})
	}

	addMetric("case_count", 1, nil)
	addMetric("reduction_percent", int64(data.GetReduction()), nil)
	addMetric("selected_law_distinct_count", int64(len(selected)), nil)

	totalCount := int64(0)
	if total := data.GetTotal(); total != nil && total.HasCount() {
		totalCount = int64(total.GetCount())
	} else {
		for _, law := range selected {
			totalCount += int64(law.GetCount())
		}
	}
	addMetric("selected_law_total_count", totalCount, nil)

	if total := data.GetTotal(); total != nil {
		if total.HasFine() {
			addMetric("fine_total", int64(total.GetFine()), nil)
		}
		if total.HasDetentionTime() {
			addMetric("detention_time_total", int64(total.GetDetentionTime()), nil)
		}
		if total.HasStvoPoints() {
			addMetric("stvo_points_total", int64(total.GetStvoPoints()), nil)
		}
	}

	for _, law := range selected {
		dimension1 := strconv.FormatInt(law.GetLawId(), 10)
		addMetric("law_count", int64(law.GetCount()), &dimension1)
	}

	return metrics, nil
}
