package cron

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const DefaultCronTimeout = 10 * time.Second

func (x *Cronjob) Merge(in *Cronjob) *Cronjob {
	x.Schedule = in.Schedule

	if in.State > CronjobState_CRONJOB_STATE_UNSPECIFIED {
		x.State = in.State
	}

	if in.NextScheduleTime != nil {
		x.NextScheduleTime = in.NextScheduleTime
	}

	if in.LastAttemptTime != nil {
		x.LastAttemptTime = in.LastAttemptTime
	}

	if in.StartedTime != nil {
		x.StartedTime = in.StartedTime
	}

	x.Timeout = in.Timeout

	if in.Data != nil {
		x.Data.Merge(in.Data)
	}

	if in.LastCompletedEvent != nil {
		x.LastCompletedEvent = in.LastCompletedEvent
	}

	return x
}

func (x *Cronjob) GetRunTimeout() time.Duration {
	if x.Timeout == nil {
		return DefaultCronTimeout
	}

	return x.Timeout.AsDuration()
}

func (x *CronjobData) Merge(in *CronjobData) *CronjobData {
	if x == nil {
		x = in
	} else {
		x.Data = in.Data
	}

	return x
}

func (x *GenericCronData) HasAttribute(key string) bool {
	if x.Attributes == nil {
		return false
	}

	_, ok := x.Attributes[key]
	return ok
}

func (x *GenericCronData) GetAttribute(key string) string {
	if x.Attributes == nil {
		return ""
	}

	return x.Attributes[key]
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

	delete(x.Attributes, key)
}

func (x *CronjobData) Unmarshal(dest proto.Message) error {
	if x == nil || dest == nil {
		return fmt.Errorf("invalid input: CronjobData or destination is nil")
	}

	expectedTypeURL := "type.googleapis.com/" + string(proto.MessageName(dest))

	if x.Data != nil && x.Data.TypeUrl == expectedTypeURL {
		// Valid type - attempt to unmarshal
		if err := x.Data.UnmarshalTo(dest); err != nil {
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
		return fmt.Errorf("invalid input: CronjobData or source is nil")
	}

	anyMsg, err := anypb.New(src)
	if err != nil {
		return fmt.Errorf("failed to marshal cron data into Any. %w", err)
	}

	x.Data = anyMsg
	return nil
}
