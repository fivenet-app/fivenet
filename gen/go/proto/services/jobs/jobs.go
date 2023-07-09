package jobs

type Server struct {
	JobsServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// TODO time clock
// TODO employee warns/notes management
