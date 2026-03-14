package jobs

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/stats"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	pkgstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetStats(
	ctx context.Context,
	req *pbstats.GetStatsRequest,
) (*pbstats.GetStatsResponse, error) {
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

	if end.Sub(start) > 365*24*time.Hour {
		return nil, status.Error(codes.InvalidArgument, "range must not exceed 365 days")
	}

	period := max(req.GetPeriod(), stats.StatsPeriod_STATS_PERIOD_DAILY)

	resp := &pbstats.GetStatsResponse{
		PeriodValues:       []*stats.DailyValue{},
		PeriodSeriesValues: []*stats.PeriodSeriesValue{},
	}

	switch req.GetCategory() {
	case stats.StatsCategory_STATS_CATEGORY_EMPLOYEE_COUNT_OVER_TIME:
		seriesValues, err := s.stats.QueryEmployeeSeriesOverTime(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			period,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		resp.PeriodSeriesValues = make([]*stats.PeriodSeriesValue, 0, len(seriesValues))
		resp.PeriodValues = make([]*stats.DailyValue, 0, len(seriesValues))
		for _, item := range seriesValues {
			resp.PeriodSeriesValues = append(
				resp.PeriodSeriesValues,
				&stats.PeriodSeriesValue{
					Day:   timestamp.New(item.Day),
					Key:   item.Key,
					Label: item.Label,
					Value: item.Value,
				},
			)

			if item.Key == "employee_count" {
				resp.PeriodValues = append(
					resp.PeriodValues,
					&stats.DailyValue{
						Day:   timestamp.New(item.Day),
						Value: item.Value,
					},
				)
			}
		}

		// Average employee count over the period
		averageValue, err := s.stats.QueryAverageValue(
			ctx,
			start,
			end,
			userInfo.GetJob(),
			pkgstats.SourceKindEmployeeCount,
			"fivenet_user_jobs",
			"employee_count",
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		resp.AverageValue = averageValue

	case stats.StatsCategory_STATS_CATEGORY_UNSPECIFIED:
		fallthrough
	default:
		return nil, status.Error(codes.InvalidArgument, "stats category is required")
	}

	return resp, nil
}
