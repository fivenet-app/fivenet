package utils

import (
	"math"
	"testing"
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
			if got := ToUint32Saturated(tt.in); got != tt.want {
				t.Fatalf("ToUint32Saturated(%d) = %d, want %d", tt.in, got, tt.want)
			}
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
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("ToUint32Checked(%d) = (%d, %t), want (%d, %t)", tt.in, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
