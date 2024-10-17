package docstore

import (
	context "context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	permscitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/perms"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) listDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	wheres := []jet.BoolExpression{}
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocumentShort.DeletedAt.IS_NULL(),
				jet.OR(
					tDocumentShort.Public.IS_TRUE(),
					tDocumentShort.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					jet.OR(
						jet.AND(
							tDUserAccess.Access.IS_NOT_NULL(),
							tDUserAccess.Access.GT(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
						jet.AND(
							tDUserAccess.Access.IS_NULL(),
							tDJobAccess.Access.IS_NOT_NULL(),
							tDJobAccess.Access.GT(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
					),
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tDocumentShort.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tDocumentShort.ID,
			tDocumentShort.CreatedAt,
			tDocumentShort.UpdatedAt,
			tDocumentShort.DeletedAt,
			tDocumentShort.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDocumentShort.Title,
			tDocumentShort.ContentType,
			tDocumentShort.Summary.AS("documentshort.content"),
			tDocumentShort.CreatorID,
			tDocumentShort.TemplateID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocumentShort.CreatorJob,
			tDocumentShort.State,
			tDocumentShort.Closed,
			tDocumentShort.Public,
			tDocumentShort.TemplateID,
		}

		if userInfo.SuperUser {
			columns = append(columns, tDocumentShort.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if slices.Contains(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tDocumentShort.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tDocumentShort.
			LEFT_JOIN(tDUserAccess,
				tDUserAccess.DocumentID.EQ(tDocumentShort.ID).
					AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tDJobAccess,
				tDJobAccess.DocumentID.EQ(tDocumentShort.ID).
					AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tDCategory,
				tDocumentShort.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocumentShort.CreatorID.EQ(tCreator.ID),
			)
	} else {
		tables = tDocumentShort.
			LEFT_JOIN(tDCategory,
				tDocumentShort.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocumentShort.CreatorID.EQ(tCreator.ID),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		))
}

func (s *Server) getDocumentQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo, withContent bool) jet.SelectStatement {
	var wheres []jet.BoolExpression
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocument.DeletedAt.IS_NULL(),
				jet.OR(
					jet.OR(
						tDocument.Public.IS_TRUE(),
						tDocument.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					),
					jet.OR(
						jet.AND(
							tDUserAccess.Access.IS_NOT_NULL(),
							tDUserAccess.Access.GT(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
						jet.AND(
							tDUserAccess.Access.IS_NULL(),
							tDJobAccess.Access.IS_NOT_NULL(),
							tDJobAccess.Access.GT(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
					),
				),
			),
		}
	}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tDocument.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tDocument.ID,
			tDocument.CreatedAt,
			tDocument.UpdatedAt,
			tDocument.DeletedAt,
			tDocument.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDCategory.Color,
			tDCategory.Icon,
			tDocument.Title,
			tDocument.ContentType,
			tDocument.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tDocument.CreatorJob,
			tDocument.State,
			tDocument.Closed,
			tDocument.Public,
			tDocument.TemplateID,
			tDPins.State.AS("document.pinned"),
		}

		if withContent {
			columns = append(columns,
				tDocument.Data,
				tDocument.Content,
			)
		}

		if userInfo.SuperUser {
			columns = append(columns, tDocument.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if slices.Contains(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tDocument.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tDocument.
			LEFT_JOIN(tDUserAccess,
				tDUserAccess.DocumentID.EQ(tDocument.ID).
					AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tDJobAccess,
				tDJobAccess.DocumentID.EQ(tDocument.ID).
					AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tDCategory,
				tDocument.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocument.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tDPins,
				tDPins.DocumentID.EQ(tDocument.ID),
			)
	} else {
		tables = tDocument.
			LEFT_JOIN(tDCategory,
				tDocument.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocument.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tDPins,
				tDPins.DocumentID.EQ(tDocument.ID),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		)).
		ORDER_BY(
			tDocument.CreatedAt.DESC(),
			tDocument.UpdatedAt.DESC(),
		)
}

func (s *Server) checkIfUserHasAccessToDoc(ctx context.Context, documentId uint64, userInfo *userinfo.UserInfo, access documents.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, access, documentId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToDocs(ctx context.Context, userInfo *userinfo.UserInfo, access documents.AccessLevel, documentIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, access, documentIds...)
	return len(out) == len(documentIds), err
}

func (s *Server) checkIfUserHasAccessToDocIDs(ctx context.Context, userInfo *userinfo.UserInfo, access documents.AccessLevel, documentIds ...uint64) ([]uint64, error) {
	if len(documentIds) == 0 {
		return documentIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return documentIds, nil
	}

	ids := make([]jet.Expression, len(documentIds))
	for i := 0; i < len(documentIds); i++ {
		ids[i] = jet.Uint64(documentIds[i])
	}

	stmt := tDocument.
		SELECT(
			tDocument.ID,
		).
		FROM(
			tDocument.
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocument.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocument.ID).
						AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(jet.AND(
			tDocument.ID.IN(ids...),
			tDocument.DeletedAt.IS_NULL(),
			jet.OR(
				tDocument.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tDocument.CreatorJob.EQ(jet.String(userInfo.Job)),
				jet.AND(
					tDUserAccess.Access.IS_NOT_NULL(),
					tDUserAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
				jet.AND(
					tDUserAccess.Access.IS_NULL(),
					tDJobAccess.Access.IS_NOT_NULL(),
					tDJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
			),
		)).
		GROUP_BY(tDocument.ID).
		ORDER_BY(tDocument.ID.DESC(), tDJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (s *Server) checkIfHasAccess(levels []string, userInfo *userinfo.UserInfo, creatorJob string, creator *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
	// and not the rank attributes checks
	if creatorJob != userInfo.Job {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if creator == nil {
		return true
	}

	// If no levels set, assume "Own" as default
	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	if slices.Contains(levels, "Any") {
		return true
	}
	if slices.Contains(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}

func (s *Server) updateDocumentOwner(ctx context.Context, tx qrm.DB, documentId uint64, userInfo *userinfo.UserInfo, newOwner *users.UserShort) error {
	stmt := tDocument.
		UPDATE(
			tDocument.CreatorID,
		).
		SET(
			newOwner.UserId,
		).
		WHERE(
			tDocument.ID.EQ(jet.Uint64(documentId)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   documentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Data: &documents.DocActivityData{
			Data: &documents.DocActivityData_OwnerChanged{
				OwnerChanged: &documents.DocOwnerChanged{
					NewOwnerId: newOwner.UserId,
					NewOwner:   newOwner,
				},
			},
		},
	}); err != nil {
		return errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	return nil
}
