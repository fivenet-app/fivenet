package dispatches

import "github.com/galexrt/arpanet/pkg/permissions"

func init() {
	permissions.RegisterPerms([]*permissions.Perm{
		{Key: "dispatches", Name: "View"},
	})
}

type Server struct {
	DispatchesServiceServer
}

func NewServer() *Server {
	return &Server{}
}
