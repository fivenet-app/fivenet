package centrum

type Server struct {
	CentrumServiceServer
	SquadServiceServer
}

func NewServer() *Server {
	return &Server{}
}
