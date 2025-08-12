package cron

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const DefaultCronTimeout = 10 * time.Second

func (x *Cronjob) Merge(in *Cronjob) *Cronjob {
	x.Schedule = in.GetSchedule()

	if in.GetState() > CronjobState_CRONJOB_STATE_UNSPECIFIED {
		x.State = in.GetState()
	}

	if in.GetNextScheduleTime() != nil {
		x.NextScheduleTime = in.GetNextScheduleTime()
	}

	if in.GetLastAttemptTime() != nil {
		x.LastAttemptTime = in.GetLastAttemptTime()
	}

	if in.GetStartedTime() != nil {
		x.StartedTime = in.GetStartedTime()
	}

	x.Timeout = in.GetTimeout()

	if in.GetData() != nil {
		x.GetData().Merge(in.GetData())
	}

	if in.GetLastCompletedEvent() != nil {
		x.LastCompletedEvent = in.GetLastCompletedEvent()
	}

	return x
}

func (x *Cronjob) GetRunTimeout() time.Duration {
	if x.GetTimeout() == nil {
		return DefaultCronTimeout
	}

	return x.GetTimeout().AsDuration()
}

func (x *CronjobData) Merge(in *CronjobData) *CronjobData {
	if x == nil {
		x = in
	} else {
		x.Data = in.GetData()
	}

	return x
}

func (x *GenericCronData) HasAttribute(key string) bool {
	if x.Attributes == nil {
		return false
	}

	_, ok := x.GetAttributes()[key]
	return ok
}

func (x *GenericCronData) GetAttribute(key string) string {
	if x.Attributes == nil {
		return ""
	}

	return x.GetAttributes()[key]
}

func (x *GenericCronData) SetAttribute(key string, value string) {
	if x.Attributes == nil {
		x.Attributes = make(map[string]string)
	}

	x.Attributes[key] = value
}

func (x *GenericCronData) DeleteAttribute(key string) {
	if x.Attributes == nil {
		return
	}

	delete(x.GetAttributes(), key)
}

func (x *CronjobData) Unmarshal(dest proto.Message) error {
	if x == nil || dest == nil {
		return errors.New("invalid input: CronjobData or destination is nil")
	}

	expectedTypeURL := "type.googleapis.com/" + string(proto.MessageName(dest))

	if x.GetData() != nil && x.GetData().GetTypeUrl() == expectedTypeURL {
		// Valid type - attempt to unmarshal
		if err := x.GetData().UnmarshalTo(dest); err != nil {
			return fmt.Errorf("failed to unmarshal cron data. %w", err)
		}
	} else {
		// Reset to empty message of expected type
		anyMsg, err := anypb.New(dest)
		if err != nil {
			return fmt.Errorf("failed to create new Any for cron data. %w", err)
		}
		x.Data = anyMsg
	}

	return nil
}

func (x *CronjobData) MarshalFrom(src proto.Message) error {
	if x == nil || src == nil {
		return errors.New("invalid input: CronjobData or source is nil")
	}

	anyMsg, err := anypb.New(src)
	if err != nil {
		return fmt.Errorf("failed to marshal cron data into Any. %w", err)
	}

	x.Data = anyMsg
	return nil
}
