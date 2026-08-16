package citizenshydrator

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
)

func TestListByUserIDSkipsDeletedUsers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	h := &Hydrator{
		db:       db,
		customDB: &config.CustomDB{},
		enricher: mstlystcdata.NewDummyUserAwareEnricher(),
	}

	mock.ExpectQuery(
		`(?s)^SELECT .*FROM fivenet_user AS user .*user\.deleted_at IS NULL.*LIMIT \?;$`,
	).
		WithArgs(int32(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	got, err := h.ListByUserID(
		t.Context(),
		nil,
		nil,
		[]int32{7},
		ResolveOpts{},
	)
	if err != nil {
		t.Fatalf("ListByUserID returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected deleted user to be skipped, got %d results", len(got))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations not met: %v", err)
	}
}
