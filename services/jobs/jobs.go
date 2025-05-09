package jobs

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tUserProps = table.FivenetUserProps
	tJobProps  = table.FivenetJobProps
)

func (s *Server) GetMOTD(ctx context.Context, req *pbjobs.GetMOTDRequest) (*pbjobs.GetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobProps.
		SELECT(
			tJobProps.Motd.AS("getmotdresponse.motd"),
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(jet.String(userInfo.Job)),
		).
		LIMIT(1)

	resp := &pbjobs.GetMOTDResponse{}
	if err := stmt.QueryContext(ctx, s.db, resp); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) SetMOTD(ctx context.Context, req *pbjobs.SetMOTDRequest) (*pbjobs.SetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "SetMOTD",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Motd,
		).
		VALUES(
			userInfo.Job,
			req.Motd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Motd.SET(jet.String(req.Motd)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.SetMOTDResponse{
		Motd: req.Motd,
	}, nil
}
