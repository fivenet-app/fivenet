package jobs

import (
	"context"
	"errors"

	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tUserProps = table.FivenetUserProps
	tJobProps  = table.FivenetJobProps
)

func (s *Server) GetMOTD(
	ctx context.Context,
	req *pbjobs.GetMOTDRequest,
) (*pbjobs.GetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobProps.
		SELECT(
			tJobProps.Motd.AS("get_motd_response.motd"),
		).
		FROM(tJobProps).
		WHERE(
			tJobProps.Job.EQ(mysql.String(userInfo.GetJob())),
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

func (s *Server) SetMOTD(
	ctx context.Context,
	req *pbjobs.SetMOTDRequest,
) (*pbjobs.SetMOTDResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tJobProps.
		INSERT(
			tJobProps.Job,
			tJobProps.Motd,
		).
		VALUES(
			userInfo.GetJob(),
			req.GetMotd(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobProps.Motd.SET(mysql.String(req.GetMotd())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.SetMOTDResponse{
		Motd: req.GetMotd(),
	}, nil
}
