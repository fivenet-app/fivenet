package cron

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

	if in.Data != nil {
		x.Data = in.Data
	}

	return x
}

func (x *CronjobData) Merge(in *CronjobData) *CronjobData {
	x.Data = in.Data

	return x
}
