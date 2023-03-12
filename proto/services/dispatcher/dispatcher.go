package dispatcher

import "github.com/galexrt/arpanet/pkg/perms"

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "dispatches", Name: "View"},
	})
}

type Server struct {
	DispatcherServiceServer
}

func NewServer() *Server {
	return &Server{}
}
