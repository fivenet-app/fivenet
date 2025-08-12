package store

import (
	"testing"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/tests"
	"github.com/fivenet-app/fivenet/v2025/internal/tests/nats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestBasicStoreCreateAndUse(t *testing.T) {
	_, js, shutdown, err := nats.NewInProcessNATSServer()
	if err != nil {
		t.Fatal(err)
	}
	defer shutdown()

	logger := zaptest.NewLogger(t)
	ctx := t.Context()

	bucket := "test1"
	store, err := New[tests.SimpleObject](ctx, logger, js, bucket)
	require.NoError(t, err)
	assert.NotNil(t, store)

	err = store.Start(ctx, false)
	require.NoError(t, err)

	// Check if the Jetstream KV was auto-created
	kv, err := js.KeyValue(ctx, bucket)
	require.NoError(t, err)
	assert.NotNil(t, kv)

	// Retrieve a non-existent key
	val, err := store.Get("non-existent-key")
	require.Error(t, err)
	assert.Nil(t, val)

	// Create and ensure two values are stored
	first := &tests.SimpleObject{
		Field1: "First",
		Field2: false,
	}
	err = store.Put(ctx, "first", first)
	require.NoError(t, err)

	second := &tests.SimpleObject{
		Field1: "Second",
		Field2: true,
	}
	err = store.Put(ctx, "second", second)
	require.NoError(t, err)

	keys := store.Keys("")
	assert.Len(t, keys, 2)

	list := store.List()
	assert.Len(t, list, 2)

	// Retrieved values are **always clones** so compare exported values
	firstRetrieved, err := store.Get("first")
	require.NoError(t, err)
	assert.NotNil(t, firstRetrieved)

	assert.EqualExportedValues(t, firstRetrieved, first)

	// Check if Get returns the correct result for a "locally cached" value
	secondRetrieved, err := store.Get("second")
	require.NoError(t, err)
	assert.NotNil(t, secondRetrieved)

	assert.EqualExportedValues(t, secondRetrieved, second)

	// Make sure that ComputeUpdate works as expected
	newField1Val := "Hello World!"
	err = store.ComputeUpdate(
		ctx,
		"second",
		func(key string, existing *tests.SimpleObject) (*tests.SimpleObject, bool, error) {
			existing.Field1 = newField1Val
			existing.Field2 = false
			return existing, true, nil
		},
	)
	require.NoError(t, err)

	secondRetrieved, err = store.Get("second")
	require.NoError(t, err)
	assert.NotNil(t, secondRetrieved)

	if secondRetrieved != nil {
		assert.Equal(t, newField1Val, secondRetrieved.GetField1())
		assert.False(t, secondRetrieved.GetField2())
	}

	// Check that range runs the callback 2 times
	runCount := 0
	store.Range(func(key string, value *tests.SimpleObject) bool {
		runCount++
		return true
	})
	assert.Equal(t, 2, runCount)

	// Check that deleting and finally clearing, removes values as expected
	err = store.Delete(ctx, "first")
	require.NoError(t, err)

	list = store.List()
	assert.Len(t, list, 1)

	err = store.Clear(ctx)
	require.NoError(t, err)

	list = store.List()
	assert.Empty(t, list)
}
