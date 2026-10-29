package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUint32Saturated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		want uint32
	}{
		{name: "negative", in: -1, want: 0},
		{name: "zero", in: 0, want: 0},
		{name: "small", in: 42, want: 42},
		{name: "max int", in: math.MaxInt, want: math.MaxUint32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ToUint32Saturated(tt.in)
			assert.Equal(
				t,
				tt.want,
				got,
				"ToUint32Saturated(%d) = %d, want %d",
				tt.in,
				got,
				tt.want,
			)
		})
	}
}

func TestToUint32Checked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     int
		want   uint32
		wantOK bool
	}{
		{name: "negative", in: -1, want: 0, wantOK: false},
		{name: "zero", in: 0, want: 0, wantOK: true},
		{name: "small", in: 42, want: 42, wantOK: true},
		{name: "max int", in: math.MaxInt, want: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := ToUint32Checked(tt.in)
			assert.Equal(
				t,
				tt.want,
				got,
				"ToUint32Checked(%d) = (%d, %t), want (%d, %t)",
				tt.in,
				got,
				ok,
				tt.want,
				tt.wantOK,
			)
			assert.Equal(
				t,
				tt.wantOK,
				ok,
				"ToUint32Checked(%d) = (%d, %t), want (%d, %t)",
				tt.in,
				got,
				ok,
				tt.want,
				tt.wantOK,
			)
		})
	}
}

func TestCheckedNarrowingFromInt64(t *testing.T) {
	t.Parallel()

	points, ok := ToUint32Checked(int64(math.MaxUint32))
	assert.True(t, ok)
	assert.Equal(t, uint32(math.MaxUint32), points)

	_, ok = ToUint32Checked(int64(-1))
	assert.False(t, ok)

	grade, ok := ToInt32Checked(int64(math.MinInt32))
	assert.True(t, ok)
	assert.Equal(t, int32(math.MinInt32), grade)

	_, ok = ToInt32Checked(int64(math.MaxInt32) + 1)
	assert.False(t, ok)
}
