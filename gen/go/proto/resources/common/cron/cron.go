package cron

func (x *Cronjob) Merge(in *Cronjob) *Cronjob {
	x.Schedule = in.Schedule

	return x
}

func (x *CronjobData) Merge(in *CronjobData) *CronjobData {
	x.Data = in.Data

	return x
}
