package dsn

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// PrepareDSN parses a MySQL DSN string, applies additional options, and ensures parse time is enabled.
// Returns the formatted DSN string or an error if parsing or applying options fails.
func PrepareDSN(inDSN string, disableLocking bool, opts ...mysql.Option) (string, error) {
	dsn, err := mysql.ParseDSN(inDSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse database DSN. %w", err)
	}

	if err := dsn.Apply(opts...); err != nil {
		return "", fmt.Errorf("failed to apply dsn option. %w", err)
	}

	// Make sure parse time is enabled
	dsn.ParseTime = true

	outDsn := dsn.FormatDSN()
	if outDsn != "" && disableLocking {
		if !strings.Contains(outDsn, "transaction_isolation") {
			outDsn += "&transaction_isolation=%27REPEATABLE%20READ%27"
		}
	}

	return outDsn, nil
}

// WithMultiStatements returns a MySQL config option that enables multi-statement execution.
func WithMultiStatements() func(c *mysql.Config) error {
	return func(c *mysql.Config) error {
		c.MultiStatements = true
		return nil
	}
}
