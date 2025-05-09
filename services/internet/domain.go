package internet

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	internet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscalendar "github.com/fivenet-app/fivenet/v2025/services/calendar/errors"
	errorsinternet "github.com/fivenet-app/fivenet/v2025/services/internet/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) CheckDomainAvailability(ctx context.Context, req *pbinternet.CheckDomainAvailabilityRequest) (*pbinternet.CheckDomainAvailabilityResponse, error) {
	domain, err := s.getDomainByTLDAndName(ctx, s.db, req.TldId, req.Name)
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
			INNER_JOIN(tTLDs,
				tTLDs.ID.EQ(tDomains.TldID),
			).
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
			tTLDs.ID,
			tTLDs.Name,
			tTLDs.Internal,
		).
		FROM(
			tDomains.
				INNER_JOIN(tTLDs,
					tTLDs.ID.EQ(tDomains.TldID),
				),
		).
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

	domain, err := s.getDomainByTLDAndName(ctx, s.db, req.TldId, req.Name)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	tDomains := table.FivenetInternetDomains

	domainId := uint64(0)
	// Domain exists
	if domain != nil {
		if domain.CreatorId != nil && *domain.CreatorId == userInfo.UserId {
			return nil, errorsinternet.ErrDomainNotTransferable
		} else if domain.TransferCode == nil {
			return nil, errorsinternet.ErrDomainNotTransferable
		} else if req.TransferCode != nil && *domain.TransferCode != *req.TransferCode {
			return nil, errorsinternet.ErrDomainWrongTransferCode
		}

		stmt := tDomains.
			UPDATE(
				tDomains.TransferCode,
				tDomains.CreatorID,
			).
			SET(
				jet.NULL,
				userInfo.UserId,
			).
			WHERE(
				tDomains.ID.EQ(jet.Uint64(domain.Id)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}

		domainId = domain.Id
	} else {
		tld, err := s.getTLD(ctx, s.db, req.TldId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}

		// If TLD is not found or internal and user is not superuser
		if tld == nil || (tld.Internal && !userInfo.SuperUser) {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}

		stmt := tDomains.
			INSERT(
				tDomains.TldID,
				tDomains.Name,
				tDomains.Active,
				tDomains.CreatorJob,
				tDomains.CreatorID,
			).
			VALUES(
				req.TldId,
				req.Name,
				userInfo.SuperUser, // Set domain active based on if user is superuser (no approval needed)
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
		domainId = uint64(lastId)
	}

	domain, err = s.getDomainById(ctx, s.db, domainId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
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

	domain, err := s.getDomainById(ctx, s.db, req.DomainId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	// Check if user owns the domain or is superuser
	if domain == nil || ((domain.CreatorId == nil && !userInfo.SuperUser) || *domain.CreatorId != userInfo.UserId) {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}

	stmt := tDomains.
		UPDATE(
			tDomains.TransferCode,
		).
		SET(
			req.Transferable,
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
