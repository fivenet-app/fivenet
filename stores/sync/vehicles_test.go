package syncstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSendVehicles(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(
		db,
		zap.NewNop(),
		&config.Config{},
		&appconfig.TestConfig{},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_owned_vehicles.*ON DUPLICATE KEY UPDATE.*`).
		WithArgs(int32(3), "police", "ABC DEF1", "adder", "car").
		WillReturnResult(sqlmock.NewResult(0, 1))

	resp, err := store.SendVehicles(t.Context(), &pbsync.SendVehiclesRequest{
		Vehicles: []*resourcesvehicles.Vehicle{{
			OwnerId: new(int32(3)),
			Job:     new("police"),
			Plate:   "ABC DEF1",
			Model:   new("adder"),
			Type:    "car",
		}},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.GetRowsAffected())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteVehicles(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(
		db,
		zap.NewNop(),
		&config.Config{},
		&appconfig.TestConfig{},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	mock.ExpectExec(`(?s)DELETE FROM .*fivenet_owned_vehicles.*plate IN \(\?, \?\).*LIMIT \?.*`).
		WithArgs("ABC DEF1", "XYZ 123", int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	resp, err := store.DeleteVehicles(t.Context(), []string{"ABC DEF1", "XYZ 123"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.GetRowsAffected())
	require.NoError(t, mock.ExpectationsWereMet())
}
