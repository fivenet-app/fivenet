package utils

import "testing"

func TestRankBetween(t *testing.T) {
	tests := []struct {
		name  string
		lower string
		upper string
		want  string
		ok    bool
	}{
		{
			name: "empty bounds",
			ok:   true,
			want: FormatRank(RankStep),
		},
		{
			name:  "before first",
			upper: FormatRank(1000),
			ok:    true,
			want:  FormatRank(500),
		},
		{
			name:  "after last",
			lower: FormatRank(1000),
			ok:    true,
			want:  FormatRank(2000),
		},
		{
			name:  "middle",
			lower: FormatRank(1000),
			upper: FormatRank(3000),
			ok:    true,
			want:  FormatRank(2000),
		},
		{
			name:  "no room",
			lower: FormatRank(1000),
			upper: FormatRank(1001),
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := RankBetween(tt.lower, tt.upper)
			if ok != tt.ok {
				t.Fatalf("ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNextRank(t *testing.T) {
	got, err := NextRank(FormatRank(1000))
	if err != nil {
		t.Fatalf("NextRank() error = %v", err)
	}
	if got != FormatRank(2000) {
		t.Fatalf("got %q, want %q", got, FormatRank(2000))
	}
}
