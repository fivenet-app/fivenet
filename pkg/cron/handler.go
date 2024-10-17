package cron

import (
	"sync"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
)

type CronjobHandlerFn func(data *cron.CronjobData) error

var (
	factoriesMu  = sync.Mutex{}
	cronHandlers = map[string]CronjobHandlerFn{}
)

func AddHandler(name string, fn CronjobHandlerFn) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()

	cronHandlers[name] = fn
}

func getCronjobHandler(name string) CronjobHandlerFn {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()

	return cronHandlers[name]
}
