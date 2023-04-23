package rector

import (
	"context"
	"errors"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/proto/resources/jobs"
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
		).
		LIMIT(1)

	resp := &GetJobPropsResponse{
		JobProps: &jobs.JobProps{
			Job:                job,
			Theme:              "default",
			LivemapMarkerColor: jobs.DefaultLivemapMarkerColor,
		},
	}
	if err := stmt.QueryContext(ctx, s.db, resp.JobProps); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return resp, nil
}
func (s *Server) SetJobProps(ctx context.Context, req *SetJobPropsRequest) (*SetJobPropsResponse, error) {
	_, job, _ := auth.GetUserInfoFromContext(ctx)
	// Ensure that the job is the user's job
	req.JobProps.Job = job

	stmt := jobProps.
		INSERT(
			jobProps.Job,
			jobProps.Theme,
			jobProps.LivemapMarkerColor,
		).
		VALUES(
			req.JobProps.Job,
			req.JobProps.Theme,
			req.JobProps.LivemapMarkerColor,
		).
		ON_DUPLICATE_KEY_UPDATE(
			jobProps.Theme.SET(jet.String(req.JobProps.Theme)),
			jobProps.LivemapMarkerColor.SET(jet.String(req.JobProps.LivemapMarkerColor)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &SetJobPropsResponse{}, nil
}
