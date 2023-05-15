package audit

import (
	"context"
	"database/sql"
	"sync"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	audit = table.FivenetAuditLog
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IAuditer interface {
	Log(service string, method string, state rector.EVENT_TYPE, targetUserId int32, data any)
	AddEntry(in *model.FivenetAuditLog)
	AddEntryWithData(in *model.FivenetAuditLog, data any)
}

type AuditStorer struct {
	logger *zap.Logger
	db     *sql.DB
	ctx    context.Context
	wg     sync.WaitGroup
	input  chan *model.FivenetAuditLog
}

func New(logger *zap.Logger, db *sql.DB) *AuditStorer {
	return &AuditStorer{
		logger: logger,
		db:     db,
		ctx:    context.Background(),
		wg:     sync.WaitGroup{},
		input:  make(chan *model.FivenetAuditLog),
	}
}

func (a *AuditStorer) Start() {
	for i := 0; i < 3; i++ {
		a.wg.Add(1)
		go a.worker()
	}
}

func (a *AuditStorer) worker() {
	defer a.wg.Done()
	for {
		select {
		case <-a.ctx.Done():
			return
		case in := <-a.input:
			if err := a.store(in); err != nil {
				a.logger.Error("failed to store audit log", zap.Error(err))
			}
		}
	}
}

func (a *AuditStorer) Stop() {
	close(a.input)
	a.wg.Wait()
}

func (a *AuditStorer) Log(service string, method string, state rector.EVENT_TYPE, targetUserId int32, data any) {
	a.input <- a.createAuditLogEntry(service, method, state, targetUserId, data)
}

func (a *AuditStorer) AddEntry(in *model.FivenetAuditLog) {
	a.AddEntryWithData(in, nil)
}

func (a *AuditStorer) AddEntryWithData(in *model.FivenetAuditLog, data any) {
	in.Data = a.toJson(data)
	a.input <- in
}

func (a *AuditStorer) store(in *model.FivenetAuditLog) error {
	stmt := audit.
		INSERT(
			audit.UserID,
			audit.UserJob,
			audit.TargetUserID,
			audit.Service,
			audit.Method,
			audit.State,
			audit.Data,
		).
		MODEL(in)

	if _, err := stmt.ExecContext(a.ctx, a.db); err != nil {
		return err
	}

	return nil
}

func (a *AuditStorer) createAuditLogEntry(service string, method string, state rector.EVENT_TYPE, targetUserId int32, data any) *model.FivenetAuditLog {
	userInfo := auth.GetUserInfoFromContext(a.ctx)

	log := &model.FivenetAuditLog{
		Service: service,
		Method:  method,
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(state),
		Data:    a.toJson(data),
	}
	if targetUserId > 0 {
		log.TargetUserID = &targetUserId
	}

	return log
}

func (a *AuditStorer) toJson(data any) *string {
	if data != nil {
		data, err := json.MarshalToString(data)
		if err != nil {
			data = "Failed to marshal data"
		}
		return &data
	}

	noData := "No Data"
	return &noData
}
