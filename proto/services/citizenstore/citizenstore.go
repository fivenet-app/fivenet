package citizenstore

import (
	context "context"
	"strconv"
	"strings"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/common/database"
	users "github.com/galexrt/arpanet/proto/resources/users"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	u   = table.Users.AS("user")
	aup = table.ArpanetUserProps
	ul  = table.UserLicenses
	aua = table.ArpanetUserActivity
)

type Server struct {
	CitizenStoreServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) FindUsers(ctx context.Context, req *FindUsersRequest) (*FindUsersResponse, error) {
	userID, _, _ := auth.GetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
		u.Firstname,
		u.Lastname,
		u.Dateofbirth,
		u.Sex,
		u.Height,
		u.PhoneNumber,
		u.Visum,
		u.Playtime,
	}
	// Field Permission Check
	if perms.P.CanID(userID, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}

	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")

	condition := jet.Bool(true)
	if req.SearchName != "" {
		condition = condition.AND(jet.BoolExp(jet.Raw("MATCH(firstname,lastname) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.SearchName})))
	}
	if req.Wanted {
		condition = condition.AND(aup.Wanted.IS_TRUE())
	}

	// Get total count of values
	countStmt := u.SELECT(
		jet.COUNT(u.ID).AS("total_count"),
	).
		FROM(
			u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID)),
		).
		WHERE(condition)

	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, query.DB, &count); err != nil {
		return nil, err
	}

	resp := &FindUsersResponse{
		Offset:     req.Offset,
		TotalCount: count.TotalCount,
		End:        0,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := u.SELECT(
		selectors[0], selectors[1:]...,
	).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(
			u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID)),
		).
		WHERE(condition).
		OFFSET(req.Offset).
		LIMIT(database.PaginationLimit)

	// Convert our proto abstracted `common.OrderBy` to actual gorm order by instructions
	orderBys := []jet.OrderByClause{}
	if len(req.OrderBy) > 0 {
		for _, orderBy := range req.OrderBy {
			var column jet.Column
			switch orderBy.Column {
			case "firstname":
				column = u.Firstname
			case "job":
				column = u.Job
			}

			if orderBy.Desc {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}

		stmt = stmt.ORDER_BY(orderBys...)
	}

	if err := stmt.QueryContext(ctx, query.DB, &resp.Users); err != nil {
		return nil, err
	}

	resp.TotalCount = count.TotalCount
	if req.Offset >= resp.TotalCount {
		resp.Offset = 0
	} else {
		resp.Offset = req.Offset
	}
	resp.End = resp.Offset + int64(len(resp.Users))

	return resp, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	userID, _, _ := auth.GetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
		u.Firstname,
		u.Lastname,
		u.Dateofbirth,
		u.Sex,
		u.Height,
		u.PhoneNumber,
		u.Visum,
		u.Playtime,
	}

	// Field Permission Check
	if perms.P.CanID(userID, CitizenStoreServicePermKey, "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}
	if perms.P.CanID(userID, CitizenStoreServicePermKey, "FindUsers", "Licenses") {
		selectors = append(selectors, ul.Type)
	}

	resp := &GetUserResponse{
		User: &users.User{},
	}
	stmt := u.SELECT(
		selectors[0], selectors[1:]...,
	).
		FROM(
			u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID)).
				LEFT_JOIN(ul, ul.Owner.EQ(u.Identifier)),
		).
		WHERE(u.ID.EQ(jet.Int32(req.UserId))).
		LIMIT(15)

	if err := stmt.QueryContext(ctx, query.DB, resp.User); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetUserActivity(ctx context.Context, req *GetUserActivityRequest) (*GetUserActivityResponse, error) {
	userID, _, _ := auth.GetUserInfoFromContext(ctx)

	resp := &GetUserActivityResponse{}
	// An user can never see their own activity on their own "profile"
	if userID == req.UserId {
		return resp, nil
	}

	uTarget := u.AS("target_user")
	uSource := u.AS("source_user")
	stmt := aua.SELECT(
		aua.AllColumns,
		uTarget.ID,
		uTarget.Identifier,
		uTarget.Job,
		uTarget.JobGrade,
		uTarget.Firstname,
		uTarget.Lastname,
		uSource.ID,
		uSource.Identifier,
		uSource.Job,
		uSource.JobGrade,
		uSource.Firstname,
		uSource.Lastname,
	).
		FROM(
			aua.LEFT_JOIN(uTarget,
				uTarget.ID.EQ(aua.TargetUserID),
			).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(aua.SourceUserID),
				),
		).
		WHERE(
			aua.TargetUserID.EQ(jet.Int32(req.UserId)),
		).
		LIMIT(12)

	if err := stmt.QueryContext(ctx, query.DB, &resp.Activity); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	userID := auth.GetUserIDFromContext(ctx)

	// Field Permission Check
	if !perms.P.CanID(userID, CitizenStoreServicePermKey, "SetUserProps", "Wanted") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set user wanted status!")
	}

	stmt := aup.INSERT(aup.AllColumns).
		MODEL(req).
		ON_DUPLICATE_KEY_UPDATE(
			aup.Wanted.SET(jet.Bool(req.Props.Wanted)),
		)
	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	// Create user activity
	key := "UserProps.Wanted"
	newValue := strconv.FormatBool(req.Props.Wanted)
	oldValue := strconv.FormatBool(!req.Props.Wanted)
	s.addUserAcitvity(ctx, &model.ArpanetUserActivity{
		SourceUserID: userID,
		TargetUserID: req.UserId,
		Type:         int16(users.USER_ACTIVITY_TYPE_CHANGED),
		Key:          key,
		OldValue:     &oldValue,
		NewValue:     &newValue,
	})

	return &SetUserPropsResponse{}, nil
}

func (s *Server) addUserAcitvity(ctx context.Context, activity *model.ArpanetUserActivity) error {
	stmt := aua.INSERT(
		aua.SourceUserID,
		aua.TargetUserID,
		aua.Type,
		aua.Key,
		aua.OldValue,
		aua.NewValue,
	).MODEL(activity)

	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return err
	}

	return nil
}
