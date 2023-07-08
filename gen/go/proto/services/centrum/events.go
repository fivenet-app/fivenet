package centrum

import (
	"fmt"
	"strings"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/nats-io/nats.go"
)

const BaseSubject = "centrum"

const (
	TopicDispatch          events.Topic = "dispatch"
	TypeDispatchUpdated    events.Type  = "updated"
	TypeDispatchStatus     events.Type  = "status"
	TypeDispatchAssigned   events.Type  = "assigned"
	TypeDispatchUnassigned events.Type  = "unassigned"

	TopicUnit        events.Topic = "unit"
	TypeUnitAssigned events.Type  = "assigned"
	TypeUnitUpdated  events.Type  = "updated"
	TypeUnitStatus   events.Type  = "status"
	TypeUnitDeleted  events.Type  = "deleted"
)

func (s *Server) registerEvents() error {
	cfg := &nats.StreamConfig{
		Name:      "CENTRUM",
		Retention: nats.InterestPolicy,
		Subjects:  []string{BaseSubject + ">"},
	}

	s.events.JS.AddStream(cfg)

	return nil
}

func (s *Server) getEventTypeFromSubject(subject string) (events.Topic, events.Type) {
	split := strings.Split(subject, ".")
	if len(split) < 3 {
		return "", ""
	}

	return events.Topic(split[2]), events.Type(split[3])
}

func (s *Server) buildSubject(topic events.Topic, tType events.Type, userInfo *userinfo.UserInfo, id uint64) string {
	format := BaseSubject + ".%s." + string(topic) + "." + string(tType)
	if id > 0 {
		return fmt.Sprintf(format+".%d", userInfo.Job, id)
	}

	return fmt.Sprintf(format, userInfo.Job)
}
