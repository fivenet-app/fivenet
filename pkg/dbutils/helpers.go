package dbutils

import (
	"github.com/go-jet/jet/v2/mysql"
)

type Columns mysql.ProjectionList

func (c Columns) Get() mysql.ProjectionList {
	out := mysql.ProjectionList{}

	for i := range c {
		if c[i] != nil {
			out = append(out, c[i])
		}
	}

	return out
}

func YEAR(column mysql.Column) mysql.Expression {
	return mysql.CustomExpression(
		mysql.Token("YEAR("),
		column,
		mysql.Token(")"),
	)
}

func WEEK(column mysql.Column) mysql.Expression {
	return mysql.CustomExpression(
		mysql.Token("WEEK("),
		column,
		mysql.Token(")"),
	)
}

// JSON_CONTAINS is a helper function to create a JSON_CONTAINS expression in go-mysql.
//
//nolint:revive // Function name is all uppercase to be consistent with go-jet package.
func JSON_CONTAINS(column mysql.Column, value mysql.Expression) mysql.Expression {
	return mysql.CustomExpression(
		mysql.Token("JSON_CONTAINS("),
		column,
		mysql.Token(", "),
		value,
		mysql.Token(")"),
	)
}

func MATCH(column mysql.Column, search mysql.Expression) mysql.BoolExpression {
	return mysql.BoolExp(mysql.CustomExpression(
		mysql.Token("MATCH("),
		column,
		mysql.Token(") AGAINST ("),
		search,
		mysql.Token(" IN BOOLEAN MODE)"),
	))
}
