package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	// RankStep keeps gaps between ranks large enough for a while so we can
	// insert between neighbors without immediately rebalancing.
	RankStep int64 = 1000

	// RankWidth keeps the string representation lexicographically sortable.
	RankWidth = 12
)

// FormatRank renders a rank value as a zero-padded sortable string.
func FormatRank(v int64) string {
	return fmt.Sprintf("%0*d", RankWidth, v)
}

// ParseRank converts a persisted rank string back to its numeric value.
func ParseRank(rank string) (int64, error) {
	if rank == "" {
		return 0, errors.New("empty rank")
	}

	return strconv.ParseInt(strings.TrimSpace(rank), 10, 64)
}

// NextRank returns the next sparse rank after the provided one.
func NextRank(last string) (string, error) {
	if last == "" {
		return FormatRank(RankStep), nil
	}

	v, err := ParseRank(last)
	if err != nil {
		return "", err
	}
	if v > math.MaxInt64-RankStep {
		return "", errors.New("rank overflow")
	}

	return FormatRank(v + RankStep), nil
}

// RankBetween returns a rank between lower and upper, both of which may be
// empty to indicate an unbounded side.
func RankBetween(lower, upper string) (string, bool) {
	switch {
	case lower == "" && upper == "":
		return FormatRank(RankStep), true
	case lower == "":
		u, err := ParseRank(upper)
		if err != nil {
			return "", false
		}

		candidate := u - RankStep
		if candidate > 0 && candidate < u {
			return FormatRank(candidate), true
		}

		candidate = u / 2
		if candidate > 0 && candidate < u {
			return FormatRank(candidate), true
		}

		return "", false
	case upper == "":
		l, err := ParseRank(lower)
		if err != nil {
			return "", false
		}

		if l > math.MaxInt64-RankStep {
			return "", false
		}

		candidate := l + RankStep
		if candidate > l {
			return FormatRank(candidate), true
		}

		return "", false
	default:
		l, err := ParseRank(lower)
		if err != nil {
			return "", false
		}
		u, err := ParseRank(upper)
		if err != nil {
			return "", false
		}
		if l >= u {
			return "", false
		}

		diff := u - l
		if diff <= 1 {
			return "", false
		}

		candidate := l + diff/2
		if candidate <= l || candidate >= u {
			return "", false
		}

		return FormatRank(candidate), true
	}
}
