package dbutils

import (
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/go-jet/jet/v2/mysql"
)

// TimestampToMySQL converts a google.protobuf.Timestamp (wrapped in our resources.timestamp.Timestamp type)
// to a Jet mysql expression that can be NULL.
func TimestampToMySQL(ts *timestamp.Timestamp) mysql.DateTimeExpression {
	if ts == nil {
		return mysql.DateTimeExp(mysql.NULL)
	}

	return mysql.DateTimeT(ts.AsTime())
}

// TimestampToMySQLDateTime converts a google.protobuf.Timestamp (wrapped in our resources.timestamp.Timestamp type)
// to a Jet mysql expression that can be NULL.
func TimestampToMySQLDateTime(ts *timestamp.Timestamp) mysql.DateTimeExpression {
	if ts == nil {
		return mysql.DateTimeExp(mysql.NULL)
	}

	return mysql.DateTimeT(ts.AsTime())
}

// TimestampToMySQLDateTimeSec converts a google.protobuf.Timestamp (wrapped in our resources.timestamp.Timestamp type)
// to a Jet mysql expression that can be NULL.
func TimestampToMySQLDateTimeSec(ts *timestamp.Timestamp) mysql.DateTimeExpression {
	if ts == nil {
		return mysql.DateTimeExp(mysql.NULL)
	}

	// Truncate to seconds
	return mysql.DateTimeT(ts.AsTime().Truncate(time.Second))
}

// Int64P helper for nullable int64 (IDs), assumes that 0 is null.
func Int64P(v int64) mysql.IntegerExpression {
	if v == 0 {
		return mysql.IntExp(mysql.NULL)
	}
	return mysql.Int64(v)
}

// Int32P helper for nullable int32, assumes that 0 is null.
func Int32P(v int32) mysql.IntegerExpression {
	if v == 0 {
		return mysql.IntExp(mysql.NULL)
	}
	return mysql.Int32(v)
}

// StringPP helper for nullable string pointers. Nil pointers are treated as NULL, non-nil pointers are dereferenced.
func StringPP(s *string) mysql.StringExpression {
	if s == nil {
		return mysql.StringExp(mysql.NULL)
	}
	return mysql.String(*s)
}
