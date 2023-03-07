package common

import (
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

type IGetUserID interface {
	GetUserID() int32
}

var CharacterBaseColumns = []jet.Projection{
	table.Users.ID,
	table.Users.Identifier,
	table.Users.Job,
	table.Users.JobGrade,
	table.Users.Firstname,
	table.Users.Lastname,
	table.Users.Dateofbirth,
	table.Users.Sex,
	table.Users.Height,
	table.Users.Visum,
	table.Users.Playtime,
}
