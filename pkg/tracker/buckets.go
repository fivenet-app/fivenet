package tracker

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	BucketUserLoc = "user_locations" // Holds UserMarker (JOB.GRADE.USER_ID)
	BucketUnitMap = "user_mappings"  // Holds UserUnitMapping (USER_ID)

	SnapshotSubject = "$KV." + BucketUserLoc + "._snapshot"

	defaultSnapEvery = 30 * time.Second
)

// ExtractUserID takes a key like "police.3.123"  ➜  123
func ExtractUserID(key string) (int32, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return 0, fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	id, err := strconv.ParseInt(key[idx+1:], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("key %q does not contain a valid numeric suffix: %w", key, err)
	}
	return int32(id), nil
}
