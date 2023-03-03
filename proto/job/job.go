package job

import "github.com/galexrt/arpanet/pkg/permissions"

func init() {
	permissions.RegisterPerms([]*permissions.Perm{
		{Key: "job", Name: "View"},
	})
}

type Server struct {
	JobServiceServer
}

func NewServer() *Server {
	return &Server{}
}
