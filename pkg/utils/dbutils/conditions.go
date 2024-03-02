package dbutils

import (
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

type CustomConditions struct {
	User UserConditions `yaml:"user"`
}

type UserConditions struct {
	FilterEmptyName bool `yaml:"filterEmptyName"`
}

func (c *UserConditions) GetFilter(alias string) jet.BoolExpression {
	condition := jet.Bool(true)

	tUser := table.Users.AS(alias)
	if c.FilterEmptyName {
		condition = condition.AND(jet.AND(
			tUser.Firstname.NOT_EQ(jet.String("")),
			tUser.Lastname.NOT_EQ(jet.String("")),
		))
	}

	return condition
}
