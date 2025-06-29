package cron

import (
	"time"
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
