package dsn

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// PrepareDSN parses a MySQL DSN string, applies additional options, and ensures parse time is enabled.
// Returns the formatted DSN string or an error if parsing or applying options fails.
func PrepareDSN(inDSN string, disableLocking bool, opts ...mysql.Option) (string, error) {
	dsn, err := mysql.ParseDSN(inDSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse database DSN. %w", err)
	}

	if disableLocking {
		opts = append(opts, WithIsolationLevel("REPEATABLE-READ"))
	}

	if err := dsn.Apply(opts...); err != nil {
		return "", fmt.Errorf("failed to apply dsn option. %w", err)
	}

	// Make sure parse time is enabled
	dsn.ParseTime = true

	return dsn.FormatDSN(), nil
}

// WithMultiStatements returns a MySQL config option that enables multi-statement execution.
func WithMultiStatements() func(c *mysql.Config) error {
	return func(c *mysql.Config) error {
		c.MultiStatements = true
		return nil
	}
}

func WithIsolationLevel(level string) func(c *mysql.Config) error {
	return func(c *mysql.Config) error {
		if c.ConnectionAttributes != "" {
			c.ConnectionAttributes += ","
		}

		switch level {
		case "UNCOMMITTED":
			c.ConnectionAttributes += "transaction_isolation=\"UNCOMMITTED\""
		case "READ-COMMITTED":
			c.ConnectionAttributes += "transaction_isolation=\"READ-COMMITTED\""
		case "REPEATABLE-READ":
			c.ConnectionAttributes += "transaction_isolation=\"REPEATABLE-READ\""
		case "SERIALIZABLE":
			c.ConnectionAttributes += "transaction_isolation=\"SERIALIZABLE\""
		default:
			return fmt.Errorf("unsupported isolation level: %s", level)
		}
		return nil
	}
}
