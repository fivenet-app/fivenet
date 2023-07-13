package centrum

import (
	"errors"
	"fmt"
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/nats-io/nats.go"
)

const BaseSubject = "centrum"

const (
	TopicGeneral           events.Topic = "general"
	TypeGeneralSettings    events.Type  = "settings"
	TypeGeneralControllers events.Type  = "controllers"

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
		Subjects:  []string{BaseSubject + ".>"},
	}

	if _, err := s.events.JS.AddStream(cfg); err != nil {
		if !errors.Is(nats.ErrStreamNameAlreadyInUse, err) {
			return err
		}

		if _, err := s.events.JS.UpdateStream(cfg); err != nil {
			return err
		}
	}

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

func (s *Server) broadcastToAllUnits(topic events.Topic, tType events.Type, userInfo *userinfo.UserInfo, data []byte) {
	jobUnits, ok := s.units.Load(userInfo.Job)
	if !ok {
		return
	}

	jobUnits.Range(func(key uint64, unit *dispatch.Unit) bool {
		s.events.JS.Publish(s.buildSubject(TopicGeneral, TypeGeneralSettings, userInfo, unit.Id), data)
		return true
	})
}
