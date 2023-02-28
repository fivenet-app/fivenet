package documents

type Server struct {
	DocumentsServiceServer
}

func NewServer() *Server {
	return &Server{}
}
