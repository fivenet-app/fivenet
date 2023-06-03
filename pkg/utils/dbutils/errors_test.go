package dbutils

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicateError(t *testing.T) {
	for _, run := range []struct {
		input  error
		output bool
		msg    string
	}{
		{
			input:  errors.New("dummy error"),
			output: false,
			msg:    "authorization string must not be empty",
		},
		{
			input:  mysql.ErrInvalidConn,
			output: false,
			msg:    "mysql error that isn't a duplicate",
		},
		{
			input: &mysql.MySQLError{
				Number: 1062,
			},
			output: true,
			msg:    "mysql duplicate error code",
		},
		{
			input: &mysql.MySQLError{
				Number: 1000,
			},
			output: false,
			msg:    "mysql non-duplicate error code",
		},
	} {
		assert.Equal(t, run.output, IsDuplicateError(run.input), run.msg)
	}
}
