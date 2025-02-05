package internet

import (
	"context"
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
	domain, err := s.getDomainByName(ctx, s.db, req.Name)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	// Check if domain is transferable (transfer code set)
	var transferable *bool
	if domain != nil && domain.TransferCode != nil {
		boolTrue := true
		transferable = &boolTrue
	}

	return &pbinternet.CheckDomainAvailabilityResponse{
		Available:    domain == nil,
		Transferable: transferable,
	}, nil
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
			tDomains.CreatedAt,
			tDomains.UpdatedAt,
			tDomains.DeletedAt,
			tDomains.ExpiresAt,
			tDomains.TldID,
			tDomains.Name,
			tDomains.Active,
			tDomains.ApproverJob,
			tDomains.ApproverID,
			tDomains.CreatorJob,
			tDomains.CreatorID,
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

	// TODO handle domain transfers

	stmt := tDomains.
		INSERT(
			tDomains.TldID,
			tDomains.Name,
			tDomains.Active,
			tDomains.CreatorJob,
			tDomains.CreatorID,
		).
		VALUES(
			req.Name,
			false,
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

	domain, err := s.getDomainById(ctx, s.db, uint64(lastId))
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

	// TODO check if user has access to domain

	domain, err := s.getDomainById(ctx, s.db, req.Domain.Id)
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

	domain, err = s.getDomainById(ctx, s.db, uint64(lastId))
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbinternet.UpdateDomainResponse{
		Domain: domain,
	}, nil
}
