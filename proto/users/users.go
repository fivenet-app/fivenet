package users

import (
	context "context"
	"strings"

	"github.com/galexrt/arpanet/api"
	"github.com/galexrt/arpanet/pkg/helpers"
)

type Server struct {
	UsersServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) FindUsers(ctx context.Context, req *FindUsersRequest) (*FindUsersResponse, error) {
	resp := &FindUsersResponse{}

	req.Firstname = strings.ReplaceAll(req.Firstname, "%", "")
	req.Lastname = strings.ReplaceAll(req.Lastname, "%", "")

	users, count, err := api.Users.SearchUsersByNamePages(req.Firstname, req.Lastname, req.Current, req.OrderBy...)
	if err != nil {
		return resp, err
	}

	resp.TotalCount = count
	resp.Current = req.Current
	resp.End = resp.Current + int64(len(users))

	for _, user := range users {
		resp.Users = append(resp.Users, helpers.ConvertModelUserToCommonCharacter(user))
	}

	return resp, nil
}
func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	resp := &GetUserResponse{}

	user, err := api.Users.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return resp, err
	}

	resp.User = helpers.ConvertModelUserToCommonCharacter(user)

	return resp, nil
}

func (s *Server) SetUserProps(ctx context.Context, req *SetUserPropsRequest) (*SetUserPropsResponse, error) {
	resp := &SetUserPropsResponse{}

	// TODO

	return resp, nil
}
