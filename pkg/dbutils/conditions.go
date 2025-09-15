package dbutils

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/go-jet/jet/v2/mysql"
)

type CustomConditions struct {
	User UserConditions `yaml:"user"`
}

type UserConditions struct {
	FilterEmptyName bool `yaml:"filterEmptyName"`
}

func (c *UserConditions) GetFilter(alias string) mysql.BoolExpression {
	condition := mysql.Bool(true)

	tUser := tables.User().AS(alias)
	if c.FilterEmptyName {
		condition = condition.AND(mysql.AND(
			tUser.Firstname.NOT_EQ(mysql.String("")),
			tUser.Lastname.NOT_EQ(mysql.String("")),
		))
	}

	return condition
}
