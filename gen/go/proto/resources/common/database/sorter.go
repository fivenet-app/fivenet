package database

import (
	"maps"
	"slices"

	"github.com/go-jet/jet/v2/mysql"
)

// ColumnSpec can turn a direction into one or more ORDER BY clauses.
type ColumnSpec interface {
	Order(desc bool) []mysql.OrderByClause
}

// Simple column = ASC/DESC with optional NULLS LAST/FIRST.
type Column struct {
	Col       mysql.Column
	NullsLast bool // Useful for timestamp columns
}

func (c Column) Order(desc bool) []mysql.OrderByClause {
	var ob mysql.OrderByClause
	if desc {
		ob = c.Col.DESC()
	} else {
		ob = c.Col.ASC()
	}
	if c.NullsLast {
		ob = ob.NULLS_LAST()
	}
	return []mysql.OrderByClause{ob}
}

// Custom can emit arbitrary ORDER BY list for ASC vs DESC.
type Custom struct {
	Asc func() []mysql.OrderByClause
	Dec func() []mysql.OrderByClause
}

func (c Custom) Order(desc bool) []mysql.OrderByClause {
	if desc {
		return c.Dec()
	}
	return c.Asc()
}

type SpecMap map[string]ColumnSpec

// Builder configures available sortable ids and default order.
type Builder struct {
	allowed SpecMap
	// defaultOrder used when the request doesn't specify anything valid.
	defaultOrder []mysql.OrderByClause
	// fallbackTiebreaker appended last (e.g., primary key) for stable order.
	tiebreaker []mysql.OrderByClause
	// limit to avoid pathological requests
	maxColumns int
}

// New returns a new Builder. maxColumns: keep small (e.g., 3).
func New(allowed SpecMap, defaultOrder, tiebreaker []mysql.OrderByClause, maxColumns int) *Builder {
	cp := maps.Clone(allowed)
	return &Builder{
		allowed:      cp,
		defaultOrder: slices.Clone(defaultOrder),
		tiebreaker:   slices.Clone(tiebreaker),
		maxColumns:   maxColumns,
	}
}

// Build converts your proto Sort into []jet.OrderByClause safely.
func (b *Builder) Build(s *Sort) []mysql.OrderByClause {
	out := make([]mysql.OrderByClause, 0, 4)

	if s != nil && len(s.Columns) > 0 {
		count := 0
		// apply in given order; unknown ids are ignored
		for _, c := range s.Columns {
			if b.maxColumns > 0 && count >= b.maxColumns {
				break
			}
			spec, ok := b.allowed[c.Id]
			if !ok {
				continue
			}
			order := spec.Order(c.Desc)
			out = append(out, order...)
			count++
		}
	}

	if len(out) == 0 {
		out = append(out, b.defaultOrder...)
	}
	// Always append a tiebreaker as last resort for stable pagination
	out = append(out, b.tiebreaker...)
	return out
}
