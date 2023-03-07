package job

import "github.com/galexrt/arpanet/pkg/perms"

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "job", Name: "View"},
	})
}

type Server struct {
	JobServiceServer
}

func NewServer() *Server {
	return &Server{}
}
