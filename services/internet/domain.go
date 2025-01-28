package internet

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	errorscalendar "github.com/fivenet-app/fivenet/services/calendar/errors"
	errorsinternet "github.com/fivenet-app/fivenet/services/internet/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) CheckDomainAvailability(ctx context.Context, req *pbinternet.CheckDomainAvailabilityRequest) (*pbinternet.CheckDomainAvailabilityResponse, error) {
	// TODO

	return nil, nil
}

func (s *Server) ListDomains(ctx context.Context, req *pbinternet.ListDomainsRequest) (*pbinternet.ListDomainsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tCreator := tables.Users().AS("creator")

	condition := jet.Bool(true)
	if !userInfo.SuperUser {
		condition = condition.AND(
			tDomains.CreatorID.EQ(jet.Int32(userInfo.UserId)),
		)
	}

	countStmt := tDomains.
		SELECT(
			jet.COUNT(tDomains.ID).AS("datacount.totalcount"),
		).
		FROM(tDomains.
			LEFT_JOIN(tCreator,
				tDomains.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tDomains.ID).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &pbinternet.ListDomainsResponse{
		Pagination: pag,
		Domains:    []*internet.Domain{},
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tDomains.
		SELECT(
			tDomains.ID,
		).
		FROM(tDomains).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Domains); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}
	}

	resp.Pagination.Update(len(resp.Domains))

	return resp, nil
}

func (s *Server) getDomain(ctx context.Context, tx *sql.DB, id uint64) (*internet.Domain, error) {
	tCreator := tables.Users().AS("creator")

	stmt := tDomains.
		SELECT(
			tDomains.ID,
		).
		FROM(
			tDomains.
				LEFT_JOIN(tCreator,
					tDomains.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(
			tDomains.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	dest := &internet.Domain{}
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Server) RegisterDomain(ctx context.Context, req *pbinternet.RegisterDomainRequest) (*pbinternet.RegisterDomainResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbinternet.InternetService_ServiceDesc.ServiceName,
		Method:  "RegisterDomain",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tDomains.
		INSERT(
			tDomains.Active,
			tDomains.Name,
			tDomains.CreatorJob,
			tDomains.CreatorID,
		).
		VALUES(
			false,
			req.Name,
			userInfo.Job,
			userInfo.UserId,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	domain, err := s.getDomain(ctx, s.db, uint64(lastId))
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &pbinternet.RegisterDomainResponse{
		Domain: domain,
	}, nil
}

func (s *Server) UpdateDomain(ctx context.Context, req *pbinternet.UpdateDomainRequest) (*pbinternet.UpdateDomainResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbinternet.InternetService_ServiceDesc.ServiceName,
		Method:  "UpdateDomain",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	domain, err := s.getDomain(ctx, s.db, req.Domain.Id)
	if err != nil {
		return nil, err
	}

	stmt := tDomains.
		UPDATE(
			tDomains.Active,
			tDomains.Name,
			tDomains.CreatorJob,
			tDomains.CreatorID,
		).
		SET(
			false,
			req.Domain.Name,
			userInfo.Job,
			userInfo.UserId,
		).
		WHERE(
			tDomains.ID.EQ(jet.Uint64(domain.Id)),
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	domain, err = s.getDomain(ctx, s.db, uint64(lastId))
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbinternet.UpdateDomainResponse{
		Domain: domain,
	}, nil
}

func (s *Server) TransferDomain(ctx context.Context, req *pbinternet.TransferDomainRequest) (*pbinternet.TransferDomainResponse, error) {
	// TODO send notification to target user and set flag on domain

	return nil, nil
}
