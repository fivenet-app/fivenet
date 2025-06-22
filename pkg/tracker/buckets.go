package tracker

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	BucketUserLoc         = "userloc"       // JOB.GRADE.USER_ID → UserMarker
	BucketUserMappingsMap = "user_mappings" // USER_ID           → UserMapping
	BucketUserLocByID     = "userloc_by_id" // USER_ID           → UserMapping (no Job/Grade)
)

func UserIdKey(id int32) string {
	return strconv.FormatInt(int64(id), 10)
}

func DecodeUserMarkerKey(key string) (int32, string, int32, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid user marker key: %s", key)
	}

	id, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid user marker id: %s", parts[2])
	}

	job := parts[0]
	grade, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid user marker grade: %s", parts[1])
	}

	return int32(id), job, int32(grade), nil
}
