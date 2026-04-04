package syncers

import "time"

const (
	maxAdaptiveDrainBatchesPerSync = 30
	lagWarnThreshold               = 5 * time.Minute
	lagPerAdditionalBatch          = 30 * time.Second
)

func calculateDrainBatchCap(
	lastCheck *time.Time,
	windowEnd time.Time,
	baseCap int,
) (int, time.Duration) {
	if lastCheck == nil || lastCheck.IsZero() {
		return baseCap, 0
	}

	lag := windowEnd.Sub(*lastCheck)
	if lag <= 0 {
		return baseCap, 0
	}

	extra := int(lag / lagPerAdditionalBatch)
	batchCap := min(baseCap+extra, maxAdaptiveDrainBatchesPerSync)

	return batchCap, lag
}
