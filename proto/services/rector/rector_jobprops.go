package rector

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	jobProps = table.FivenetJobProps
)

func (s *Server) GetJobProps(ctx context.Context, req *GetJobPropsRequest) (*GetJobPropsResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)

	stmt := jobProps.
		SELECT(
			jobProps.AllColumns,
		).
		FROM(jobProps).
		WHERE(
			jobProps.Job.EQ(jet.String(job)),
		)

	resp := &GetJobPropsResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.JobProps); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return resp, nil
}
func (s *Server) SetJobProps(ctx context.Context, req *SetJobPropsRequest) (*SetJobPropsResponse, error) {
	resp := &SetJobPropsResponse{}

	// TODO

	return resp, nil
}
