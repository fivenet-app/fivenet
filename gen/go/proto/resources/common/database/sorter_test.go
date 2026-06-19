package database

import (
	"testing"

	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/require"
)

// TestSorterBuilderWithNoColumns uses default and tiebreaker when no columns requested.
func TestSorterBuilderWithNoColumns(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}

	builder := New(
		SpecMap{
			"col1": col1,
		},
		[]mysql.OrderByClause{col1.ASC()},
		[]mysql.OrderByClause{col1.DESC()},
		"",
		0,
	)

	result := builder.Build(nil)
	require.Len(t, result, 2) // default + tiebreaker
}

// TestSorterBuilderWithRequestedColumnsIgnoresDefault skips default when cols present.
func TestSorterBuilderWithRequestedColumnsIgnoresDefault(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}
	col2 := &testColumn{id: 2}

	builder := New(
		SpecMap{
			"col1": col1,
			"col2": col2,
		},
		[]mysql.OrderByClause{col1.ASC()},
		[]mysql.OrderByClause{col1.DESC()},
		"",
		0,
	)

	result := builder.Build(&Sort{
		Columns: []*SortByColumn{
			{Id: "col1", Desc: true},
		},
	})

	// col1.DESC + tiebreaker (default skipped)
	require.Len(t, result, 2)
}

// TestSorterBuilderUnknownFallback applies unknownDefaultID on unknown columns.
func TestSorterBuilderUnknownFallback(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}
	colFallback := &testColumn{id: 99}

	builder := New(
		SpecMap{
			"col1":     col1,
			"fallback": colFallback,
		},
		[]mysql.OrderByClause{},
		[]mysql.OrderByClause{col1.DESC()},
		"fallback",
		0,
	)

	result := builder.Build(&Sort{
		Columns: []*SortByColumn{
			{Id: "unknown", Desc: true},
			{Id: "col1", Desc: false},
		},
	})

	// unknown→fallback.DESC + col1.ASC + tiebreaker
	require.Len(t, result, 3)
}

// TestSorterBuilderFallbackAppliedOnce ensures fallback only applies to first unknown.
func TestSorterBuilderFallbackAppliedOnce(t *testing.T) {
	t.Parallel()

	colFallback := &testColumn{id: 99}

	builder := New(
		SpecMap{
			"fallback": colFallback,
		},
		[]mysql.OrderByClause{},
		[]mysql.OrderByClause{colFallback.DESC()},
		"fallback",
		0,
	)

	result := builder.Build(&Sort{
		Columns: []*SortByColumn{
			{Id: "unknown1", Desc: true},
			{Id: "unknown2", Desc: false}, // skipped, fallback already applied
		},
	})

	// fallback.DESC + tiebreaker (unknown2 ignored)
	require.Len(t, result, 2)
}

// TestSorterBuilderMaxColumnsLimit stops at column limit.
func TestSorterBuilderMaxColumnsLimit(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}
	col2 := &testColumn{id: 2}
	col3 := &testColumn{id: 3}

	builder := New(
		SpecMap{
			"col1": col1,
			"col2": col2,
			"col3": col3,
		},
		[]mysql.OrderByClause{},
		[]mysql.OrderByClause{col1.DESC()},
		"",
		2, // maxColumns = 2
	)

	result := builder.Build(&Sort{
		Columns: []*SortByColumn{
			{Id: "col1", Desc: false},
			{Id: "col2", Desc: false},
			{Id: "col3", Desc: false}, // ignored due to maxColumns
		},
	})

	// col1 + col2 + tiebreaker (col3 ignored)
	require.Len(t, result, 3)
}

// TestSorterBuilderIgnoresUnknownWithoutFallback skips unknown columns if no fallback.
func TestSorterBuilderIgnoresUnknownWithoutFallback(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}

	builder := New(
		SpecMap{
			"col1": col1,
		},
		[]mysql.OrderByClause{},
		[]mysql.OrderByClause{col1.DESC()},
		"", // no fallback
		0,
	)

	result := builder.Build(&Sort{
		Columns: []*SortByColumn{
			{Id: "unknown", Desc: true},
			{Id: "col1", Desc: false},
		},
	})

	// unknown ignored, col1 + tiebreaker
	require.Len(t, result, 2)
}

// TestSorterBuilderAlwaysAppendsTiebreaker ensures tiebreaker is always present.
func TestSorterBuilderAlwaysAppendsTiebreaker(t *testing.T) {
	t.Parallel()

	col1 := &testColumn{id: 1}

	builder := New(
		SpecMap{},
		[]mysql.OrderByClause{col1.ASC()},
		[]mysql.OrderByClause{col1.DESC()},
		"",
		0,
	)

	result := builder.Build(nil)
	// Even with empty request, tiebreaker is always there
	require.Len(t, result, 2)
}

// testColumn is a stub for testing that implements ColumnSpec.
type testColumn struct {
	id int
}

func (tc *testColumn) Order(desc bool) []mysql.OrderByClause {
	if desc {
		return []mysql.OrderByClause{tc.DESC()}
	}
	return []mysql.OrderByClause{tc.ASC()}
}

func (tc *testColumn) ASC() mysql.OrderByClause {
	return mysql.IntegerColumn("test").ASC()
}

func (tc *testColumn) DESC() mysql.OrderByClause {
	return mysql.IntegerColumn("test").DESC()
}
