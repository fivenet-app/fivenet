package audit

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/query/fivenet/table"
)

var (
	audit = table.FivenetAuditLog
)

type AuditStorer struct {
	db  *sql.DB
	ctx context.Context
}

func New(db *sql.DB) *AuditStorer {
	return &AuditStorer{
		db:  db,
		ctx: context.Background(),
	}
}

func (a *AuditStorer) Stop() {
	// TODO
}

func (a *AuditStorer) storeInDatabase() error {

	stmt := audit.
		INSERT(
			audit.UserID,
			audit.UserJob,
			audit.TargetJob,
			audit.Service,
			audit.Method,
			audit.State,
			audit.Data,
		).
		VALUES(
			"",
		)

	if _, err := stmt.ExecContext(a.ctx, a.db); err != nil {
		return err
	}

	return nil
}
