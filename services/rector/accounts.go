package rector

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/accounts"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	rector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbrector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsrector "github.com/fivenet-app/fivenet/v2025/services/rector/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tAccounts = table.FivenetAccounts.AS("account")

func (s *Server) ListAccounts(ctx context.Context, req *pbrector.ListAccountsRequest) (*pbrector.ListAccountsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "ListAccounts",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
	}, req)

	condition := jet.Bool(true)
	if req.License != nil {
		condition = condition.AND(tAccounts.License.LIKE(jet.String(fmt.Sprintf("%%%s%%", *req.License))))
	}

	if req.Enabled != nil {
		condition = condition.AND(tAccounts.Enabled.EQ(jet.Bool(*req.Enabled)))
	}

	countStmt := tAccounts.
		SELECT(
			jet.COUNT(tAccounts.ID).AS("datacount.totalcount"),
		).
		FROM(tAccounts).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, AuditLogPageSize)
	resp := &pbrector.ListAccountsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "license":
			column = tAccounts.License
		case "username":
			fallthrough
		case "id":
			fallthrough
		default:
			column = tAccounts.ID
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tAccounts.CreatedAt.DESC())
	}

	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.License,
			tAccounts.LastChar,
		).
		FROM(tAccounts).
		WHERE(condition).
		ORDER_BY(orderBys...).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Accounts); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	resp.Pagination.Update(len(resp.Accounts))

	return resp, nil
}

func (s *Server) getAccount(ctx context.Context, id uint64) (*accounts.Account, error) {
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
			tAccounts.CreatedAt,
			tAccounts.UpdatedAt,
			tAccounts.Enabled,
			tAccounts.Username,
			tAccounts.License,
			tAccounts.LastChar,
		).
		FROM(tAccounts).
		WHERE(
			tAccounts.ID.EQ(jet.Uint64(id)),
		)

	var account accounts.Account
	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	if account.Id == 0 {
		return nil, nil
	}

	return &account, nil
}

func (s *Server) UpdateAccount(ctx context.Context, req *pbrector.UpdateAccountRequest) (*pbrector.UpdateAccountResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "UpdateAccount",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	tAccounts := table.FivenetAccounts

	updateSets := []interface{}{}

	if req.Enabled != nil {
		updateSets = append(updateSets, tAccounts.Enabled.SET(jet.Bool(*req.Enabled)))
	}

	if req.LastChar != nil && *req.LastChar > 0 {
		updateSets = append(updateSets, tAccounts.LastChar.SET(jet.Int32(*req.LastChar)))
	}

	if len(updateSets) > 0 {
		stmt := tAccounts.
			UPDATE()
		if len(updateSets) == 1 {
			stmt = stmt.SET(updateSets[0])
		} else {
			stmt = stmt.SET(updateSets[0], updateSets[1:]...)
		}

		stmt = stmt.
			WHERE(
				tAccounts.ID.EQ(jet.Uint64(req.Id)),
			)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
		}
	}

	acc, err := s.getAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbrector.UpdateAccountResponse{
		Account: acc,
	}, nil
}

func (s *Server) DeleteAccount(ctx context.Context, req *pbrector.DeleteAccountRequest) (*pbrector.DeleteAccountResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteAccount",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	tAccounts := table.FivenetAccounts

	stmt := tAccounts.
		DELETE().
		WHERE(tAccounts.ID.EQ(jet.Uint64(req.Id))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsrector.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbrector.DeleteAccountResponse{}, nil
}
