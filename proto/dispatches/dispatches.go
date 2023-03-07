package dispatches

import "github.com/galexrt/arpanet/pkg/perms"

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "dispatches", Name: "View"},
	})
}

type Server struct {
	DispatchesServiceServer
}

func NewServer() *Server {
	return &Server{}
}
