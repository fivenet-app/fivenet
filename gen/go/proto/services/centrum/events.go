package centrum

import (
	"fmt"

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
	fmt.Printf("CENTRUM MSG: %+v\n", msg)
}
