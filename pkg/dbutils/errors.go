package dbutils

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// IsDuplicateError returns true if the error is a MySQL duplicate entry error (error code 1062).
func IsDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	// Check if it is a duplicate error by returned number/code
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return false
}
