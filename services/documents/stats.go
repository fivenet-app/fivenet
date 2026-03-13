package documents

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetStats(
	ctx context.Context,
	req *pbdocuments.GetStatsRequest,
) (*pbdocuments.GetStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	start := time.Now().UTC().AddDate(0, 0, -14)
	end := time.Now().UTC()
	if req.GetStart() != nil {
		start = req.GetStart().AsTime()
	}
	if req.GetEnd() != nil {
		end = req.GetEnd().AsTime()
	}

	if end.Before(start) {
		return nil, status.Error(codes.InvalidArgument, "end must not be before start")
	}
	period := max(req.GetPeriod(), pbdocuments.StatsPeriod_STATS_PERIOD_DAILY)

	categories, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.StatsServicePerm,
		permsdocuments.StatsServiceGetStatsPerm,
		permsdocuments.StatsServiceGetStatsCategoriesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	_ = categories

	resp := &pbdocuments.GetStatsResponse{
		TopLaws:             []*pbdocuments.KeyValue{},
		FinesOverTime:       []*pbdocuments.DailyValue{},
		DocumentsByCategory: []*pbdocuments.CategoryValue{},
		PeriodValues:        []*pbdocuments.DailyValue{},
		PeriodSeriesValues:  []*pbdocuments.PeriodSeriesValue{},
	}

	switch req.GetCategory() {
	case pbdocuments.StatsCategory_STATS_CATEGORY_DOCUMENTS_BY_CATEGORY:
		byCategory, err := s.stats.QueryDocumentsByCategory(ctx, start, end, userInfo.GetJob())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.DocumentsByCategory = make([]*pbdocuments.CategoryValue, 0, len(byCategory))
		for _, item := range byCategory {
			resp.DocumentsByCategory = append(resp.DocumentsByCategory, &pbdocuments.CategoryValue{
				Id:    item.ID,
				Name:  item.Name,
				Color: item.Color,
				Value: item.Value,
			})
		}
		periodValues, err := s.stats.QueryPeriodValues(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			docstats.SourceKindDocumentColumn,
			"documents",
			"document_count",
			period,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.PeriodValues = make([]*pbdocuments.DailyValue, 0, len(periodValues))
		for _, item := range periodValues {
			resp.PeriodValues = append(resp.PeriodValues, &pbdocuments.DailyValue{
				Day:   timestamp.New(item.Day),
				Value: item.Value,
			})
		}
		totalValue, err := s.stats.QueryTotalValue(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			docstats.SourceKindDocumentColumn,
			"documents",
			"document_count",
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.TotalValue = totalValue

	case pbdocuments.StatsCategory_STATS_CATEGORY_TOP_LAWS:
		if !categories.Contains("PenaltyCalculator") {
			return nil, status.Error(
				codes.PermissionDenied,
				"user does not have permission to view penalty calculator stats",
			)
		}

		topLaws, err := s.stats.QueryTopLaws(ctx, start, end, userInfo.GetJob(), 10)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.TopLaws = make([]*pbdocuments.KeyValue, 0, len(topLaws))
		for _, item := range topLaws {
			resp.TopLaws = append(resp.TopLaws, &pbdocuments.KeyValue{
				Key:   item.Key,
				Value: item.Value,
			})
		}
		periodValues, err := s.stats.QueryPeriodValues(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			docstats.SourceKindDocumentMetric,
			docstats.PenaltyCalculatorSourceKey,
			"law_count",
			period,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.PeriodValues = make([]*pbdocuments.DailyValue, 0, len(periodValues))
		for _, item := range periodValues {
			resp.PeriodValues = append(resp.PeriodValues, &pbdocuments.DailyValue{
				Day:   timestamp.New(item.Day),
				Value: item.Value,
			})
		}
		totalValue, err := s.stats.QueryTotalValue(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			docstats.SourceKindDocumentMetric,
			docstats.PenaltyCalculatorSourceKey,
			"law_count",
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.TotalValue = totalValue

	case pbdocuments.StatsCategory_STATS_CATEGORY_PENALTIES_OVER_TIME:
		if !categories.Contains("PenaltyCalculator") {
			return nil, status.Error(
				codes.PermissionDenied,
				"user does not have permission to view penalty calculator stats",
			)
		}

		seriesValues, err := s.stats.QueryPenaltySeriesOverTime(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			period,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.PeriodSeriesValues = make([]*pbdocuments.PeriodSeriesValue, 0, len(seriesValues))
		resp.FinesOverTime = []*pbdocuments.DailyValue{}
		for _, item := range seriesValues {
			resp.PeriodSeriesValues = append(
				resp.PeriodSeriesValues,
				&pbdocuments.PeriodSeriesValue{
					Day:   timestamp.New(item.Day),
					Key:   item.Key,
					Label: item.Label,
					Value: item.Value,
				},
			)
			if item.Key == "fine_total" {
				resp.FinesOverTime = append(resp.FinesOverTime, &pbdocuments.DailyValue{
					Day:   timestamp.New(item.Day),
					Value: item.Value,
				})
			}
		}
		totalValue, err := s.stats.QueryTotalValue(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			docstats.SourceKindDocumentMetric,
			docstats.PenaltyCalculatorSourceKey,
			"fine_total",
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.TotalValue = totalValue

	// case pbdocuments.StatsCategory_STATS_CATEGORY_PENALTY_REDUCTION_SUMS:
	// 	reductionPercentSum, caseCountSum, err := s.stats.QueryPenaltyReductionSums(
	// 		ctx,
	// 		start,
	// 		end,
	// 		userInfo.GetJob(),
	// 	)
	// 	if err != nil {
	// 		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	// 	}
	// 	resp.ReductionPercentSum = reductionPercentSum
	// 	resp.CaseCountSum = caseCountSum

	case pbdocuments.StatsCategory_STATS_CATEGORY_UNSPECIFIED:
		fallthrough
	default:
		return nil, status.Error(codes.InvalidArgument, "stats category is required")
	}

	return resp, nil
}
