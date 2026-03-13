package stats

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	SourceKindDocumentColumn = "document_column"
	SourceKindDocumentMetric = "document_metric"
)

type Service struct {
	mu sync.Mutex

	db         *sql.DB
	extractors []DocumentMetricExtractor
}

func NewService(db *sql.DB, extractors ...DocumentMetricExtractor) *Service {
	return &Service{
		db:         db,
		extractors: extractors,
	}
}

func (s *Service) RebuildDocumentMetrics(ctx context.Context, doc *documents.Document) error {
	if doc == nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.RebuildDocumentMetricsTx(ctx, tx, doc); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) RebuildDocumentMetricsTx(
	ctx context.Context,
	tx qrm.DB,
	doc *documents.Document,
) error {
	if doc == nil {
		return nil
	}

	if !isPublishedDocument(doc) {
		return s.clearDocumentMetrics(ctx, tx, doc.GetId())
	}

	grouped := map[string][]*DocumentMetric{}
	for _, extractor := range s.extractors {
		source := extractor.SourceKey()
		grouped[source] = []*DocumentMetric{}

		if !extractor.Supports(doc) {
			continue
		}

		metrics, err := extractor.Extract(ctx, doc)
		if err != nil {
			return err
		}
		for _, metric := range metrics {
			if metric == nil {
				continue
			}
			grouped[source] = append(grouped[source], metric)
		}
	}

	sources := make([]string, 0, len(grouped))
	for source := range grouped {
		sources = append(sources, source)
	}
	sort.Strings(sources)

	tMetric := table.FivenetDocumentsStatsMetric
	for _, source := range sources {
		if _, err := tMetric.
			DELETE().
			WHERE(mysql.AND(
				tMetric.DocumentID.EQ(mysql.Int64(doc.GetId())),
				tMetric.SourceKey.EQ(mysql.String(source)),
			)).
			ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	now := time.Now().UTC()
	for _, source := range sources {
		metrics := grouped[source]
		if len(metrics) == 0 {
			continue
		}

		stmt := tMetric.
			INSERT(
				tMetric.DocumentID,
				tMetric.Job,
				tMetric.SourceKey,
				tMetric.MetricKey,
				tMetric.Dimension1,
				tMetric.Dimension2,
				tMetric.Dimension3,
				tMetric.Value,
				tMetric.OccurredAt,
				tMetric.CreatedAt,
				tMetric.UpdatedAt,
			)

		for _, metric := range metrics {
			documentID := metric.DocumentID
			if documentID == 0 {
				documentID = doc.GetId()
			}
			job := metric.Job
			if strings.TrimSpace(job) == "" {
				job = doc.GetCreatorJob()
			}
			sourceKey := metric.SourceKey
			if strings.TrimSpace(sourceKey) == "" {
				sourceKey = source
			}
			occurredAt := metric.OccurredAt
			if occurredAt.IsZero() {
				occurredAt = now
				if doc.GetCreatedAt() != nil {
					occurredAt = doc.GetCreatedAt().AsTime()
				}
			}

			stmt = stmt.VALUES(
				documentID,
				job,
				sourceKey,
				metric.MetricKey,
				dbutils.StringPEmpty(metric.Dimension1),
				dbutils.StringPEmpty(metric.Dimension2),
				dbutils.StringPEmpty(metric.Dimension3),
				metric.Value,
				occurredAt,
				now,
				now,
			)
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) RebuildDocumentColumnRollups(
	ctx context.Context,
	startDay, endDay time.Time,
) error {
	startDay = timeutils.StartOfDay(startDay)
	endDay = timeutils.StartOfDay(endDay)

	return s.rebuildRollupsWithRange(
		ctx,
		startDay,
		endDay,
		SourceKindDocumentColumn,
		`INSERT INTO fivenet_documents_stats_daily_rollup (day, job, source_kind, source_key, metric_key, dimension1, dimension2, dimension3, value)
SELECT DATE(d.created_at), d.creator_job, 'document_column', 'documents', 'document_count', '', '', '', COUNT(*)
FROM fivenet_documents d
WHERE d.deleted_at IS NULL
  AND d.draft = FALSE
  AND DATE(d.created_at) >= ?
  AND DATE(d.created_at) <= ?
GROUP BY DATE(d.created_at), d.creator_job
UNION ALL
SELECT DATE(d.created_at), d.creator_job, 'document_column', 'documents', 'category_count', CAST(d.category_id AS CHAR), '', '', COUNT(*)
FROM fivenet_documents d
WHERE d.deleted_at IS NULL
  AND d.draft = FALSE
  AND d.category_id IS NOT NULL
  AND DATE(d.created_at) >= ?
  AND DATE(d.created_at) <= ?
GROUP BY DATE(d.created_at), d.creator_job, d.category_id
UNION ALL
SELECT DATE(d.created_at), d.creator_job, 'document_column', 'documents', 'template_count', CAST(d.template_id AS CHAR), '', '', COUNT(*)
FROM fivenet_documents d
WHERE d.deleted_at IS NULL
  AND d.draft = FALSE
  AND d.template_id IS NOT NULL
  AND DATE(d.created_at) >= ?
  AND DATE(d.created_at) <= ?
GROUP BY DATE(d.created_at), d.creator_job, d.template_id
UNION ALL
SELECT DATE(d.created_at), d.creator_job, 'document_column', 'documents', 'word_count_sum', '', '', '', SUM(COALESCE(d.word_count, 0))
FROM fivenet_documents d
WHERE d.deleted_at IS NULL
  AND d.draft = FALSE
  AND DATE(d.created_at) >= ?
  AND DATE(d.created_at) <= ?
GROUP BY DATE(d.created_at), d.creator_job`,
		dateRangeArgs(startDay, endDay, 4),
	)
}

func (s *Service) RebuildDocumentMetricRollups(
	ctx context.Context,
	startDay, endDay time.Time,
) error {
	startDay = timeutils.StartOfDay(startDay)
	endDay = timeutils.StartOfDay(endDay)

	return s.rebuildRollupsWithRange(
		ctx,
		startDay,
		endDay,
		SourceKindDocumentMetric,
		`INSERT INTO fivenet_documents_stats_daily_rollup (day, job, source_kind, source_key, metric_key, dimension1, dimension2, dimension3, value)
SELECT DATE(m.occurred_at), m.job, 'document_metric', m.source_key, m.metric_key,
       COALESCE(m.dimension1, ''), COALESCE(m.dimension2, ''), COALESCE(m.dimension3, ''), SUM(m.value)
FROM fivenet_documents_stats_metric m
INNER JOIN fivenet_documents d ON d.id = m.document_id
WHERE d.deleted_at IS NULL
  AND d.draft = FALSE
  AND DATE(m.occurred_at) >= ?
  AND DATE(m.occurred_at) <= ?
GROUP BY DATE(m.occurred_at), m.job, m.source_key, m.metric_key,
         COALESCE(m.dimension1, ''), COALESCE(m.dimension2, ''), COALESCE(m.dimension3, '')`,
		dateRangeArgs(startDay, endDay, 1),
	)
}

func (s *Service) QueryTopLaws(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	limit int64,
) ([]*KeyValue, error) {
	if limit <= 0 {
		limit = 10
	}

	items := []*KeyValue{}
	query := `
SELECT
  CONCAT(COALESCE(b.name, 'Unknown'), '::', COALESCE(l.name, CONCAT('#', x.dimension1))) AS ` + "`key`" + `,
  x.value AS ` + "`value`" + `
FROM (
  SELECT
    r.dimension1,
    SUM(r.value) AS value
  FROM fivenet_documents_stats_daily_rollup r
  WHERE r.day >= ?
    AND r.day <= ?
    AND r.job = ?
    AND r.source_kind = ?
    AND r.source_key = ?
    AND r.metric_key = 'law_count'
    AND r.dimension1 <> ''
  GROUP BY r.dimension1
  ORDER BY SUM(r.value) DESC
  LIMIT ?
) x
LEFT JOIN fivenet_lawbooks_laws l ON l.id = CAST(x.dimension1 AS UNSIGNED)
LEFT JOIN fivenet_lawbooks b ON b.id = l.lawbook_id
ORDER BY x.value DESC
`
	rows, err := s.db.QueryContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		SourceKindDocumentMetric,
		PenaltyCalculatorSourceKey,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &KeyValue{}
		if err := rows.Scan(&item.Key, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) QueryFinesOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	period documentspb.StatsPeriod,
) ([]*DailyValue, error) {
	return s.QueryPeriodValues(
		ctx,
		startDay,
		endDay,
		job,
		SourceKindDocumentMetric,
		PenaltyCalculatorSourceKey,
		"fine_total",
		period,
	)
}

func (s *Service) QueryPenaltySeriesOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	period documentspb.StatsPeriod,
) ([]*PeriodSeriesValue, error) {
	periodExpr := periodStartExpr(period)
	query := fmt.Sprintf(`
SELECT
  %s AS day,
  metric_key AS `+"`key`"+`,
  SUM(value) AS value
FROM fivenet_documents_stats_daily_rollup
WHERE day >= ?
  AND day <= ?
  AND job = ?
  AND source_kind = ?
  AND source_key = ?
  AND metric_key IN ('fine_total', 'detention_time_total', 'stvo_points_total')
GROUP BY %s, metric_key
ORDER BY %s ASC, metric_key ASC
`, periodExpr, periodExpr, periodExpr)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		SourceKindDocumentMetric,
		PenaltyCalculatorSourceKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*PeriodSeriesValue{}
	for rows.Next() {
		item := &PeriodSeriesValue{}
		if err := rows.Scan(&item.Day, &item.Key, &item.Value); err != nil {
			return nil, err
		}

		item.Label = item.Key
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) QueryDocumentsByCategory(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
) ([]*CategoryValue, error) {
	items := []*CategoryValue{}
	query := `
SELECT
  CAST(r.dimension1 AS UNSIGNED) AS ` + "`id`" + `,
  COALESCE(c.name, CONCAT('#', r.dimension1)) AS ` + "`name`" + `,
  c.color AS ` + "`color`" + `,
  c.icon AS ` + "`icon`" + `,
  SUM(r.value) AS ` + "`value`" + `
FROM fivenet_documents_stats_daily_rollup r
LEFT JOIN fivenet_documents_categories c ON c.id = CAST(r.dimension1 AS UNSIGNED) AND c.deleted_at IS NULL
WHERE r.day >= ?
  AND r.day <= ?
  AND r.job = ?
  AND r.source_kind = ?
  AND r.source_key = 'documents'
  AND r.metric_key = 'category_count'
  AND r.dimension1 <> ''
GROUP BY CAST(r.dimension1 AS UNSIGNED), COALESCE(c.name, CONCAT('#', r.dimension1)), c.color, c.icon
ORDER BY SUM(r.value) DESC
`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		SourceKindDocumentColumn,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &CategoryValue{}
		if err := rows.Scan(&item.ID, &item.Name, &item.Color, &item.Icon, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) QueryPenaltyReductionAverage(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
) (int64, int64, error) {
	tRollup := table.FivenetDocumentsStatsDailyRollup
	var sums struct {
		ReductionSum int64 `alias:"reduction_sum"`
		CaseCountSum int64 `alias:"case_count_sum"`
	}

	stmt := tRollup.
		SELECT(
			mysql.AVG(mysql.CASE().
				WHEN(tRollup.MetricKey.EQ(mysql.String("reduction_percent"))).
				THEN(tRollup.Value).
				ELSE(mysql.Int(0))).AS("reduction_sum"),
			mysql.SUM(mysql.CASE().
				WHEN(tRollup.MetricKey.EQ(mysql.String("case_count"))).
				THEN(tRollup.Value).
				ELSE(mysql.Int(0))).AS("case_count_sum"),
		).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.EQ(mysql.String(job)),
			tRollup.SourceKind.EQ(mysql.String(SourceKindDocumentMetric)),
			tRollup.SourceKey.EQ(mysql.String(PenaltyCalculatorSourceKey)),
			tRollup.MetricKey.IN(
				mysql.String("reduction_percent"),
				mysql.String("case_count"),
			),
		))

	if err := stmt.QueryContext(ctx, s.db, &sums); err != nil {
		return 0, 0, err
	}

	return sums.ReductionSum, sums.CaseCountSum, nil
}

func (s *Service) QueryTotalValue(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	sourceKind string,
	sourceKey string,
	metricKey string,
) (int64, error) {
	var total sql.NullInt64
	query := `
SELECT SUM(value)
FROM fivenet_documents_stats_daily_rollup
WHERE day >= ?
  AND day <= ?
  AND job = ?
  AND source_kind = ?
  AND source_key = ?
  AND metric_key = ?
`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		sourceKind,
		sourceKey,
		metricKey,
	).Scan(&total); err != nil {
		return 0, err
	}

	if !total.Valid {
		return 0, nil
	}

	return total.Int64, nil
}

func (s *Service) QueryPeriodValues(
	ctx context.Context,
	startDay, endDay time.Time,
	job string,
	sourceKind string,
	sourceKey string,
	metricKey string,
	period documentspb.StatsPeriod,
) ([]*DailyValue, error) {
	periodExpr := periodStartExpr(period)
	query := fmt.Sprintf(`
SELECT
  %s AS day,
  SUM(value) AS value
FROM fivenet_documents_stats_daily_rollup
WHERE day >= ?
  AND day <= ?
  AND job = ?
  AND source_kind = ?
  AND source_key = ?
  AND metric_key = ?
GROUP BY %s
ORDER BY %s ASC
`, periodExpr, periodExpr, periodExpr)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		timeutils.StartOfDay(startDay).Format(time.DateOnly),
		timeutils.StartOfDay(endDay).Format(time.DateOnly),
		job,
		sourceKind,
		sourceKey,
		metricKey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*DailyValue{}
	for rows.Next() {
		item := &DailyValue{}
		if err := rows.Scan(&item.Day, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) rebuildRollupsWithRange(
	ctx context.Context,
	startDay, endDay time.Time,
	sourceKind string,
	insertSQL string,
	insertArgs []any,
) error {
	startDay = timeutils.StartOfDay(startDay)
	endDay = timeutils.StartOfDay(endDay)
	if endDay.Before(startDay) {
		return fmt.Errorf("end day before start day")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`DELETE FROM fivenet_documents_stats_daily_rollup WHERE day >= ? AND day <= ? AND source_kind = ?`,
		startDay.Format(time.DateOnly),
		endDay.Format(time.DateOnly),
		sourceKind,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, insertSQL, insertArgs...); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) clearDocumentMetrics(ctx context.Context, tx qrm.DB, documentID int64) error {
	tMetric := table.FivenetDocumentsStatsMetric
	_, err := tMetric.
		DELETE().
		WHERE(tMetric.DocumentID.EQ(mysql.Int64(documentID))).
		ExecContext(ctx, tx)
	return err
}

func isPublishedDocument(doc *documents.Document) bool {
	if doc == nil {
		return false
	}
	if doc.GetDeletedAt() != nil {
		return false
	}
	if strings.TrimSpace(doc.GetCreatorJob()) == "" {
		return false
	}
	if doc.GetMeta() == nil {
		return false
	}

	return !doc.GetMeta().GetDraft()
}

func periodStartExpr(period documentspb.StatsPeriod) string {
	switch period {
	case documentspb.StatsPeriod_STATS_PERIOD_MONTHLY:
		return "DATE_SUB(day, INTERVAL DAYOFMONTH(day) - 1 DAY)"

	case documentspb.StatsPeriod_STATS_PERIOD_WEEKLY:
		return "DATE_SUB(day, INTERVAL WEEKDAY(day) DAY)"

	case documentspb.StatsPeriod_STATS_PERIOD_DAILY:
		fallthrough
	default:
		return "day"
	}
}

func dateRangeArgs(startDay, endDay time.Time, repeats int) []any {
	args := make([]any, 0, repeats*2)
	start := startDay.Format(time.DateOnly)
	end := endDay.Format(time.DateOnly)

	for range repeats {
		args = append(args, start, end)
	}

	return args
}
