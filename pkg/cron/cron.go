package cron

import (
	"github.com/fivenet-app/fivenet/pkg/events"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	LC fx.Lifecycle

	JS *events.JSWrapper
}

type Cron struct {
	js *events.JSWrapper

	// TODO
}

func New(p Params) (*Cron, error) {
	return &Cron{
		js: p.JS,
	}, nil
}
