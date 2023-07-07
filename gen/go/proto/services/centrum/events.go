package centrum

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject = "centrum"
)

func (s *Server) registerEvents() error {
	var err error
	s.eventSub, err = s.events.NC.Subscribe(BaseSubject+".>", s.handleMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleMessage(msg *nats.Msg) {
	spew.Dump(msg)
}
