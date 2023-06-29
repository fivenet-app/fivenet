package centrum

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tUnits      = table.FivenetCentrumUnits
	tUnitsUsers = table.FivenetCentrumUnitsUsers
)

func (s *Server) ListUnits(ctx context.Context, req *ListUnitsRequest) (*ListUnitsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tUnits.
		SELECT(
			tUnits.ID,
			tUnits.Job,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Description,
			tUnits.Status,
			tUnits.Reason,
		).
		FROM(tUnits).
		WHERE(
			tUnits.Job.EQ(jet.String(userInfo.Job)),
		)

	resp := &ListUnitsResponse{}

	if err := stmt.QueryContext(ctx, s.db, &resp.Units); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) getUnit(ctx context.Context, userInfo *userinfo.UserInfo, id uint64) (*dispatch.Unit, error) {
	stmt := tUnits.
		SELECT(
			tUnits.ID,
		).
		FROM(tUnits).
		WHERE(tUnits.ID.EQ(jet.Uint64(id))).
		LIMIT(1)

	var unit dispatch.Unit
	if err := stmt.QueryContext(ctx, s.db, &unit); err != nil {
		return nil, err
	}

	return &unit, nil
}

func (s *Server) CreateUnit(ctx context.Context, req *CreateUnitRequest) (*CreateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tUnits.
		INSERT(
			tUnits.Job,
			tUnits.Name,
			tUnits.Initials,
			tUnits.Description,
			tUnits.Status,
		).
		VALUES(
			userInfo.Job,
			req.Unit.Name,
			req.Unit.Initials,
			req.Unit.Description,
			dispatch.UNIT_STATUS_UNAVAILABLE,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	unit, err := s.getUnit(ctx, userInfo, uint64(lastId))
	if err != nil {
		return nil, err
	}

	resp := &CreateUnitResponse{
		Unit: unit,
	}

	return resp, nil
}

func (s *Server) UpdateUnit(ctx context.Context, req *UpdateUnitRequest) (*UpdateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tUnits.
		UPDATE(
			tUnits.Name,
			tUnits.Initials,
			tUnits.Description,
			tUnits.Status,
		).
		SET(
			userInfo.Job,
			req.Unit.Name,
			req.Unit.Initials,
			req.Unit.Description,
			dispatch.UNIT_STATUS_UNAVAILABLE,
		).
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(userInfo.Job)),
			tUnits.Job.EQ(jet.String(userInfo.Job)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	unit, err := s.getUnit(ctx, userInfo, req.Unit.Id)
	if err != nil {
		return nil, err
	}

	resp := &UpdateUnitResponse{
		Unit: unit,
	}

	return resp, nil
}

func (s *Server) DeleteUnit(ctx context.Context, req *DeleteUnitRequest) (*DeleteUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tUnits.
		DELETE().
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(userInfo.Job)),
			tUnits.ID.EQ(jet.Uint64(req.UnitId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &DeleteUnitResponse{}, nil
}

func (s *Server) AssignUnit(ctx context.Context, req *AssignUnitRequest) (*AssignUnitResponse, error) {
	resp := &AssignUnitResponse{}

	// TODO

	return resp, nil
}

func (s *Server) StreamUnits(req *UnitStreamRequest, srv UnitService_StreamUnitsServer) error {

	// TODO

	return nil
}
