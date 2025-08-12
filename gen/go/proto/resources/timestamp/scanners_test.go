// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timestamp

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_TimestampValue(t *testing.T) {
	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		ts := Timestamp{Timestamp: &timestamppb.Timestamp{Seconds: 0, Nanos: 0}}
		v, err := ts.Value()
		require.NoError(t, err)
		assert.Equal(t, v, utcDate(1970, 1, 1))
	})
	t.Run("valid nil ts", func(t *testing.T) {
		t.Parallel()

		var ts *Timestamp
		v, err := ts.Value()
		require.NoError(t, err)
		assert.Nil(t, v)
	})
	t.Run("invalid ts", func(t *testing.T) {
		t.Parallel()

		ts := Timestamp{Timestamp: &timestamppb.Timestamp{Seconds: maxValidSeconds, Nanos: 0}}
		v, err := ts.Value()
		require.NoError(t, err)
		assert.Equal(t, v, utcDate(10000, 1, 1))
	})
}

func Test_TimestampScan(t *testing.T) {
	t.Parallel()
	assert, require := assert.New(t), require.New(t)
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		v := time.Unix(0, 0)
		ts := Timestamp{}
		err := ts.Scan(v)
		require.NoError(err)
		assert.True(
			reflect.DeepEqual(ts.GetTimestamp(), &timestamppb.Timestamp{Seconds: 0, Nanos: 0}),
		)
	})
	t.Run("valid default time", func(t *testing.T) {
		t.Parallel()

		var v time.Time
		ts := Timestamp{}
		err := ts.Scan(v)
		require.NoError(err)
		assert.True(
			reflect.DeepEqual(
				ts.GetTimestamp(),
				&timestamppb.Timestamp{Seconds: -62135596800, Nanos: 0},
			),
		)
	})
	t.Run("invalid type", func(t *testing.T) {
		t.Parallel()

		v := 1
		ts := Timestamp{}
		err := ts.Scan(v)
		require.Error(err)
		assert.Equal("not a protobuf Timestamp", err.Error())
	})
	t.Run("invalid time", func(t *testing.T) {
		t.Parallel()

		v := time.Unix(maxValidSeconds, 0)
		ts := Timestamp{}
		err := ts.Scan(v)
		require.NoError(err)
	})
}

const maxValidSeconds = 253402300800

func utcDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
