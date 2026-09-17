package housekeeper

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	cronjobDispatchAssignmentExpiration = "centrum.housekeeper.dispatch_assignment_expiration"
	cronjobCleanupUnits                 = "centrum.housekeeper.cleanup_units"
	cronjobAuditUnitMembership          = "centrum.housekeeper.audit_unit_membership"
	cronjobAuditEmptyUnitDispatches     = "centrum.housekeeper.audit_empty_unit_dispatches"
	cronjobCancelOldDispatches          = "centrum.housekeeper.cancel_old_dispatches"
	cronjobDeleteOldDispatches          = "centrum.housekeeper.delete_old_dispatches"
	cronjobDeleteOldDispatchesFromKV    = "centrum.housekeeper.delete_old_dispatches_from_kv"
	cronjobReconcileUserInfoState       = "centrum.housekeeper.reconcile_userinfo_state"
	cronjobCleanupDispatchers           = "centrum.housekeeper.cleanup_dispatchers"

	legacyCronjobLoadNewDispatches = "centrum.housekeeper.load_new_dispatches"
)

var legacyManagerCronjobs = []string{
	"centrum.manager_housekeeper.dispatch_deduplication",
	"centrum.manager_housekeeper.load_new_dispatches",
	"centrum.manager_housekeeper.dispatch_assignment_expiration",
	"centrum.manager_housekeeper.cleanup_units",
	"centrum.manager_housekeeper.audit_unit_membership",
	"centrum.manager_housekeeper.audit_empty_unit_dispatches",
	"centrum.manager_housekeeper.cancel_old_dispatches",
	"centrum.manager_housekeeper.delete_old_dispatches",
	"centrum.manager_housekeeper.delete_old_dispatches_from_kv",
	"centrum.manager_housekeeper.reconcile_userinfo_state",
}

type housekeeperCronjob struct {
	name     string
	schedule string
	timeout  time.Duration
	handler  croner.CronjobHandlerFn
}

func (s *Housekeeper) cronjobs() []housekeeperCronjob {
	return []housekeeperCronjob{
		{
			name:     cronjobDispatchAssignmentExpiration,
			schedule: "*/2 * * * * * *",
			timeout:  3 * time.Second,
			handler:  s.runHandleDispatchAssignmentExpiration,
		},
		{
			name:     cronjobCleanupUnits,
			schedule: "15 * * * *",
			timeout:  6 * time.Second,
			handler:  s.runCleanupUnits,
		},
		{
			name:     cronjobReconcileUserInfoState,
			schedule: reconcileUserInfoStateSchedule,
			timeout:  2 * time.Minute,
			handler:  s.runReconcileUserInfoState,
		},
		{
			name:     cronjobAuditEmptyUnitDispatches,
			schedule: auditEmptyUnitDispatchesSchedule,
			timeout:  30 * time.Second,
			handler:  s.runAuditEmptyUnitDispatches,
		},
		{
			name:     cronjobAuditUnitMembership,
			schedule: "30 3 * * *",
			timeout:  30 * time.Second,
			handler:  s.runAuditUnitMembership,
		},
		{
			name:     cronjobCancelOldDispatches,
			schedule: cancelOldDispatchesSchedule,
			timeout:  30 * time.Second,
			handler:  s.runCancelOldDispatches,
		},
		{
			name:     cronjobDeleteOldDispatches,
			schedule: "@hourly",
			timeout:  30 * time.Second,
			handler:  s.runDeleteOldDispatches,
		},
		{
			name:     cronjobDeleteOldDispatchesFromKV,
			schedule: deleteOldDispatchesKVSchedule,
			timeout:  30 * time.Second,
			handler:  s.runDeleteOldDispatchesFromKV,
		},
		{
			name:     cronjobCleanupDispatchers,
			schedule: "45 3 * * *",
			timeout:  30 * time.Second,
			handler:  s.runCleanupDispatchers,
		},
	}
}

func (s *Housekeeper) registerCronjobHandlers(h *croner.Handlers) {
	for _, job := range s.cronjobs() {
		h.Add(job.name, job.handler)
	}
}

func newCronjob(name, schedule string, timeout time.Duration) *cron.Cronjob {
	return &cron.Cronjob{Name: name, Schedule: schedule, Timeout: durationpb.New(timeout)}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	for _, c := range legacyManagerCronjobs {
		if err := registry.UnregisterCronjob(ctx, c); err != nil {
			return err
		}
	}
	if err := registry.UnregisterCronjob(ctx, legacyCronjobLoadNewDispatches); err != nil {
		return err
	}
	for _, job := range s.cronjobs() {
		if err := registry.RegisterCronjob(
			ctx,
			newCronjob(job.name, job.schedule, job.timeout),
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	s.registerCronjobHandlers(h)
	return nil
}
