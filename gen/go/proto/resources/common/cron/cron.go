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
