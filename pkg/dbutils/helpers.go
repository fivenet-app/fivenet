package dbutils

import jet "github.com/go-jet/jet/v2/mysql"

type Columns jet.ProjectionList

func (c Columns) Get() jet.ProjectionList {
	out := jet.ProjectionList{}

	for i := range c {
		if c[i] != nil {
			out = append(out, c[i])
		}
	}

	return out
}

// JSON_CONTAINS is a helper function to create a JSON_CONTAINS expression in go-jet.
//
//nolint:revive // Function name is all uppercase to be consistent with go-jet package.
func JSON_CONTAINS(column jet.Column, value jet.Expression) jet.Expression {
	return jet.CustomExpression(
		jet.Token("JSON_CONTAINS("),
		column,
		jet.Token(", "),
		value,
		jet.Token(")"),
	)
}
