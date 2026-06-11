package stats

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/housekeeper"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	SourceKindDocumentColumn = "document_column"
	SourceKindDocumentMetric = "document_metric"
	SourceKindEmployeeCount  = "employee_count"
)

func init() {
	housekeeper.AddTable(&housekeeper.Table{
		Table:           table.FivenetDocumentsStatsMetric,
		IDColumn:        table.FivenetDocumentsStatsMetric.ID,
		TimestampColumn: table.FivenetDocumentsStatsMetric.CreatedAt,

		MinDays: 366,
	})

	housekeeper.AddTable(&housekeeper.Table{
		Table:      table.FivenetStatsDailyRollup,
		DateColumn: table.FivenetStatsDailyRollup.Day,

		MinDays: 366,
	})
}

type Service struct {
	mu sync.Mutex

	db     *sql.DB
	appCfg appconfig.IConfig

	extractors []DocumentMetricExtractor
}

func NewService(db *sql.DB, appCfg appconfig.IConfig) *Service {
	return &Service{
		db:     db,
		appCfg: appCfg,

		extractors: []DocumentMetricExtractor{
			NewPenaltyCalculatorExtractor(),
		},
	}
}

func normalizeJobFilters(jobs []string) []string {
	if len(jobs) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(jobs))
	seen := make(map[string]struct{}, len(jobs))

	for _, job := range jobs {
		job = strings.TrimSpace(job)
		if job == "" {
			continue
		}
		if _, ok := seen[job]; ok {
			continue
		}

		seen[job] = struct{}{}
		normalized = append(normalized, job)
	}

	return normalized
}

func jobExpressions(jobs []string) []mysql.Expression {
	expressions := make([]mysql.Expression, 0, len(jobs))
	for _, job := range jobs {
		expressions = append(expressions, mysql.String(job))
	}

	return expressions
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
		deleteCondition := mysql.AND(
			tMetric.DocumentID.EQ(mysql.Int64(doc.GetId())),
			tMetric.SourceKey.EQ(mysql.String(source)),
		)

		if _, err := tMetric.
			DELETE().
			WHERE(deleteCondition).
			LIMIT(10000).
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
		`INSERT INTO fivenet_stats_daily_rollup (day, job, source_kind, source_key, metric_key, dimension1, dimension2, dimension3, value)
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
		`INSERT INTO fivenet_stats_daily_rollup (day, job, source_kind, source_key, metric_key, dimension1, dimension2, dimension3, value)
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
	jobs []string,
	limit int64,
) ([]*KeyValue, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return []*KeyValue{}, nil
	}

	if limit <= 0 {
		limit = 20
	}

	tRollup := table.FivenetStatsDailyRollup.AS("r")
	tLaws := table.FivenetLawbooksLaws.AS("l")
	tLawbooks := table.FivenetLawbooks.AS("b")
	jobNames := jobExpressions(jobs)

	lawNameExpr := mysql.COALESCE(
		tLaws.Name,
		mysql.CONCAT(mysql.String("#"), tRollup.Dimension1),
	)
	lawbookNameExpr := mysql.COALESCE(tLawbooks.Name, mysql.String("Unknown"))
	keyExpr := mysql.CONCAT(lawbookNameExpr, mysql.String("::"), lawNameExpr)

	stmt := tRollup.
		SELECT(
			keyExpr.AS("keyvalue.key"),
			mysql.SUM(tRollup.Value).AS("keyvalue.value"),
		).
		FROM(
			tRollup.
				LEFT_JOIN(
					tLaws,
					tLaws.ID.EQ(mysql.CAST(tRollup.Dimension1).AS_UNSIGNED()),
				).
				LEFT_JOIN(
					tLawbooks,
					tLawbooks.ID.EQ(tLaws.LawbookID),
				),
		).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(SourceKindDocumentMetric)),
			tRollup.SourceKey.EQ(mysql.String(PenaltyCalculatorSourceKey)),
			tRollup.MetricKey.EQ(mysql.String("law_count")),
			tRollup.Dimension1.NOT_EQ(mysql.String("")),
		)).
		GROUP_BY(tRollup.Dimension1, tLawbooks.Name, tLaws.Name).
		ORDER_BY(mysql.SUM(tRollup.Value).DESC()).
		LIMIT(limit)

	items := []*KeyValue{}
	if err := stmt.QueryContext(ctx, s.db, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) QueryFinesOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
	period pbstats.StatsPeriod,
) ([]*DailyValue, error) {
	return s.QueryPeriodValues(
		ctx,
		startDay,
		endDay,
		jobs,
		SourceKindDocumentMetric,
		PenaltyCalculatorSourceKey,
		"fine_total",
		period,
	)
}

func (s *Service) QueryPenaltySeriesOverTime(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
	period pbstats.StatsPeriod,
) ([]*PeriodSeriesValue, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return []*PeriodSeriesValue{}, nil
	}

	tRollup := table.FivenetStatsDailyRollup
	periodExpr := periodStartDateExpr(period)
	jobNames := jobExpressions(jobs)

	stmt := tRollup.
		SELECT(
			periodExpr.AS("periodseriesvalue.day"),
			tRollup.MetricKey.AS("periodseriesvalue.key"),
			mysql.SUM(tRollup.Value).AS("periodseriesvalue.value"),
		).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(SourceKindDocumentMetric)),
			tRollup.SourceKey.EQ(mysql.String(PenaltyCalculatorSourceKey)),
			tRollup.MetricKey.IN(
				mysql.String("fine_total"),
				mysql.String("detention_time_total"),
				mysql.String("stvo_points_total"),
			),
		)).
		GROUP_BY(periodExpr, tRollup.MetricKey).
		ORDER_BY(periodExpr.ASC(), tRollup.MetricKey.ASC())

	rowsDest := []*PeriodSeriesValue{}
	if err := stmt.QueryContext(ctx, s.db, &rowsDest); err != nil {
		return nil, err
	}

	items := make([]*PeriodSeriesValue, 0, len(rowsDest))
	for i := range rowsDest {
		rowsDest[i].Label = rowsDest[i].Key
		items = append(items, rowsDest[i])
	}

	return items, nil
}

func (s *Service) QueryDocumentsByCategory(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
) ([]*CategoryValue, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return []*CategoryValue{}, nil
	}

	tRollup := table.FivenetStatsDailyRollup.AS("r")
	tCategory := table.FivenetDocumentsCategories.AS("c")
	jobNames := jobExpressions(jobs)

	categoryIDExpr := mysql.CAST(tRollup.Dimension1).AS_UNSIGNED()
	categoryNameExpr := mysql.COALESCE(
		mysql.MAX(tCategory.Name),
		mysql.CONCAT(mysql.String("#"), tRollup.Dimension1),
	)
	categoryColorExpr := mysql.MAX(tCategory.Color)
	categoryIconExpr := mysql.MAX(tCategory.Icon)

	stmt := tRollup.
		SELECT(
			categoryIDExpr.AS("category_value.id"),
			categoryNameExpr.AS("category_value.name"),
			categoryColorExpr.AS("category_value.color"),
			categoryIconExpr.AS("category_value.icon"),
			tRollup.Job.AS("category_value.job"),
			mysql.SUM(tRollup.Value).AS("category_value.value"),
		).
		FROM(
			tRollup.LEFT_JOIN(
				tCategory,
				mysql.AND(
					tCategory.ID.EQ(categoryIDExpr),
					tCategory.DeletedAt.IS_NULL(),
				),
			),
		).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(SourceKindDocumentColumn)),
			tRollup.SourceKey.EQ(mysql.String("documents")),
			tRollup.MetricKey.EQ(mysql.String("category_count")),
			tRollup.Dimension1.NOT_EQ(mysql.String("")),
		)).
		GROUP_BY(tRollup.Dimension1, tRollup.Job, categoryIDExpr).
		ORDER_BY(mysql.SUM(tRollup.Value).DESC())

	items := []*CategoryValue{}
	if err := stmt.QueryContext(ctx, s.db, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) QueryPenaltyReductionAverage(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
) (int64, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return 0, nil
	}

	tRollup := table.FivenetStatsDailyRollup
	jobNames := jobExpressions(jobs)

	stmt := tRollup.
		SELECT(
			mysql.AVG(mysql.CASE().
				WHEN(tRollup.MetricKey.EQ(mysql.String("reduction_percent"))).
				THEN(tRollup.Value).
				ELSE(mysql.Int(0))).AS("reduction_sum"),
		).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(SourceKindDocumentMetric)),
			tRollup.SourceKey.EQ(mysql.String(PenaltyCalculatorSourceKey)),
			tRollup.MetricKey.EQ(mysql.String("reduction_percent")),
		))

	var sums struct {
		ReductionSum int64 `alias:"reduction_sum"`
	}
	if err := stmt.QueryContext(ctx, s.db, &sums); err != nil {
		return 0, err
	}

	return sums.ReductionSum, nil
}

func (s *Service) QueryTotalValue(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
	sourceKind string,
	sourceKey string,
	metricKey string,
) (int64, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return 0, nil
	}

	tRollup := table.FivenetStatsDailyRollup
	jobNames := jobExpressions(jobs)

	var dest struct {
		Total sql.NullInt64 `alias:"total"`
	}
	stmt := tRollup.
		SELECT(mysql.SUM(tRollup.Value).AS("total")).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(sourceKind)),
			tRollup.SourceKey.EQ(mysql.String(sourceKey)),
			tRollup.MetricKey.EQ(mysql.String(metricKey)),
		))
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	if !dest.Total.Valid {
		return 0, nil
	}

	return dest.Total.Int64, nil
}

func (s *Service) QueryAverageValue(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
	sourceKind string,
	sourceKey string,
	metricKey string,
) (float64, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return 0, nil
	}

	tRollup := table.FivenetStatsDailyRollup
	jobNames := jobExpressions(jobs)

	var dest struct {
		Average sql.NullFloat64 `alias:"average"`
	}
	stmt := tRollup.
		SELECT(mysql.AVG(tRollup.Value).AS("average")).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(sourceKind)),
			tRollup.SourceKey.EQ(mysql.String(sourceKey)),
			tRollup.MetricKey.EQ(mysql.String(metricKey)),
		))
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	if !dest.Average.Valid {
		return 0, nil
	}

	return dest.Average.Float64, nil
}

func (s *Service) QueryPeriodValues(
	ctx context.Context,
	startDay, endDay time.Time,
	jobs []string,
	sourceKind string,
	sourceKey string,
	metricKey string,
	period pbstats.StatsPeriod,
) ([]*DailyValue, error) {
	jobs = normalizeJobFilters(jobs)
	if len(jobs) == 0 {
		return []*DailyValue{}, nil
	}

	tRollup := table.FivenetStatsDailyRollup
	periodExpr := periodStartDateExpr(period)
	jobNames := jobExpressions(jobs)

	stmt := tRollup.
		SELECT(
			periodExpr.AS("dailyvalue.day"),
			mysql.SUM(tRollup.Value).AS("dailyvalue.value"),
		).
		FROM(tRollup).
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(timeutils.StartOfDay(startDay))),
			tRollup.Day.LT_EQ(mysql.DateT(timeutils.StartOfDay(endDay))),
			tRollup.Job.IN(jobNames...),
			tRollup.SourceKind.EQ(mysql.String(sourceKind)),
			tRollup.SourceKey.EQ(mysql.String(sourceKey)),
			tRollup.MetricKey.EQ(mysql.String(metricKey)),
		)).
		GROUP_BY(periodExpr).
		ORDER_BY(periodExpr.ASC())

	items := []*DailyValue{}
	if err := stmt.QueryContext(ctx, s.db, &items); err != nil {
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
	tRollup := table.FivenetStatsDailyRollup

	startDay = timeutils.StartOfDay(startDay)
	endDay = timeutils.StartOfDay(endDay)
	if endDay.Before(startDay) {
		return errors.New("end day before start day")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := tRollup.
		DELETE().
		WHERE(mysql.AND(
			tRollup.Day.GT_EQ(mysql.DateT(startDay)),
			tRollup.Day.LT_EQ(mysql.DateT(endDay)),
			tRollup.SourceKind.EQ(mysql.String(sourceKind)),
		))
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
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
		LIMIT(10000).
		ExecContext(ctx, tx)
	return err
}
