package dsn

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var prepareDSNTests = []struct {
	Input    string
	Opts     []mysql.Option
	Expected string
	Error    bool
}{
	{
		// Empty DSN should result in a localhost
		Input:    "",
		Expected: "tcp(127.0.0.1:3306)/?parseTime=true",
	},
	{
		// Normal DSN gets parseTime=true added
		Input:    "username:password@tcp(localhost:3306)/fivenet?charset=utf8mb4&loc=Europe%2FBerlin",
		Expected: "username:password@tcp(localhost:3306)/fivenet?charset=utf8mb4&loc=Europe%2FBerlin&parseTime=true",
	},
	{
		// Make sure parseTime is overriden
		Input:    "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=false",
		Expected: "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=true",
	},
	{
		// Make sure "override" options work
		Input:    "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=false",
		Opts:     []mysql.Option{WithMultiStatements()},
		Expected: "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&multiStatements=true&parseTime=true",
	},
}

func TestPrepareDSN(t *testing.T) {
	for _, test := range prepareDSNTests {
		out, err := PrepareDSN(test.Input, test.Opts...)
		if test.Error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, test.Expected, out)
	}
}
