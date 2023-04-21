package audit

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/proto/resources/rector"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	audit = table.FivenetAuditLog
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IAuditer interface {
	Log(ctx context.Context, serviceDesc grpc.ServiceDesc, state rector.EVENT_TYPE, targetJob string, data interface{})
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
	a.wg.Done()
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

func (a *AuditStorer) Log(ctx context.Context, serviceDesc grpc.ServiceDesc, state rector.EVENT_TYPE, targetJob string, data interface{}) {

	// TODO don't use interceptors, use defer in the methods we want to audit log

}

func (a *AuditStorer) store(in *model.FivenetAuditLog) error {
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
		MODEL(in)

	if _, err := stmt.ExecContext(a.ctx, a.db); err != nil {
		return err
	}

	return nil
}

func (a *AuditStorer) createAuditLogEntry(ctx context.Context, fullMethod string, in interface{}) *model.FivenetAuditLog {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)
	service, method := a.getServiceAndMethodFromFullMethod(fullMethod)
	data, err := json.MarshalToString(in)
	if err != nil {
		data = "Failed to marshal request data"
	}
	return &model.FivenetAuditLog{
		Service: service,
		Method:  method,
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_UNKNOWN),
		Data:    &data,
	}
}

func (a *AuditStorer) getServiceAndMethodFromFullMethod(fm string) (string, string) {
	parts := strings.Split(fm, "/")

	return parts[1], parts[2]
}
