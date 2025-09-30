package access

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

// tQualifications is the table alias for qualifications.
var (
	tQualifications = table.FivenetQualifications.AS("qualification_short")
	// tQualiResults is the table for qualification results.
	tQualiResults = table.FivenetQualificationsResults
)

// QualificationsAccessProtoMessage defines the interface for qualification access proto messages.
// It extends ProtoMessage and provides accessors for qualification and access fields.
type QualificationsAccessProtoMessage[T any, V protoutils.ProtoEnum] interface {
	protoutils.ProtoMessage[T]

	GetId() int64
	GetTargetId() int64

	GetQualificationId() int64
	SetQualificationId(id int64)
	GetAccess() V
	SetAccess(access V)
}

// Qualifications provides access control logic for qualification-based permissions.
type Qualifications[U any, T QualificationsAccessProtoMessage[U, V], V protoutils.ProtoEnum] struct {
	// table is the main table for qualification access entries.
	table mysql.Table
	// columns holds column references for the main table.
	columns *QualificationAccessColumns
	// selectTable is the table used for select queries (may be aliased).
	selectTable mysql.Table
	// selectColumns holds column references for select queries (may be aliased).
	selectColumns *QualificationAccessColumns
}

// NewQualifications creates a new Qualifications instance for qualification-based access control.
func NewQualifications[U any, T QualificationsAccessProtoMessage[U, V], V protoutils.ProtoEnum](
	table mysql.Table,
	columns *QualificationAccessColumns,
	tableAlias mysql.Table,
	columnsAlias *QualificationAccessColumns,
) *Qualifications[U, T, V] {
	return &Qualifications[U, T, V]{
		table:         table,
		columns:       columns,
		selectTable:   tableAlias,
		selectColumns: columnsAlias,
	}
}

// List returns all qualification access entries for a given targetId.
// If user info is present in context, also joins with qualification results for that user.
func (a *Qualifications[U, T, V]) List(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
) ([]T, error) {
	tQualiResults := tQualiResults.AS("qualification_result")

	var stmt mysql.SelectStatement

	userInfo, ok := auth.GetUserInfoFromContext(ctx)
	if ok {
		stmt = a.selectTable.
			SELECT(
				a.selectColumns.ID,
				a.selectColumns.TargetID,
				a.selectColumns.Access,
				a.selectColumns.QualificationID,
				tQualifications.ID,
				tQualifications.Job,
				tQualifications.Abbreviation,
				tQualifications.Title,
				tQualiResults.ID,
				tQualiResults.QualificationID,
				tQualiResults.Status,
			).
			FROM(
				a.selectTable.
					INNER_JOIN(tQualifications,
						tQualifications.ID.EQ(a.selectColumns.QualificationID),
					).
					LEFT_JOIN(tQualiResults,
						tQualiResults.QualificationID.EQ(a.selectColumns.QualificationID),
					),
			).
			WHERE(mysql.AND(
				a.selectColumns.TargetID.EQ(mysql.Int64(targetId)),
				a.selectColumns.QualificationID.IS_NOT_NULL(),
				tQualifications.DeletedAt.IS_NULL(),
				tQualiResults.DeletedAt.IS_NULL(),
				tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
			))
	} else {
		stmt = a.selectTable.
			SELECT(
				a.selectColumns.ID,
				a.selectColumns.TargetID,
				a.selectColumns.Access,
				a.selectColumns.QualificationID,
				tQualifications.ID,
				tQualifications.Job,
				tQualifications.Abbreviation,
				tQualifications.Title,
			).
			FROM(
				a.selectTable.
					INNER_JOIN(tQualifications,
						tQualifications.ID.EQ(a.selectColumns.QualificationID),
					),
			).
			WHERE(mysql.AND(
				a.selectColumns.TargetID.EQ(mysql.Int64(targetId)),
				tQualifications.DeletedAt.IS_NULL(),
			))
	}

	var dest []T
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

// Clear deletes all qualification access entries for a given targetId.
func (a *Qualifications[U, T, V]) Clear(
	ctx context.Context,
	tx qrm.DB,
	targetId int64,
) (T, error) {
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

// Compare compares the current qualification access entries in the database with the provided input.
// Returns slices of entries to create, update, and delete.
func (a *Qualifications[U, T, V]) Compare(
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

// compare performs a comparison between current and input qualification access entries.
// Returns entries to create, update, and delete. Handles matching by qualification ID and access level.
func (a *Qualifications[U, T, V]) compare(
	current, in []T,
) ([]T, []T, []T) {
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
			if cj.GetQualificationId() != uj.GetQualificationId() {
				continue
			}
			found = uj
			foundIdx = i
			break
		}
		// No match in incoming qualification access, needs to be deleted
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

// HandleAccessChanges applies the necessary create, update, and delete operations for qualification access entries.
// Returns the created, updated, and deleted entries, or an error if any operation fails.
func (a *Qualifications[U, T, AccessLevel]) HandleAccessChanges(
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
