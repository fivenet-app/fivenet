package jobs

type Server struct {
	JobsServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// TODO employee warns/notes management
// TODO time clock
