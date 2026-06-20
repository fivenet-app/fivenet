package dbutils

import (
	"github.com/go-jet/jet/v2/mysql"
)

// YEAR is a helper function to create a YEAR expression in go-jet.
func YEAR(column mysql.Column) mysql.Expression {
	return mysql.CustomExpression(
		mysql.Token("YEAR("),
		column,
		mysql.Token(")"),
	)
}

// WEEK is a helper function to create a WEEK expression in go-jet.
func WEEK(column mysql.Column) mysql.Expression {
	return mysql.CustomExpression(
		mysql.Token("WEEK("),
		column,
		mysql.Token(")"),
	)
}

// JSON_CONTAINS is a helper function to create a JSON_CONTAINS expression in go-jet.
//
//nolint:revive // Function name is all uppercase to be consistent with go-jet package.
func JSON_CONTAINS(column mysql.Column, value mysql.Expression) mysql.BoolExpression {
	return mysql.BoolExp(mysql.CustomExpression(
		mysql.Token("JSON_CONTAINS("),
		column,
		mysql.Token(", "),
		value,
		mysql.Token(")"),
	))
}

// MATCH is a helper function to create a boolean mode MATCH expression in go-jet.
func MATCH(column mysql.Column, search mysql.Expression) mysql.BoolExpression {
	return mysql.BoolExp(mysql.CustomExpression(
		mysql.Token("MATCH("),
		column,
		mysql.Token(") AGAINST ("),
		search,
		mysql.Token(" IN BOOLEAN MODE)"),
	))
}

// LAST_INSERT_ID is a helper function to create a LAST_INSERT_ID expression in go-jet.
//
//nolint:revive // Function name is all uppercase to be consistent with go-jet package.
func LAST_INSERT_ID(column mysql.Column) mysql.IntegerExpression {
	return mysql.IntExp(mysql.CustomExpression(
		mysql.Token("LAST_INSERT_ID("),
		column,
		mysql.Token(")"),
	))
}
