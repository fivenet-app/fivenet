package coords

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestPointer struct {
	X, Y float64
}

func (t *TestPointer) Point() orb.Point {
	return orb.Point{t.X, t.Y}
}

func TestBasicFunctionality(t *testing.T) {
	cs := New[*TestPointer]()
	require.NotNil(t, cs)

	// Test 1: No points
	result := cs.Has(orb.Point{0.0, 0.0}, nil)
	assert.False(t, result)
	result = cs.Has(orb.Point{1.0, 1.0}, nil)
	assert.False(t, result)

	// Test 2: Add and check points
	point1 := &TestPointer{
		X: 123.456,
		Y: -123.456,
	}
	err := cs.Add(point1)
	assert.NoError(t, err)
	point2 := &TestPointer{
		X: 789.001,
		Y: -789.001,
	}
	err = cs.Add(point2)
	assert.NoError(t, err)

	result = cs.Has(point1, nil)
	assert.True(t, result)

	closeBy, ok := cs.Closest(point1.X, point1.Y)
	assert.NotNil(t, closeBy)
	assert.True(t, ok)
	assert.Equal(t, point1, closeBy)

	point3 := &TestPointer{
		X: 780.0,
		Y: -780.0,
	}
	// Radius is 1, no points should be returned
	closeBys := cs.KNearest(point3.Point(), 3, nil, 1)
	assert.Nil(t, closeBys)
	assert.Len(t, closeBys, 0)
	// Just crank up the radius, so we should get all 2 added points back
	closeBys = cs.KNearest(point3.Point(), 3, nil, 1500)
	assert.NotNil(t, closeBys)
	assert.Len(t, closeBys, 2)

	// Test 3: Add a new point and remove old one
	err = cs.Add(point3)
	assert.NoError(t, err)
	result = cs.Remove(point1, nil)
	assert.True(t, result)

	// We should get 2 points back
	closeBys = cs.KNearest(point3.Point(), 3, nil, 150)
	assert.NotNil(t, closeBys)
	assert.Len(t, closeBys, 2)
}
