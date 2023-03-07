package users

import (
	context "context"
	"strings"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "users", Name: "View"},
		{Key: "users", Name: "FindUsers", Fields: []string{"Licenses", "UserProps"}},
		{Key: "users", Name: "SetUserProps", Fields: []string{"Wanted"}},
		{Key: "users", Name: "GetUserActivityRequest", Fields: []string{"CauseUser", ""}},
	})
}

var (
	u   = table.Users
	aup = table.ArpanetUserProps
	ul  = table.UserLicenses
)

type Server struct {
	UsersServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) FindUsers(ctx context.Context, req *FindUsersRequest) (*FindUsersResponse, error) {
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Permission check
	if !perms.P.Can(user, "users", "FindUsers") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to find users")
	}

	selectors := make(jet.ProjectionList, len(common.CharacterBaseColumns))
	copy(selectors, common.CharacterBaseColumns)
	if perms.P.Can(user, "users", "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}

	req.Firstname = strings.ReplaceAll(req.Firstname, "%", "")
	req.Lastname = strings.ReplaceAll(req.Lastname, "%", "")

	condition := jet.Bool(true)
	if req.Firstname != "" {
		condition = condition.AND(u.Firstname.LIKE(jet.String("%" + req.Firstname + "%")))
	}
	if req.Lastname != "" {
		condition = condition.AND(u.Lastname.LIKE(jet.String("%" + req.Lastname + "%")))
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
	}

	if len(orderBys) == 0 {
		orderBys = append(orderBys, u.Firstname.ASC())
	}

	stmt := u.SELECT(
		selectors[0], selectors[1:]...,
	).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(
			u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID)),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(req.Offset).
		LIMIT(common.DefaultPageLimit)

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
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Permission check
	if !perms.P.Can(user, "users", "FindUsers") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to find users")
	}

	selectors := make(jet.ProjectionList, len(common.CharacterBaseColumns))
	copy(selectors, common.CharacterBaseColumns)
	if perms.P.Can(user, "users", "FindUsers", "UserProps") {
		selectors = append(selectors, aup.Wanted)
	}
	if perms.P.Can(user, "users", "FindUsers", "Licenses") {
		selectors = append(selectors, ul.Type)
	}

	resp := &GetUserResponse{}

	stmt := u.SELECT(
		selectors[0], selectors[1:]...,
	).
		FROM(
			u.LEFT_JOIN(aup, aup.UserID.EQ(u.ID)).
				LEFT_JOIN(ul, ul.Owner.EQ(u.Identifier)),
		).
		WHERE(u.ID.EQ(jet.Int32(req.UserID))).
		LIMIT(15)

	var dest common.User
	if err := stmt.QueryContext(ctx, query.DB, &dest); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Permission check
	if !perms.P.Can(user, "users.SetUserProps") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set user properties!")
	}
	if !perms.P.Can(user, "users", "SetUserProps", "Wanted") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set an user wanted!")
	}

	stmt := aup.INSERT(aup.AllColumns).
		MODEL(req).
		ON_DUPLICATE_KEY_UPDATE(
			aup.Wanted.SET(jet.Bool(*req.Wanted)),
		)
	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	return &SetUserPropsResponse{}, nil
}

func (s *Server) GetUserActivity(ctx context.Context, req *GetUserActivityRequest) (*GetUserActivityResponse, error) {
	var activities []*UserActivity

	ua := table.ArpanetUserActivity
	if err := ua.SELECT(ua.AllColumns).
		FROM(ua).
		WHERE(
			ua.TargetUserID.EQ(jet.Int32(req.UserID)),
		).
		LIMIT(10).
		QueryContext(ctx, query.DB, &activities); err != nil {
		return nil, err
	}

	resp := &GetUserActivityResponse{
		Activity: activities,
	}
	return resp, nil
}
