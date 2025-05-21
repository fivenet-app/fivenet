package events

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/multierr"
)

type ITransaction interface {
	Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) error
	Commit(ctx context.Context) error
}

type tx struct {
	Subject string
	Payload []byte
	Opts    []jetstream.PublishOpt
}

type Transaction struct {
	js *JSWrapper

	txs []*tx
}

func NewTransaction(js *JSWrapper) *Transaction {
	return &Transaction{
		js: js,
	}
}

func (t *Transaction) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) error {
	t.txs = append(t.txs, &tx{
		Subject: subject,
		Payload: payload,
		Opts:    opts,
	})
	return nil
}

func (t *Transaction) Commit(ctx context.Context) error {
	var errs error
	for _, tx := range t.txs {
		if _, err := t.js.Publish(ctx, tx.Subject, tx.Payload, tx.Opts...); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}
