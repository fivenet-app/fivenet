package dsn

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func PrepareDSN(inDSN string, opts ...mysql.Option) (string, error) {
	dsn, err := mysql.ParseDSN(inDSN)
	if err != nil {
		return "", fmt.Errorf("failed to parse database DSN. %w", err)
	}

	if err := dsn.Apply(opts...); err != nil {
		return "", fmt.Errorf("failed to apply dsn option. %w", err)
	}

	// Make sure parse time is enabled
	dsn.ParseTime = true

	return dsn.FormatDSN(), nil
}

func WithMultiStatements() func(c *mysql.Config) error {
	return func(c *mysql.Config) error {
		c.MultiStatements = true
		return nil
	}
}
