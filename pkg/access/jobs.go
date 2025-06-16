package access

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// JobsAccessProtoMessage defines the interface for job access proto messages.
// It extends ProtoMessage and provides accessors for job and access fields.
type JobsAccessProtoMessage[T any, V protoutils.ProtoEnum] interface {
	protoutils.ProtoMessage[T]

	GetId() uint64
	GetTargetId() uint64
	GetJob() string

	GetMinimumGrade() int32
	SetMinimumGrade(int32)
	GetAccess() V
	SetAccess(V)
}

// Jobs provides access control logic for job-based permissions.
type Jobs[U any, T JobsAccessProtoMessage[U, V], V protoutils.ProtoEnum] struct {
	// table is the main table for job access entries.
	table jet.Table
	// columns holds column references for the main table.
	columns *JobAccessColumns
	// selectTable is the table used for select queries (may be aliased).
	selectTable jet.Table
	// selectColumns holds column references for select queries (may be aliased).
	selectColumns *JobAccessColumns
}

// NewJobs creates a new Jobs instance for job-based access control.
func NewJobs[U any, T JobsAccessProtoMessage[U, V], V protoutils.ProtoEnum](table jet.Table, columns *JobAccessColumns, tableAlias jet.Table, columnsAlias *JobAccessColumns) *Jobs[U, T, V] {
	return &Jobs[U, T, V]{
		table:         table,
		columns:       columns,
		selectTable:   tableAlias,
		selectColumns: columnsAlias,
	}
}

// List returns all job access entries for a given targetId.
func (a *Jobs[U, T, V]) List(ctx context.Context, tx qrm.DB, targetId uint64) ([]T, error) {
	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.Job,
			a.selectColumns.MinimumGrade,
		).
		FROM(a.selectTable).
		WHERE(jet.AND(
			a.selectColumns.TargetID.EQ(jet.Uint64(targetId)),
			a.selectColumns.Job.IS_NOT_NULL(),
			a.selectColumns.MinimumGrade.IS_NOT_NULL(),
		))

	var dest []T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// Clear deletes all job access entries for a given targetId.
func (a *Jobs[U, T, V]) Clear(ctx context.Context, tx qrm.DB, targetId uint64) (T, error) {
	stmt := a.table.
		DELETE().
		WHERE(
			a.columns.TargetID.EQ(jet.Uint64(targetId)),
		)

	var dest T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// Compare compares the current job access entries in the database with the provided input.
// Returns slices of entries to create, update, and delete.
func (a *Jobs[U, T, V]) Compare(ctx context.Context, tx qrm.DB, targetId uint64, in []T) (toCreate []T, toUpdate []T, toDelete []T, err error) {
	current, err := a.List(ctx, tx, targetId)
	if err != nil {
		return nil, nil, nil, err
	}

	toCreate, toUpdate, toDelete = a.compare(current, in)
	return toCreate, toUpdate, toDelete, nil
}

// compare performs a comparison between current and input job access entries.
// Returns entries to create, update, and delete. Handles matching by job and minimum grade.
func (a *Jobs[U, T, V]) compare(current, in []T) (toCreate []T, toUpdate []T, toDelete []T) {
	toCreate = []T{}
	toUpdate = []T{}
	toDelete = []T{}

	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b T) int {
		return int(a.GetId() - b.GetId())
	})

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cj := range current {
			var found T
			var foundIdx int
			for i, uj := range in {
				if cj.GetJob() != uj.GetJob() {
					continue
				}
				if cj.GetMinimumGrade() != uj.GetMinimumGrade() {
					continue
				}
				found = uj
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete = append(toDelete, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.GetMinimumGrade() != found.GetMinimumGrade() {
				cj.SetMinimumGrade(found.GetMinimumGrade())
				changed = true
			}
			if cj.GetAccess().Number() != found.GetAccess().Number() {
				cj.SetAccess(found.GetAccess())
				changed = true
			}

			if changed {
				toUpdate = append(toUpdate, cj)
			}
		}

		for i, uj := range in {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	return
}

// HandleAccessChanges applies the necessary create, update, and delete operations for job access entries.
// Returns the created, updated, and deleted entries, or an error if any operation fails.
func (a *Jobs[U, T, AccessLevel]) HandleAccessChanges(ctx context.Context, tx qrm.DB, targetId uint64, access []T) (toCreate []T, toUpdate []T, toDelete []T, err error) {
	toCreate, toUpdate, toDelete, err = a.Compare(ctx, tx, targetId, access)
	if err != nil {
		return
	}

	for _, entry := range toCreate {
		if err = a.CreateEntry(ctx, tx, targetId, entry); err != nil {
			return
		}
	}

	for _, entry := range toUpdate {
		if err = a.UpdateEntry(ctx, tx, targetId, entry); err != nil {
			return
		}
	}

	for _, entry := range toDelete {
		if err = a.DeleteEntry(ctx, tx, targetId, entry.GetId()); err != nil {
			return
		}
	}

	return
}
