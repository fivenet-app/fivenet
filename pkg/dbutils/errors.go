package dbutils

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// IsDuplicateError returns true if the error is a MySQL duplicate entry error (error code 1062).
func IsDuplicateError(err error) bool {
	// Check if it is a duplicate error by returned number/code
	if mysqlErr, ok := errors.AsType[*mysql.MySQLError](err); ok {
		return mysqlErr.Number == 1062
	}

	return false
}
