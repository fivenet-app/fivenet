package housekeeper

import (
	"context"
	"fmt"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
)

type fakeUnitDispatchAssignments struct {
	update func(context.Context, *string, *int32, int64, []int64, []int64, time.Time) error
	delete func(context.Context, int64, bool) error
}

func (f fakeUnitDispatchAssignments) UpdateAssignments(
	ctx context.Context,
	creatorJob *string,
	creatorID *int32,
	dispatchID int64,
	toAdd []int64,
	toRemove []int64,
	expiresAt time.Time,
) error {
	if f.update == nil {
		return nil
	}
	return f.update(ctx, creatorJob, creatorID, dispatchID, toAdd, toRemove, expiresAt)
}

func (f fakeUnitDispatchAssignments) Delete(
	ctx context.Context,
	id int64,
	removeFromDB bool,
) error {
	if f.delete == nil {
		return nil
	}
	return f.delete(ctx, id, removeFromDB)
}

type fakeUnitUserState struct {
	sync      func(context.Context, int32) error
	reconcile func(context.Context, int32, string) (bool, error)
	rangeFn   func(context.Context, func(string, *centrumunits.Unit) bool)
}

func (f fakeUnitUserState) SyncUserUnitMapping(ctx context.Context, userID int32) error {
	if f.sync == nil {
		return nil
	}
	return f.sync(ctx, userID)
}

func (f fakeUnitUserState) ReconcileUserJobChange(
	ctx context.Context,
	userID int32,
	job string,
) (bool, error) {
	if f.reconcile == nil {
		return false, nil
	}
	return f.reconcile(ctx, userID, job)
}

func (f fakeUnitUserState) Range(ctx context.Context, fn func(string, *centrumunits.Unit) bool) {
	if f.rangeFn != nil {
		f.rangeFn(ctx, fn)
	}
}

type fakeDispatcherUserState struct {
	set     func(context.Context, string, int32, bool) error
	rangeFn func(func(string, *centrumdispatchers.Dispatchers) bool)
}

func (f fakeDispatcherUserState) SetUserState(
	ctx context.Context,
	job string,
	userID int32,
	active bool,
) error {
	if f.set == nil {
		return nil
	}
	return f.set(ctx, job, userID, active)
}

func (f fakeDispatcherUserState) Range(fn func(string, *centrumdispatchers.Dispatchers) bool) {
	if f.rangeFn != nil {
		f.rangeFn(fn)
	}
}

type fakeDispatchLifecycle struct {
	get        func(context.Context, int64) (*centrumdispatches.Dispatch, error)
	update     func(context.Context, int64, *centrumdispatches.DispatchStatus) (*centrumdispatches.DispatchStatus, error)
	updateInDB func(context.Context, int64, *centrumdispatches.DispatchStatus) (*centrumdispatches.DispatchStatus, error)
	attribute  func(context.Context, *centrumdispatches.Dispatch, centrumdispatches.DispatchAttribute) error
	delete     func(context.Context, int64, bool) error
	schedule   func(context.Context, int64, *timestamp.Timestamp) error
}

func (f fakeDispatchLifecycle) Get(
	ctx context.Context,
	id int64,
) (*centrumdispatches.Dispatch, error) {
	if f.get == nil {
		return nil, fmt.Errorf("unexpected dispatch lookup for %d", id)
	}
	return f.get(ctx, id)
}

func (f fakeDispatchLifecycle) UpdateStatus(
	ctx context.Context,
	id int64,
	status *centrumdispatches.DispatchStatus,
) (*centrumdispatches.DispatchStatus, error) {
	if f.update == nil {
		return nil, fmt.Errorf("unexpected dispatch status update for %d", id)
	}
	return f.update(ctx, id, status)
}

func (f fakeDispatchLifecycle) UpdateStatusInDB(
	ctx context.Context,
	id int64,
	status *centrumdispatches.DispatchStatus,
) (*centrumdispatches.DispatchStatus, error) {
	if f.updateInDB == nil {
		return nil, fmt.Errorf("unexpected db-only dispatch status update for %d", id)
	}
	return f.updateInDB(ctx, id, status)
}

func (f fakeDispatchLifecycle) AddAttributeToDispatch(
	ctx context.Context,
	dispatch *centrumdispatches.Dispatch,
	attribute centrumdispatches.DispatchAttribute,
) error {
	if f.attribute == nil {
		return nil
	}
	return f.attribute(ctx, dispatch, attribute)
}

func (f fakeDispatchLifecycle) Delete(ctx context.Context, id int64, removeFromDB bool) error {
	if f.delete == nil {
		return nil
	}
	return f.delete(ctx, id, removeFromDB)
}

func (f fakeDispatchLifecycle) ScheduleProjectionCleanup(
	ctx context.Context,
	id int64,
	createdAt *timestamp.Timestamp,
) error {
	if f.schedule == nil {
		return nil
	}
	return f.schedule(ctx, id, createdAt)
}

type fakeEmptyUnitCleaner struct {
	remove func(context.Context, *centrumunits.Unit) (int, error)
}

func (f fakeEmptyUnitCleaner) Remove(ctx context.Context, unit *centrumunits.Unit) (int, error) {
	if f.remove == nil {
		return 0, nil
	}
	return f.remove(ctx, unit)
}
