package jobs

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/stretchr/testify/require"
)

func TestDeleteJobAssetRejectsAssetOutsideTheUserJob(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	server := &Server{
		db: db,
		store: jobsstore.New(
			db,
			&config.CustomDB{},
			access.NewJobGroupsSubjectObjectAccess(db),
		).Store,
	}

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_assets INNER JOIN fivenet_files")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"job_asset.id",
			"job_asset.created_by_user_id",
			"job_asset.display_name",
			"job_asset.file_path",
			"job_asset.byte_size",
			"job_asset.content_type",
			"job_asset.created_at",
		}))

	_, err = server.DeleteJobAsset(
		grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{Job: "police"}),
		&pbsettings.DeleteJobAssetRequest{Id: 42},
	)
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.NoError(t, mock.ExpectationsWereMet())
}
