package audit

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/proto/resources/rector"
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
	// TODO finish the queue and close the input channels
}

func (a *AuditStorer) Log(userId int32, job string, targetJob string, service string, method string, state rector.EVENT_TYPE, data map[string]interface{}) error {

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
			userId,
			job,
			targetJob,
			service,
			method,
			state,
			data,
		)

	if _, err := stmt.ExecContext(a.ctx, a.db); err != nil {
		return err
	}

	return nil
}
