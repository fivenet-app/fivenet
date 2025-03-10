package centrum

func (x *Settings) Default(job string) {
	x.Job = job

	if x.Mode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.Mode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.FallbackMode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.FallbackMode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.PredefinedStatus == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	}

	if x.Timings == nil {
		x.Timings = &Timings{}
	}
	if x.Timings.DispatchMaxWait == 0 {
		x.Timings.DispatchMaxWait = 900
	}
	if x.Timings.RequireUnitReminderSeconds == 0 {
		x.Timings.RequireUnitReminderSeconds = 180
	}
}
