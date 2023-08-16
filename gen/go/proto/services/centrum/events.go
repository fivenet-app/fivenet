package centrum

import (
	"errors"
	"fmt"
	"strings"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject events.Subject = "centrum"

	TopicGeneral          events.Topic = "general"
	TypeGeneralSettings   events.Type  = "settings"
	TypeGeneralDisponents events.Type  = "disponents"

	TopicDispatch       events.Topic = "dispatch"
	TypeDispatchCreated events.Type  = "created"
	TypeDispatchDeleted events.Type  = "deleted"
	TypeDispatchUpdated events.Type  = "updated"
	TypeDispatchStatus  events.Type  = "status"

	TopicUnit            events.Topic = "unit"
	TypeUnitDeleted      events.Type  = "deleted"
	TypeUnitUpdated      events.Type  = "updated"
	TypeUnitStatus       events.Type  = "status"
	TypeUnitUserAssigned events.Type  = "user_assigned"
)

func (s *Server) registerEvents() error {
	cfg := &nats.StreamConfig{
		Name:      "CENTRUM",
		Retention: nats.InterestPolicy,
		Subjects:  []string{fmt.Sprintf("%s.>", BaseSubject)},
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
	_, topic, eType := s.splitSubject(subject)
	return topic, eType
}

func (s *Server) splitSubject(subject string) (string, events.Topic, events.Type) {
	split := strings.Split(subject, ".")
	if len(split) < 3 {
		return "", "", ""
	}

	return split[1], events.Topic(split[2]), events.Type(split[3])
}

func (s *Server) buildSubject(topic events.Topic, tType events.Type, job string, id uint64) string {
	format := "%s.%s." + string(topic) + "." + string(tType)
	if id > 0 {
		return fmt.Sprintf(format+".%d", BaseSubject, job, id)
	}

	return fmt.Sprintf(format, BaseSubject, job)
}

func (s *Server) broadcastToAllUnits(topic events.Topic, tType events.Type, job string, data []byte) {
	units, ok := s.units.Load(job)
	if !ok {
		return
	}

	units.Range(func(key uint64, unit *dispatch.Unit) bool {
		s.events.JS.Publish(s.buildSubject(TopicGeneral, TypeGeneralSettings, job, unit.Id), data)
		return true
	})
}
