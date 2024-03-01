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
