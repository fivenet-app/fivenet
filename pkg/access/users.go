package access

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// UsersAccessProtoMessage defines the interface for user access proto messages.
// It extends ProtoMessage and provides accessors for user and access fields.
type UsersAccessProtoMessage[T any, V protoutils.ProtoEnum] interface {
	protoutils.ProtoMessage[T]

	GetId() int64
	GetTargetId() int64

	GetUserId() int32
	SetUserId(userId int32)
	GetAccess() V
	SetAccess(access V)
}

// Users provides access control logic for user-based permissions.
type Users[U any, T UsersAccessProtoMessage[U, V], V protoutils.ProtoEnum] struct {
	// table is the main table for user access entries.
	table mysql.Table
	// columns holds column references for the main table.
	columns *UserAccessColumns
	// selectTable is the table used for select queries (may be aliased).
	selectTable mysql.Table
	// selectColumns holds column references for select queries (may be aliased).
	selectColumns *UserAccessColumns
}

// NewUsers creates a new Users instance for user-based access control.
func NewUsers[U any, T UsersAccessProtoMessage[U, V], V protoutils.ProtoEnum](
	table mysql.Table,
	columns *UserAccessColumns,
	tableAlias mysql.Table,
	columnsAlias *UserAccessColumns,
) *Users[U, T, V] {
	return &Users[U, T, V]{
		table:         table,
		columns:       columns,
		selectTable:   tableAlias,
		selectColumns: columnsAlias,
	}
}

// List returns all user access entries for a given targetId, joining with user_short for additional user info.
func (a *Users[U, T, V]) List(ctx context.Context, tx qrm.DB, targetId int64) ([]T, error) {
	tUsers := tables.User().AS("user_short")

	stmt := a.selectTable.
		SELECT(
			a.selectColumns.ID,
			a.selectColumns.TargetID,
			a.selectColumns.Access,
			a.selectColumns.UserId,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			a.selectTable.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(a.selectColumns.UserId),
				),
		).
		WHERE(mysql.AND(
			a.selectColumns.TargetID.EQ(mysql.Int64(targetId)),
			a.selectColumns.UserId.IS_NOT_NULL(),
		))

	var dest []T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// Clear deletes all user access entries for a given targetId.
func (a *Users[U, T, V]) Clear(ctx context.Context, tx qrm.DB, targetId int64) (T, error) {
	stmt := a.table.
		DELETE().
		WHERE(
			a.columns.TargetID.EQ(mysql.Int64(targetId)),
		)

	var dest T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// Compare compares the current user access entries in the database with the provided input.
// Returns slices of entries to create, update, and delete.
func (a *Users[U, T, V]) Compare(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	in []T,
) ([]T, []T, []T, error) {
	current, err := a.List(ctx, tx, targetId)
	if err != nil {
		return nil, nil, nil, err
	}

	toCreate, toUpdate, toDelete := a.compare(current, in)
	return toCreate, toUpdate, toDelete, nil
}

// compare performs a comparison between current and input user access entries.
// Returns entries to create, update, and delete. Handles matching by user ID and access level.
func (a *Users[U, T, V]) compare(current, in []T) ([]T, []T, []T) {
	toCreate := []T{}
	toUpdate := []T{}
	toDelete := []T{}

	if len(current) == 0 {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current, func(a, b T) int {
		return int(a.GetId() - b.GetId())
	})

	foundTracker := []int{}
	for _, cj := range current {
		var found T
		var foundIdx int
		for i, uj := range in {
			if cj.GetUserId() != uj.GetUserId() {
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

	return toCreate, toUpdate, toDelete
}

// HandleAccessChanges applies the necessary create, update, and delete operations for user access entries.
// Returns the created, updated, and deleted entries, or an error if any operation fails.
func (a *Users[U, T, AccessLevel]) HandleAccessChanges(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
	access []T,
) ([]T, []T, []T, error) {
	toCreate, toUpdate, toDelete, err := a.Compare(ctx, tx, targetId, access)
	if err != nil {
		return toCreate, toUpdate, toDelete, err
	}

	for _, entry := range toCreate {
		if err := a.CreateEntry(ctx, tx, targetId, entry); err != nil {
			return toCreate, toUpdate, toDelete, err
		}
	}

	for _, entry := range toUpdate {
		if err := a.UpdateEntry(ctx, tx, targetId, entry); err != nil {
			return toCreate, toUpdate, toDelete, err
		}
	}

	for _, entry := range toDelete {
		if err := a.DeleteEntry(ctx, tx, targetId, entry.GetId()); err != nil {
			return toCreate, toUpdate, toDelete, err
		}
	}

	return toCreate, toUpdate, toDelete, nil
}
