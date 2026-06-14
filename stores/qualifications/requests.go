package qualificationsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListQualificationRequests(
	ctx context.Context,
	req *pbqualifications.ListQualificationRequestsRequest,
	userInfo *userinfo.UserInfo,
	where mysql.BoolExpression,
	includePhoneNumber bool,
) (*pbqualifications.ListQualificationRequestsResponse, error) {
	tQuali := tQuali.AS("qualification_short")
	tUser := table.FivenetUser.AS("user")
	tApprover := tUser.AS("approver")

	condition := mysql.AND(
		tQualiReq.DeletedAt.IS_NULL(),
		tQualiReq.Status.NOT_EQ(
			mysql.Int32(int32(resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED)),
		),
	)
	if where != nil {
		condition = condition.AND(where)
	}

	countColumn := mysql.Expression(tQualiReq.QualificationID)
	if req.GetUserId() != 0 {
		condition = condition.AND(mysql.AND(
			tUser.Job.EQ(mysql.String(userInfo.GetJob())),
			tQualiReq.UserID.EQ(mysql.Int32(req.GetUserId())),
		))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).
				AND(tQualiReq.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
			countColumn = mysql.DISTINCT(tQualiReq.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiReq.UserID)
		}
	}

	if len(req.GetStatus()) > 0 {
		statuses := []mysql.Expression{}
		for i := range req.GetStatus() {
			statuses = append(statuses, mysql.Int32(int32(req.GetStatus()[i])))
		}
		condition = condition.AND(tQualiReq.Status.IN(statuses...))
	} else {
		condition = condition.AND(
			tQualiReq.Status.NOT_EQ(
				mysql.Int32(int32(resqualifications.RequestStatus_REQUEST_STATUS_COMPLETED)),
			),
		)
	}

	countStmt := tQualiReq.
		SELECT(mysql.COUNT(countColumn).AS("data_count.total")).
		FROM(
			tQualiReq.
				INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiReq.QualificationID)).
				LEFT_JOIN(tUser, tQualiReq.UserID.EQ(tUser.ID)),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*resqualifications.QualificationRequest{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil && len(req.GetSort().GetColumns()) > 0 {
		for _, sc := range req.GetSort().GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "status":
				column = tQualiReq.Status
			case "approvedAt":
				column = tQualiReq.ApprovedAt
			case "createdAt":
				fallthrough
			default:
				column = tQualiReq.CreatedAt
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys, tQualiReq.CreatedAt.DESC())
	}

	columns := mysql.ProjectionList{
		tQualiReq.CreatedAt,
		tQualiReq.QualificationID,
		tQuali.ID,
		tQuali.CreatedAt,
		tQuali.UpdatedAt,
		tQuali.Job,
		tQuali.Closed,
		tQuali.Draft,
		tQuali.Public,
		tQuali.Abbreviation,
		tQuali.Title,
		tQuali.Description,
		tQualiReq.UserID,
		tUser.ID,
		tUser.Job,
		tUser.JobGrade,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Dateofbirth,
		tQualiReq.UserComment,
		tQualiReq.Status,
		tQualiReq.ApprovedAt,
		tQualiReq.ApproverComment,
		tQualiReq.ApproverID,
		tApprover.ID,
		tApprover.Job,
		tApprover.JobGrade,
		tApprover.Firstname,
		tApprover.Lastname,
		tApprover.Dateofbirth,
		tQualiReq.ApproverJob,
	}
	if includePhoneNumber {
		columns = append(columns, tUser.PhoneNumber, tApprover.PhoneNumber)
	}

	stmt := tQualiReq.
		SELECT(columns[0], columns[1:]...).
		FROM(
			tQualiReq.
				INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiReq.QualificationID)).
				LEFT_JOIN(tUser, tQualiReq.UserID.EQ(tUser.ID)).
				LEFT_JOIN(tApprover, tQualiReq.ApproverID.EQ(tApprover.ID)),
		).
		GROUP_BY(tQualiReq.QualificationID, tQualiReq.UserID).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Store) GetQualificationRequest(
	ctx context.Context,
	qualificationId int64,
	userId int32,
	includePhoneNumber bool,
) (*resqualifications.QualificationRequest, error) {
	tQuali := tQuali.AS("qualificationshort")
	tUser := table.FivenetUser.AS("user")
	tApprover := tUser.AS("approver")

	columns := mysql.ProjectionList{
		tQualiReq.CreatedAt,
		tQualiReq.DeletedAt,
		tQualiReq.QualificationID,
		tQuali.ID,
		tQuali.CreatedAt,
		tQuali.UpdatedAt,
		tQuali.Job,
		tQuali.Closed,
		tQuali.Draft,
		tQuali.Public,
		tQuali.Abbreviation,
		tQuali.Title,
		tQuali.Description,
		tQualiReq.UserID,
		tUser.ID,
		tUser.Job,
		tUser.JobGrade,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Dateofbirth,
		tQualiReq.UserComment,
		tQualiReq.Status,
		tQualiReq.ApprovedAt,
		tQualiReq.ApproverComment,
		tQualiReq.ApproverID,
		tQualiReq.ApproverJob,
		tApprover.ID,
		tApprover.Job,
		tApprover.JobGrade,
		tApprover.Firstname,
		tApprover.Lastname,
		tApprover.Dateofbirth,
	}
	if includePhoneNumber {
		columns = append(columns, tUser.PhoneNumber, tApprover.PhoneNumber)
	}

	stmt := tQualiReq.
		SELECT(columns[0], columns[1:]...).
		FROM(tQualiReq.
			INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiReq.QualificationID)).
			LEFT_JOIN(tUser, tUser.ID.EQ(tQualiReq.UserID)).
			LEFT_JOIN(tApprover, tApprover.ID.EQ(tQualiReq.ApproverID)),
		).
		GROUP_BY(tQualiReq.QualificationID, tQualiReq.UserID).
		ORDER_BY(tQualiReq.CreatedAt.DESC()).
		WHERE(mysql.AND(
			tQualiReq.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualiReq.UserID.EQ(mysql.Int32(userId)),
			tQualiReq.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var request resqualifications.QualificationRequest
	if err := stmt.QueryContext(ctx, s.db, &request); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if request.GetQualificationId() == 0 {
		return nil, nil
	}
	return &request, nil
}

func (s *Store) DeleteQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	stmt := tQualiReq.
		UPDATE(tQualiReq.DeletedAt).
		SET(mysql.CURRENT_TIMESTAMP()).
		WHERE(mysql.AND(
			tQualiReq.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualiReq.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpdateRequestStatus(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	status resqualifications.RequestStatus,
) error {
	stmt := tQualiReq.
		INSERT(
			tQualiReq.QualificationID,
			tQualiReq.UserID,
			tQualiReq.Status,
		).
		VALUES(
			qualificationId,
			userId,
			status,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tQualiReq.Status.SET(mysql.Int32(int32(status))),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) ApproveQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	req *resqualifications.QualificationRequest,
	userInfo *userinfo.UserInfo,
) error {
	stmt := tQualiReq.
		UPDATE(
			tQualiReq.Status,
			tQualiReq.ApprovedAt,
			tQualiReq.ApproverComment,
			tQualiReq.ApproverID,
			tQualiReq.ApproverJob,
		).
		SET(
			req.GetStatus(),
			mysql.CURRENT_TIMESTAMP(),
			req.ApproverComment,
			userInfo.GetUserId(),
			userInfo.GetJob(),
		).
		WHERE(mysql.AND(
			tQualiReq.QualificationID.EQ(mysql.Int64(req.GetQualificationId())),
			tQualiReq.UserID.EQ(mysql.Int32(req.GetUserId())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpsertQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	req *resqualifications.QualificationRequest,
) error {
	stmt := tQualiReq.
		INSERT(
			tQualiReq.QualificationID,
			tQualiReq.UserID,
			tQualiReq.UserComment,
			tQualiReq.Status,
		).
		VALUES(
			req.GetQualificationId(),
			req.GetUserId(),
			req.UserComment,
			resqualifications.RequestStatus_REQUEST_STATUS_PENDING,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tQualiReq.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
			tQualiReq.UserComment.SET(mysql.RawString("VALUES(`user_comment`)")),
			tQualiReq.Status.SET(mysql.Int32(int32(resqualifications.RequestStatus_REQUEST_STATUS_PENDING))),
			tQualiReq.ApprovedAt.SET(mysql.DateTimeExp(mysql.NULL)),
			tQualiReq.ApproverComment.SET(mysql.StringExp(mysql.NULL)),
			tQualiReq.ApproverID.SET(mysql.IntExp(mysql.NULL)),
			tQualiReq.ApproverJob.SET(mysql.StringExp(mysql.NULL)),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
