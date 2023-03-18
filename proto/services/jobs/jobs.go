package jobs

type Server struct {
	JobsServiceServer
}

func NewServer() *Server {
	return &Server{}
}
