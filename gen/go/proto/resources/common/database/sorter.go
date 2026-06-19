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

// Column contains info for sorting ASC/DESC by that column with optional NULLS LAST/FIRST.
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
	Asc  func() []mysql.OrderByClause
	Desc func() []mysql.OrderByClause
}

func (c Custom) Order(desc bool) []mysql.OrderByClause {
	if desc {
		return c.Desc()
	}
	return c.Asc()
}

type SpecMap map[string]ColumnSpec

// SorterBuilder configures available sortable ids and default order.
type SorterBuilder struct {
	allowed SpecMap
	// defaultOrder used when the request doesn't specify anything valid.
	defaultOrder []mysql.OrderByClause
	// tiebreaker appended last (e.g., primary key) for stable order.
	tiebreaker []mysql.OrderByClause
	// unknownDefaultID used when request sends unknown sort id.
	unknownDefaultID string
	// maxColumns limits sort column count to avoid pathological requests.
	maxColumns int
}

// New returns a new Builder. maxColumns: keep small (e.g., 3).
// unknownDefaultID: fallback column id for unknown sort column ids (empty string = ignore unknowns).
func New(
	allowed SpecMap,
	defaultOrder []mysql.OrderByClause,
	tiebreaker []mysql.OrderByClause,
	unknownDefaultID string,
	maxColumns int,
) *SorterBuilder {
	cp := maps.Clone(allowed)
	return &SorterBuilder{
		allowed:          cp,
		defaultOrder:     slices.Clone(defaultOrder),
		tiebreaker:       slices.Clone(tiebreaker),
		unknownDefaultID: unknownDefaultID,
		maxColumns:       maxColumns,
	}
}

// Build converts proto Sort into []mysql.OrderByClause safely.
// Order: 1. Requested columns (fallback to unknownDefaultID if unknown), 2. defaultOrder if empty, 3. tiebreaker always appended.
// unknownDefaultID applied once per unknown id (no dedup); unknown ids without default are skipped.
func (b *SorterBuilder) Build(s *Sort) []mysql.OrderByClause {
	out := make([]mysql.OrderByClause, 0, 4)
	fallbackApplied := false

	if s != nil && len(s.Columns) > 0 {
		count := 0
		// apply in given order; unknown ids use unknownDefaultID fallback if set
		for _, c := range s.Columns {
			if b.maxColumns > 0 && count >= b.maxColumns {
				break
			}
			spec, ok := b.allowed[c.Id]
			if !ok && b.unknownDefaultID != "" && !fallbackApplied {
				spec, ok = b.allowed[b.unknownDefaultID]
				fallbackApplied = true
			}
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
	// Always append tiebreaker for stable pagination
	out = append(out, b.tiebreaker...)
	return out
}
