package jobs

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) addGroupActivity(
	ctx context.Context,
	db qrm.DB,
	job string,
	groupID int64,
	activityType jobsgroups.GroupActivityType,
	actorUserID int32,
	targetUserID int32,
	ruleID int64,
	reason *string,
	data *jobsgroups.GroupActivityData,
) error {
	return s.store.CreateGroupActivity(ctx, db, &jobsgroups.GroupActivity{
		Job:          job,
		GroupId:      groupID,
		Type:         activityType,
		ActorUserId:  &actorUserID,
		TargetUserId: int32PtrOrNil(targetUserID),
		RuleId:       int64PtrOrNil(ruleID),
		Reason:       reason,
		Data:         data,
	})
}

func int32PtrOrNil(value int32) *int32 {
	if value == 0 {
		return nil
	}
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}

func int64PtrOrNil(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func (s *Server) ListGroupActivity(
	ctx context.Context,
	req *pbjobs.ListGroupActivityRequest,
) (*pbjobs.ListGroupActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", req.GetGroupId(),
	})

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: true,
	}, req.GetGroupId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	tActivity := table.FivenetJobGroupActivity.AS("group_activity")
	condition := mysql.AND(
		tActivity.Job.EQ(mysql.String(userInfo.GetJob())),
		tActivity.GroupID.EQ(mysql.Int64(req.GetGroupId())),
	)

	if len(req.GetTypes()) > 0 {
		types := make([]mysql.Expression, 0, len(req.GetTypes()))
		for _, activityType := range req.GetTypes() {
			if activityType == jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_UNSPECIFIED {
				continue
			}
			types = append(types, mysql.Int32(int32(activityType)))
		}
		if len(types) == 0 {
			resp := &pbjobs.ListGroupActivityResponse{
				Pagination: &database.PaginationResponse{PageSize: defaultPageSize},
				Activity:   []*jobsgroups.GroupActivity{},
			}
			grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
			return resp, nil
		}
		condition = condition.AND(tActivity.ActivityType.IN(types...))
	}

	if req.HasUserId() {
		userID := req.GetUserId()
		condition = condition.AND(mysql.OR(
			tActivity.ActorUserID.EQ(mysql.Int32(userID)),
			tActivity.TargetUserID.EQ(mysql.Int32(userID)),
		))
	}
	if req.HasFrom() {
		condition = condition.AND(
			tActivity.CreatedAt.GT_EQ(dbutils.TimestampToMySQL(req.GetFrom())),
		)
	}
	if req.HasTo() {
		condition = condition.AND(tActivity.CreatedAt.LT_EQ(dbutils.TimestampToMySQL(req.GetTo())))
	}

	query := jobsstore.ListQuery{
		Job:   userInfo.GetJob(),
		Where: condition,
		Sort:  req.GetSort(),
	}
	count, err := s.store.CountGroupActivity(ctx, s.db, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, defaultPageSize)
	resp := &pbjobs.ListGroupActivityResponse{
		Pagination: pag,
		Activity:   []*jobsgroups.GroupActivity{},
	}
	if count <= 0 {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
		return resp, nil
	}

	query.Offset = pag.GetOffset()
	query.Limit = limit
	resp.Activity, err = s.store.ListGroupActivity(ctx, s.db, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	resp.Pagination.Update(len(resp.GetActivity()))

	targets := appendGroupActivityColleagueTargets(nil, resp.GetActivity())
	if err := s.hydrateGroupColleagueTargets(ctx, userInfo, targets); err != nil {
		return nil, err
	}
	s.enrichGroupActivityRules(userInfo.GetJob(), resp.GetActivity())

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}

func (s *Server) enrichGroupActivityRules(job string, activities []*jobsgroups.GroupActivity) {
	for _, activity := range activities {
		if activity.GetData().GetRule() == nil {
			continue
		}

		s.enrichGroupRuleGradeLabels(job, activity.GetData().GetRule())
	}
}

func groupActivityRuleData(rule *jobsgroups.GroupRule) *jobsgroups.GroupActivityData {
	if rule == nil {
		return nil
	}

	return &jobsgroups.GroupActivityData{
		Data: &jobsgroups.GroupActivityData_Rule{
			Rule: rule,
		},
	}
}
