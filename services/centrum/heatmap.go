package centrum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) GetDispatchHeatmap(
	ctx context.Context,
	req *pbcentrum.GetDispatchHeatmapRequest,
) (*pbcentrum.GetDispatchHeatmapResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "GetDispatchHeatmap",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	resp := &pbcentrum.GetDispatchHeatmapResponse{}

	tDispatchHeatmap := tDispatchHeatmap.AS("coords")
	stmt := tDispatchHeatmap.
		SELECT(
			tDispatchHeatmap.Max.AS("max"),
			tDispatchHeatmap.HeatJSON.AS("data"),
		).
		FROM(tDispatchHeatmap).
		WHERE(
			tDispatchHeatmap.Job.EQ(mysql.String(userInfo.GetJob())),
		).
		LIMIT(1)

	var raw struct {
		Max  int32
		Data []byte
	}
	if err := stmt.QueryContext(ctx, s.db, &raw); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if len(raw.Data) > 0 {
		resp.MaxEntries = raw.Max
		if err := json.Unmarshal(raw.Data, &resp.Entries); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

const binSize = float64(10)

// heatmapQuery is the SQL to rebuild the per‐job Y×X‐bin heatmaps.
// note how we break out the `max` identifier so the raw string literal isn’t terminated early.
var heatmapQuery = `
REPLACE INTO fivenet_centrum_dispatches_heatmaps (job, heat_json, ` + "`max`" + `)
WITH exploded AS (
    SELECT
        jt.job AS job,
        f.x,
        f.y
    FROM fivenet_centrum_dispatches AS f
    JOIN JSON_TABLE(
        f.jobs, '$[*]'
        COLUMNS (
            job VARCHAR(50) PATH '$'
        )
    ) AS jt
),
bins AS (
    SELECT
        job,
        ROUND(x/?)* ? AS x_bin,
        ROUND(y/?)* ? AS y_bin,
        COUNT(*)      AS w
    FROM exploded
    GROUP BY job, x_bin, y_bin
),
maxw AS (
    SELECT job, MAX(w) AS mw
    FROM bins
    GROUP BY job
)
SELECT
    b.job,
    JSON_ARRAYAGG(
        JSON_OBJECT('x', b.x_bin, 'y', b.y_bin, 'w', b.w)
    ) AS heat_json,
    m.mw AS ` + "`max`" + `
FROM bins AS b
JOIN maxw AS m USING (job)
GROUP BY b.job;
`

func (s *Server) generateDispatchHeatmaps(ctx context.Context) error {
	// Four placeholders → Four copies of the grid value
	args := []any{binSize, binSize, binSize, binSize}

	if _, err := s.db.ExecContext(ctx, heatmapQuery, args...); err != nil {
		return fmt.Errorf("failed to generate dispatch heatmaps. %w", err)
	}

	return nil
}
