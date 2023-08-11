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
	TypeUnitCreated      events.Type  = "created"
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
	split := strings.Split(subject, ".")
	if len(split) < 3 {
		return "", ""
	}

	return events.Topic(split[2]), events.Type(split[3])
}

func (s *Server) buildSubject(topic events.Topic, tType events.Type, userInfo *userinfo.UserInfo, id uint64) string {
	format := "%s.%s." + string(topic) + "." + string(tType)
	if id > 0 {
		return fmt.Sprintf(format+".%d", BaseSubject, userInfo.Job, id)
	}

	return fmt.Sprintf(format, BaseSubject, userInfo.Job)
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
