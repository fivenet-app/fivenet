package tracker

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/internal/modules"
	"github.com/fivenet-app/fivenet/internal/tests/servers"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/services/centrum/centrumstate"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	// Enable ESX compatibility for database tables
	tables.EnableESXCompat()

	code := m.Run()
	os.Exit(code)
}

func TestRefreshUserLocations(t *testing.T) {
	dbServer := servers.NewDBServer(t, true)
	natsServer := servers.NewNATSServer(t, true)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var manager *Manager
	app := fxtest.New(t,
		modules.GetFxTestOpts(
			dbServer.FxProvide(),
			natsServer.FxProvide(),
			centrumstate.StateModule,
			fx.Provide(NewManager),
			fx.Invoke(func(m *Manager) {
				manager = m
			}),
		)...,
	)
	require.NotNil(t, app)

	app.RequireStart()
	defer app.RequireStop()
	require.NotNil(t, manager)

	msgCh := make(chan int)
	consumer, err := manager.js.CreateConsumer(ctx, StreamName, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		require.NoError(t, err)
	}
	assert.NotNil(t, consumer)

	eventCount := 0
	jsCons, err := consumer.Consume(func(msg jetstream.Msg) {
		eventCount++

		if err := msg.Ack(); err != nil {
			manager.logger.Error("failed to ack message", zap.Error(err))
		}

		dest := &livemap.UsersUpdateEvent{}
		if err := proto.Unmarshal(msg.Data(), dest); err != nil {
			manager.logger.Error("failed to unmarshal nats user update response", zap.Error(err))
			return
		}
		assert.NoError(t, err)

		msgCh <- eventCount
	})
	require.NoError(t, err)
	require.NotNil(t, jsCons)
	defer jsCons.Stop()

	// Run the refreshUserLocations method to make sure the database state has been loaded
	err = manager.refreshUserLocations(ctx)
	assert.NoError(t, err)

	list := manager.userStore.List()
	assert.Len(t, list, 0)

	db, err := dbServer.DB()
	assert.NoError(t, err)

	// Insert user locations
	assert.NoError(t, insertUserLocation(ctx, db, "char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4", "ambulance", 1.0, 1.0, false))
	assert.NoError(t, insertUserLocation(ctx, db, "char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57", "ambulance", 1.0, 1.0, true))

	// Wait for users to appear (an event is sent for this)
	err = retry.Do(ctx, retry.WithMaxRetries(10, retry.NewConstant(1*time.Second)), func(ctx context.Context) error {
		select {
		case <-msgCh:
			return nil

		case <-time.After(1 * time.Second):
			return retry.RetryableError(fmt.Errorf("no user event received"))
		}
	})
	require.NoError(t, err)

	user1, ok := manager.userStore.Get(userIdKey(int32(1)))
	assert.True(t, ok)
	assert.NotNil(t, user1)
	assert.Equal(t, 1.0, user1.X)
	assert.Equal(t, 1.0, user1.Y)

	user2, ok := manager.userStore.Get(userIdKey(int32(2)))
	assert.True(t, ok)
	assert.NotNil(t, user2)
	assert.Equal(t, 1.0, user2.X)
	assert.Equal(t, 1.0, user2.Y)

	list = manager.userStore.List()
	assert.Len(t, list, 2)

	// Update user location (no event is sent for updates)
	assert.NoError(t, insertUserLocation(ctx, db, "char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57", "ambulance", 5.0, 5.0, true))

	// Wait for user2 to be updated
	err = retry.Do(ctx, retry.WithMaxRetries(10, retry.NewConstant(1*time.Second)), func(ctx context.Context) error {
		user2, ok := manager.userStore.Get(userIdKey(int32(2)))
		if !ok {
			return fmt.Errorf("user2 is nil in retry")
		}

		if user2.X == 5.0 && user2.Y == 5.0 {
			return nil
		}

		return retry.RetryableError(fmt.Errorf("user2 location not updated"))
	})
	require.NoError(t, err)

	user1, ok = manager.userStore.Get(userIdKey(int32(1)))
	assert.True(t, ok)
	assert.NotNil(t, user1)
	assert.Equal(t, 1.0, user1.X)
	assert.Equal(t, 1.0, user1.Y)

	user2, ok = manager.userStore.Get(userIdKey(int32(2)))
	assert.True(t, ok)
	assert.NotNil(t, user2)
	assert.Equal(t, 5.0, user2.X)
	assert.Equal(t, 5.0, user2.Y)

	assert.NoError(t, removeUserLocations(ctx, db))

	// Wait for users to be removed (it takes at least 15 seconds from the updatedAt time of each user location)
	err = retry.Do(ctx, retry.WithMaxRetries(45, retry.NewConstant(1*time.Second)), func(ctx context.Context) error {
		list := manager.userStore.List()
		if len(list) == 0 {
			return nil
		}

		stmt := tLocs.SELECT(jet.COUNT(tLocs.Identifier).AS("total_count"))
		var dest database.DataCount
		if err := stmt.QueryContext(ctx, db, &dest); err != nil {
			return err
		}

		return retry.RetryableError(fmt.Errorf("user list isn't empty yet. count %d", dest.TotalCount))
	})
	require.NoError(t, err)

	list = manager.userStore.List()
	assert.Len(t, list, 0)

	user1, ok = manager.userStore.Get(userIdKey(int32(1)))
	assert.False(t, ok)
	assert.Nil(t, user1)

	user2, ok = manager.userStore.Get(userIdKey(int32(2)))
	assert.False(t, ok)
	assert.Nil(t, user2)
}

func insertUserLocation(ctx context.Context, db *sql.DB, identifier string, job string, x float64, y float64, hidden bool) error {
	stmt := tLocs.
		INSERT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.X,
			tLocs.Y,
			tLocs.Hidden,
		).
		VALUES(
			identifier,
			job,
			x,
			y,
			hidden,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tLocs.Job.SET(jet.StringExp(jet.Raw("VALUES(`job`)"))),
			tLocs.X.SET(jet.FloatExp(jet.Raw("VALUES(`x`)"))),
			tLocs.Y.SET(jet.FloatExp(jet.Raw("VALUES(`y`)"))),
			tLocs.Hidden.SET(jet.BoolExp(jet.Raw("VALUES(`hidden`)"))),
		)

	_, err := stmt.ExecContext(ctx, db)

	return err
}

func removeUserLocations(ctx context.Context, db *sql.DB) error {
	stmt := tLocs.
		DELETE().
		WHERE(tLocs.Identifier.IS_NOT_NULL().OR(tLocs.Identifier.IS_NULL()))

	_, err := stmt.ExecContext(ctx, db)

	return err
}
