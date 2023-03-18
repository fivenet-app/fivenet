package dispatcher

type Server struct {
	DispatcherServiceServer
}

func NewServer() *Server {
	return &Server{}
}
