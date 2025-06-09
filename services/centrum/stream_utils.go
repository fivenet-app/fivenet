package centrum

import (
	"fmt"
	"strconv"
	"strings"
)

func kvSubjects(bucket string, jobs []string) []string {
	subs := make([]string, 0, len(jobs))
	for _, j := range jobs {
		// JetStream turns each KV key into: $KV.<BUCKET>.<key>
		subs = append(subs, fmt.Sprintf("$KV.%s.%s.*", bucket, j))
	}
	return subs
}

// extractID takes a key like "police.123"  ➜  123
func extractID(key string) (uint64, error) {
	idx := strings.LastIndexByte(key, '.')
	if idx < 0 || idx+1 >= len(key) {
		return 0, fmt.Errorf("key %q does not contain a numeric suffix", key)
	}

	return strconv.ParseUint(key[idx+1:], 10, 64)
}
