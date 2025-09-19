package centrum

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

const defaultDeduplicationDuration = 3 * time.Minute

func (x *Settings) Default(job string) {
	x.Job = job

	if x.GetMode() <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.Mode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.GetFallbackMode() <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.FallbackMode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.GetPredefinedStatus() == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	}

	if x.GetTimings() == nil {
		x.Timings = &Timings{}
	}
	if x.GetTimings().GetDispatchMaxWait() <= 0 {
		x.Timings.DispatchMaxWait = 900
	}
	if x.GetTimings().GetRequireUnitReminderSeconds() <= 0 {
		x.Timings.RequireUnitReminderSeconds = 180
	}

	if x.GetConfiguration() == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(defaultDeduplicationDuration),
		}
	} else if x.GetConfiguration().GetDeduplicationDuration() == nil {
		x.Configuration.DeduplicationDuration = durationpb.New(defaultDeduplicationDuration)
	}
}

func (x *Settings) Merge(in *Settings) *Settings {
	x.Job = in.GetJob()
	x.Enabled = in.GetEnabled()

	x.Mode = in.GetMode()
	x.FallbackMode = in.GetFallbackMode()

	x.Type = in.GetType()
	x.Public = in.GetPublic()

	if in.GetPredefinedStatus() == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	} else {
		x.PredefinedStatus = in.GetPredefinedStatus()
	}

	if in.GetTimings() == nil {
		x.Timings = &Timings{}
	} else {
		x.Timings = in.GetTimings()
	}

	if in.GetAccess() == nil {
		x.Access = &CentrumAccess{}
	} else {
		x.Access = in.GetAccess()
	}

	if in.GetOfferedAccess() == nil {
		x.OfferedAccess = &CentrumAccess{}
	} else {
		x.OfferedAccess = in.GetOfferedAccess()
	}

	if in.GetEffectiveAccess() == nil {
		x.EffectiveAccess = &EffectiveAccess{}
	} else {
		x.EffectiveAccess = in.GetEffectiveAccess()
	}

	if in.GetConfiguration() == nil {
		x.Configuration = &Configuration{
			DeduplicationEnabled:  true,
			DeduplicationRadius:   45,
			DeduplicationDuration: durationpb.New(defaultDeduplicationDuration),
		}
	} else {
		x.Configuration = in.GetConfiguration()
	}

	return x
}

func (x *Settings) JobHasAccess(job string, access CentrumAccessLevel) bool {
	if x.GetJob() == job {
		return true // Same job is explicitly allowed
	}

	if x.GetEffectiveAccess() == nil || x.Access.Jobs == nil {
		return x.GetJob() == job // No access restrictions defined, only the job itself is allowed
	}

	if x.GetEffectiveAccess().GetDispatches() == nil {
		return false // No dispatch access defined
	}

	for _, j := range x.GetEffectiveAccess().GetDispatches().GetJobs() {
		if j.GetJob() == job && j.GetAccess() >= access {
			return true // Job is explicitly allowed in the access list
		}
	}

	return false // Fallback access is to deny
}
