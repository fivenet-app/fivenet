package users

import (
	context "context"
	"errors"
	"strings"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/helpers"
	"github.com/galexrt/arpanet/pkg/permissions"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

func init() {
	permissions.RegisterPerms([]*permissions.Perm{
		{Key: "users", Name: "View"},
		{Key: "users", Name: "FindUsers", Fields: []string{"Licenses", "UserProps"}},
		{Key: "users", Name: "SetUserProps", Fields: []string{"Wanted"}},
		{Key: "users", Name: "GetUserActivityRequest", Fields: []string{"CauseUser", ""}},
	})
}

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

	if !auth.CanUser(user, "users.FindUsers") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to find users")
	}

	req.Firstname = strings.ReplaceAll(req.Firstname, "%", "")
	req.Lastname = strings.ReplaceAll(req.Lastname, "%", "")

	u := query.User
	q := u.Clauses(hints.UseIndex("idx_users_firstname_lastname"))

	if auth.CanUserAccessField(user, "users.FindUsers", "Licenses") {
		q = q.Preload(u.UserLicenses.RelationField)
	}
	if auth.CanUserAccessField(user, "users.FindUsers", "UserProps") {
		q = q.Preload(u.UserProps.RelationField)
	}

	if req.Firstname != "" {
		q = q.Where(u.Firstname.Like("%" + req.Firstname + "%"))
	}
	if req.Lastname != "" {
		q = q.Where(u.Lastname.Like("%" + req.Lastname + "%"))
	}

	// Convert our proto abstracted `common.OrderBy` to actual gorm order by instructions
	if len(req.OrderBy) > 0 {
		for _, orderBy := range req.OrderBy {
			field, ok := u.GetFieldByName(orderBy.Column)
			if !ok {
				return nil, errors.New("orderBy column not found")
			}

			if orderBy.Desc {
				q = q.Order(field.Desc())
			} else {
				q = q.Order(field)
			}
		}
	} else {
		q = q.Order(u.Firstname)
	}

	users, count, err := q.FindByPage(int(req.Current), common.DefaultPageLimit)
	if err != nil {
		return nil, err
	}

	resp := &FindUsersResponse{}
	resp.TotalCount = count
	resp.Current = req.Current
	resp.End = resp.Current + int64(len(users))

	for _, user := range users {
		resp.Users = append(resp.Users, helpers.ConvertModelUserToCommonCharacter(user))
	}

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !auth.CanUser(user, "users.SetUserProps") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set user properties!")
	}
	if !auth.CanUserAccessField(user, "users.SetUserProps", "Wanted") {
		return nil, status.Error(codes.PermissionDenied, "You are not allowed to set an user wanted!")
	}

	userProps := &model.UserProps{
		UserID: int32(req.UserID),
		Wanted: req.Wanted,
	}

	ups := query.UserProps
	if err := ups.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(userProps); err != nil {
		return nil, err
	}

	resp := &SetUserPropsResponse{}
	return resp, nil
}

func (s *Server) GetUserActivity(ctx context.Context, req *GetUserActivityRequest) (*GetUserActivityResponse, error) {
	ua := query.UserActivity
	activities, err := ua.Where(ua.TargetUserID.Eq(int32(req.UserID))).Limit(10).Find()
	if err != nil {
		return nil, err
	}

	protoActivity := make([]*UserActivity, len(activities))
	for k := 0; k < len(activities); k++ {
		protoActivity[k] = &UserActivity{
			UserID:    uint64(activities[k].ID),
			CreatedAt: timestamppb.New(activities[k].CreatedAt),
			TargetUser: &common.ShortCharacter{
				UserID: 0,
			},
			CauseUser: &common.ShortCharacter{},
			Key:       activities[k].Key,
			OldValue:  activities[k].OldValue,
			NewValue:  activities[k].NewValue,
			Reason:    activities[k].Reason,
		}
	}

	resp := &GetUserActivityResponse{
		Activity: protoActivity,
	}
	return resp, nil
}
