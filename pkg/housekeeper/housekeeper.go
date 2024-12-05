package housekeeper

import (
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
)

var Module = fx.Module("db_housekeeper",
	fx.Provide(
		New,
	),
)

type Table struct {
	Table     jet.Table
	DeletedAt jet.Column
	MinDays   int
}

type Housekeeper struct{}

type Params struct {
	fx.In
}

func New(p Params) *Housekeeper {
	h := &Housekeeper{}

	return h
}
