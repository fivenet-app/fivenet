package dsn

import (
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var prepareDSNTests = []struct {
	Input          string
	DisableLocking bool
	Opts           []mysql.Option
	Expected       string
	Error          bool
}{
	{
		// Empty DSN should result in a localhost
		Input:          "",
		DisableLocking: false,
		Expected:       "tcp(127.0.0.1:3306)/?parseTime=true",
	},
	{
		// Normal DSN gets parseTime=true added
		Input:          "username:password@tcp(localhost:3306)/fivenet?charset=utf8mb4&loc=Europe%2FBerlin",
		DisableLocking: false,
		Expected:       "username:password@tcp(localhost:3306)/fivenet?charset=utf8mb4&loc=Europe%2FBerlin&parseTime=true",
	},
	{
		// Make sure parseTime is overridden
		Input:          "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=false",
		DisableLocking: false,
		Expected:       "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=true",
	},
	{
		// Make sure "override" options work
		Input:          "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=false",
		Opts:           []mysql.Option{WithMultiStatements()},
		DisableLocking: false,
		Expected:       "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&multiStatements=true&parseTime=true",
	},
	{
		// Make sure disable locking works
		Input:          "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&parseTime=false",
		Opts:           []mysql.Option{WithMultiStatements()},
		DisableLocking: true,
		Expected:       "tcp(localhost:3306)/fivenet?loc=Europe%2FBerlin&multiStatements=true&parseTime=true&transaction_isolation=%27REPEATABLE-READ%27",
	},
}

func TestPrepareDSN(t *testing.T) {
	for _, test := range prepareDSNTests {
		out, err := PrepareDSN(test.Input, test.DisableLocking, test.Opts...)
		if test.Error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		assert.Equal(t, test.Expected, out)
	}
}
