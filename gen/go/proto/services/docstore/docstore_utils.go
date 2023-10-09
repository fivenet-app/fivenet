package docstore

import (
	context "context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	permscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) listDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, contentLength int, userInfo *userinfo.UserInfo) jet.SelectStatement {
	tDocs := tDocs.AS("documentshort")
	wheres := []jet.BoolExpression{}
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocs.DeletedAt.IS_NULL(),
				jet.OR(
					jet.OR(
						tDocs.Public.IS_TRUE(),
						tDocs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					),
					jet.OR(
						jet.AND(
							tDUserAccess.Access.IS_NOT_NULL(),
							tDUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
						jet.AND(
							tDUserAccess.Access.IS_NULL(),
							tDJobAccess.Access.IS_NOT_NULL(),
							tDJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
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
		q = tDocs.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tDocs.ID,
			tDocs.CreatedAt,
			tDocs.UpdatedAt,
			tDocs.DeletedAt,
			tDocs.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDocs.Title,
			tDocs.ContentType,
			tDocs.Summary.AS("documentshort.content"),
			tDocs.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tDocs.CreatorJob,
			tDocs.State,
			tDocs.Closed,
			tDocs.Public,
		}

		if userInfo.SuperUser {
			columns = append(columns, tDocs.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if utils.InSlice(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tDocs.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tDocs.
			LEFT_JOIN(tDUserAccess,
				tDUserAccess.DocumentID.EQ(tDocs.ID).
					AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
			LEFT_JOIN(tDJobAccess,
				tDJobAccess.DocumentID.EQ(tDocs.ID).
					AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tDCategory,
				tDocs.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocs.CreatorID.EQ(tCreator.ID),
			)
	} else {
		tables = tDocs.
			LEFT_JOIN(tDCategory,
				tDocs.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocs.CreatorID.EQ(tCreator.ID),
			)
	}

	return q.
		FROM(tables).
		WHERE(
			jet.AND(
				wheres...,
			),
		).
		ORDER_BY(
			tDocs.CreatedAt.DESC(),
			tDocs.UpdatedAt.DESC(),
		)
}

func (s *Server) getDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	var wheres []jet.BoolExpression
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tDocs.DeletedAt.IS_NULL(),
				jet.OR(
					jet.OR(
						tDocs.Public.IS_TRUE(),
						tDocs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					),
					jet.OR(
						jet.AND(
							tDUserAccess.Access.IS_NOT_NULL(),
							tDUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
						),
						jet.AND(
							tDUserAccess.Access.IS_NULL(),
							tDJobAccess.Access.IS_NOT_NULL(),
							tDJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.AccessLevel_ACCESS_LEVEL_BLOCKED))),
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
		q = tDocs.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tDocs.ID,
			tDocs.CreatedAt,
			tDocs.UpdatedAt,
			tDocs.DeletedAt,
			tDocs.CategoryID,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Job,
			tDocs.Title,
			tDocs.ContentType,
			tDocs.Data,
			tDocs.Content,
			tDocs.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tDocs.CreatorJob,
			tDocs.State,
			tDocs.Closed,
			tDocs.Public,
		}

		if userInfo.SuperUser {
			columns = append(columns, tDocs.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if utils.InSlice(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tDocs.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tDocs.
			LEFT_JOIN(tDUserAccess,
				tDUserAccess.DocumentID.EQ(tDocs.ID).
					AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
			LEFT_JOIN(tDJobAccess,
				tDJobAccess.DocumentID.EQ(tDocs.ID).
					AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tDCategory,
				tDocs.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocs.CreatorID.EQ(tCreator.ID),
			)
	} else {
		tables = tDocs.
			LEFT_JOIN(tDCategory,
				tDocs.CategoryID.EQ(tDCategory.ID),
			).
			LEFT_JOIN(tCreator,
				tDocs.CreatorID.EQ(tCreator.ID),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		)).
		ORDER_BY(
			tDocs.CreatedAt.DESC(),
			tDocs.UpdatedAt.DESC(),
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

	condition := jet.AND(
		tDocs.ID.IN(ids...),
		tDocs.DeletedAt.IS_NULL(),
		jet.OR(
			tDocs.CreatorID.EQ(jet.Int32(userInfo.UserId)),
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
	)

	stmt := tDocs.
		SELECT(
			tDocs.ID,
		).
		FROM(
			tDocs.
				LEFT_JOIN(tDUserAccess,
					tDUserAccess.DocumentID.EQ(tDocs.ID).
						AND(tDUserAccess.UserID.EQ(jet.Int32(userInfo.UserId)))).
				LEFT_JOIN(tDJobAccess,
					tDJobAccess.DocumentID.EQ(tDocs.ID).
						AND(tDJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tDJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition).
		GROUP_BY(tDocs.ID).
		ORDER_BY(tDocs.ID.DESC(), tDJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (s *Server) checkIfHasAccess(levels []string, userInfo *userinfo.UserInfo, creator *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// If no levels set, assume "Own" as default
	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	// If both have the same job, the rank checks are executed, otherwise it is "just another document"
	// the user has access to
	if creator.Job != userInfo.Job {
		return true
	}

	if utils.InSlice(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if utils.InSlice(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if utils.InSlice(levels, "Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
