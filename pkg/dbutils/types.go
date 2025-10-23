package dbutils

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/go-jet/jet/v2/mysql"
)

// TimestampToMySQL converts a google.protobuf.Timestamp (wrapped in our resources.timestamp.Timestamp type)
// to a Jet mysql expression that can be NULL.
func TimestampToMySQL(ts *timestamp.Timestamp) mysql.Expression {
	if ts == nil {
		return mysql.NULL
	}

	return mysql.TimestampT(ts.AsTime())
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
