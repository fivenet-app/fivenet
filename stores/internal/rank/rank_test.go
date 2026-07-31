package rank

import (
	"context"
	"errors"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testGroup struct {
	rows []Row
}

func (g *testGroup) ListRanks(_ context.Context, _ qrm.DB, excludeID int64) ([]Row, error) {
	rows := make([]Row, 0, len(g.rows))
	for _, row := range g.rows {
		if row.ID == excludeID {
			continue
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func (g *testGroup) UpdateRank(_ context.Context, _ qrm.DB, id int64, sortRank string) error {
	for idx, row := range g.rows {
		if row.ID == id {
			g.rows[idx].SortRank = sortRank
			return nil
		}
	}

	return errors.New("row not found")
}

func TestNext(t *testing.T) {
	t.Parallel()

	group := &testGroup{
		rows: []Row{
			{ID: 1, SortRank: utils.FormatRank(1000)},
			{ID: 2, SortRank: utils.FormatRank(2000)},
		},
	}

	got, err := Next(t.Context(), nil, group, 0)
	require.NoError(t, err)
	assert.Equal(t, utils.FormatRank(3000), got)
}

func TestInsertUsesGapWhenAvailable(t *testing.T) {
	t.Parallel()

	group := &testGroup{
		rows: []Row{
			{ID: 1, SortRank: utils.FormatRank(1000)},
			{ID: 2, SortRank: utils.FormatRank(3000)},
		},
	}

	afterID := int64(1)
	got, err := Insert(t.Context(), nil, group, 0, nil, &afterID, errors.New("not found"), errors.New("failed"))
	require.NoError(t, err)
	assert.Equal(t, utils.FormatRank(2000), got)
}

func TestInsertRebalancesWhenNoGap(t *testing.T) {
	t.Parallel()

	group := &testGroup{
		rows: []Row{
			{ID: 1, SortRank: utils.FormatRank(1000)},
			{ID: 2, SortRank: utils.FormatRank(1001)},
		},
	}

	afterID := int64(1)
	got, err := Insert(t.Context(), nil, group, 0, nil, &afterID, errors.New("not found"), errors.New("failed"))
	require.NoError(t, err)
	assert.Equal(t, utils.FormatRank(1500), got)
	assert.Equal(t, utils.FormatRank(1000), group.rows[0].SortRank)
	assert.Equal(t, utils.FormatRank(2000), group.rows[1].SortRank)
}

func TestInsertReturnsNotFoundWhenNeighborIsMissing(t *testing.T) {
	t.Parallel()

	notFoundErr := errors.New("not found")
	group := &testGroup{
		rows: []Row{
			{ID: 1, SortRank: utils.FormatRank(1000)},
		},
	}

	afterID := int64(2)
	_, err := Insert(t.Context(), nil, group, 0, nil, &afterID, notFoundErr, errors.New("failed"))
	require.ErrorIs(t, err, notFoundErr)
}
