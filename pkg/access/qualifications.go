package access

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQualifications = table.FivenetQualifications.AS("qualification_short")
	tQualiResults   = table.FivenetQualificationsResults
)

type QualificationsAccessProtoMessage[T any, V protoutils.ProtoEnum] interface {
	protoutils.ProtoMessage[T]

	GetId() uint64
	GetTargetId() uint64

	GetQualificationId() uint64
	SetQualificationId(uint64)
	GetAccess() V
	SetAccess(V)
}

type Qualifications[U any, T QualificationsAccessProtoMessage[U, V], V protoutils.ProtoEnum] struct {
	table         jet.Table
	columns       *QualificationAccessColumns
	selectTable   jet.Table
	selectColumns *QualificationAccessColumns
}

func NewQualifications[U any, T QualificationsAccessProtoMessage[U, V], V protoutils.ProtoEnum](table jet.Table, columns *QualificationAccessColumns, tableAlias jet.Table, columnsAlias *QualificationAccessColumns) *Qualifications[U, T, V] {
	return &Qualifications[U, T, V]{
		table:         table,
		columns:       columns,
		selectTable:   tableAlias,
		selectColumns: columnsAlias,
	}
}

func (a *Qualifications[U, T, V]) List(ctx context.Context, tx qrm.DB, targetId uint64) ([]T, error) {
	tQualiResults := tQualiResults.AS("qualificationresult")

	var stmt jet.SelectStatement

	userInfo, ok := auth.GetUserInfoFromContext(ctx)
	if ok {
		stmt = a.selectTable.
			SELECT(
				a.selectColumns.ID,
				a.selectColumns.CreatedAt,
				a.selectColumns.TargetID,
				a.selectColumns.Access,
				a.selectColumns.QualificationId,
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
						tQualifications.ID.EQ(a.selectColumns.QualificationId),
					).
					LEFT_JOIN(tQualiResults,
						tQualiResults.QualificationID.EQ(a.selectColumns.QualificationId),
					),
			).
			WHERE(jet.AND(
				a.selectColumns.TargetID.EQ(jet.Uint64(targetId)),
				tQualifications.DeletedAt.IS_NULL(),
				tQualiResults.DeletedAt.IS_NULL(),
				tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)),
			))
	} else {
		stmt = a.selectTable.
			SELECT(
				a.selectColumns.ID,
				a.selectColumns.CreatedAt,
				a.selectColumns.TargetID,
				a.selectColumns.Access,
				a.selectColumns.QualificationId,
				tQualifications.ID,
				tQualifications.Job,
				tQualifications.Abbreviation,
				tQualifications.Title,
			).
			FROM(
				a.selectTable.
					INNER_JOIN(tQualifications,
						tQualifications.ID.EQ(a.selectColumns.QualificationId),
					),
			).
			WHERE(jet.AND(
				a.selectColumns.TargetID.EQ(jet.Uint64(targetId)),
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

func (a *Qualifications[U, T, V]) Clear(ctx context.Context, tx qrm.DB, targetId uint64) (T, error) {
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

func (a *Qualifications[U, T, V]) Compare(ctx context.Context, tx qrm.DB, targetId uint64, in []T) (toCreate []T, toUpdate []T, toDelete []T, err error) {
	current, err := a.List(ctx, tx, targetId)
	if err != nil {
		return nil, nil, nil, err
	}

	toCreate, toUpdate, toDelete = a.compare(current, in)
	return toCreate, toUpdate, toDelete, nil
}

func (a *Qualifications[U, T, V]) compare(current, in []T) (toCreate []T, toUpdate []T, toDelete []T) {
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
				if cj.GetQualificationId() != uj.GetQualificationId() {
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
	}

	return
}

func (a *Qualifications[U, T, AccessLevel]) HandleAccessChanges(ctx context.Context, tx qrm.DB, targetId uint64, access []T) ([]T, []T, []T, error) {
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
