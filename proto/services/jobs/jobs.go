package jobs

import "github.com/galexrt/arpanet/pkg/perms"

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "jobs", Name: "View"},
	})
}

type Server struct {
	JobsServiceServer
}

func NewServer() *Server {
	return &Server{}
}
