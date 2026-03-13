package stats

import "time"

type DocumentMetric struct {
	DocumentID int64

	Job       string
	SourceKey string
	MetricKey string

	Dimension1 *string
	Dimension2 *string
	Dimension3 *string

	Value int64

	OccurredAt time.Time
}

type KeyValue struct {
	Key   string `alias:"key"`
	Value int64  `alias:"value"`
}

type CategoryValue struct {
	ID    int64   `alias:"id"`
	Name  string  `alias:"name"`
	Color *string `alias:"color"`
	Icon  *string `alias:"icon"`
	Value int64   `alias:"value"`
}

type DailyValue struct {
	Day   time.Time `alias:"day"`
	Value int64     `alias:"value"`
}

type PeriodSeriesValue struct {
	Day   time.Time `alias:"day"`
	Key   string    `alias:"key"`
	Label string    `alias:"label"`
	Value int64     `alias:"value"`
}
