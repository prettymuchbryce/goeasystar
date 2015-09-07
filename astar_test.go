package goeasystar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// It should find a path successfully with corner cutting enabled
func TestCornerCutting(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 0, 0, 0, 0})
	grid = append(grid, []int{0, 1, 0, 0, 0})
	grid = append(grid, []int{0, 0, 1, 0, 0})
	grid = append(grid, []int{0, 0, 0, 1, 0})
	grid = append(grid, []int{0, 0, 0, 0, 1})

	pathfinder.SetGrid(grid)

	pathfinder.EnableCornerCutting()

	pathfinder.EnableDiagonals()

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(0, 0, 4, 4)

	assert.Nil(t, err)
	assert.NotNil(t, path)
	assert.Equal(t, 5, len(path))
	assert.Equal(t, 0, path[0].X)
	assert.Equal(t, 0, path[0].Y)
	assert.Equal(t, 3, path[3].X)
	assert.Equal(t, 3, path[3].Y)
}

// It should fail to find a path successfully with corner cutting disabled
func TestFailCornerCutting(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 0, 0, 0, 0})
	grid = append(grid, []int{0, 1, 0, 0, 0})
	grid = append(grid, []int{0, 0, 1, 0, 0})
	grid = append(grid, []int{0, 0, 0, 1, 0})
	grid = append(grid, []int{0, 0, 0, 0, 1})

	pathfinder.SetGrid(grid)

	pathfinder.DisableCornerCutting()

	pathfinder.EnableDiagonals()

	pathfinder.SetAcceptableTiles([]int{1})

	_, err := pathfinder.FindPath(0, 0, 4, 4)

	assert.EqualError(t, err, "path not found")
}

// It should find a path successfully
func TestSuccessfulPath(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})

	pathfinder.SetGrid(grid)

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(1, 2, 3, 2)

	assert.Nil(t, err)
	assert.NotNil(t, path)
	assert.Equal(t, 5, len(path))
	assert.Equal(t, 1, path[0].X)
	assert.Equal(t, 2, path[0].Y)
	assert.Equal(t, 2, path[2].X)
	assert.Equal(t, 3, path[2].Y)
}

// It should be able to avoid an additional point successfully
func TestAvoidAdditionalPoint(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})

	pathfinder.SetGrid(grid)

	pathfinder.AvoidAdditionalPoint(2, 3)

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(1, 2, 3, 2)

	assert.Nil(t, err)
	assert.NotNil(t, path)
	assert.Equal(t, 7, len(path))
	assert.Equal(t, 1, path[0].X)
	assert.Equal(t, 2, path[0].Y)
	assert.Equal(t, 1, path[2].X)
	assert.Equal(t, 4, path[2].Y)
}

// It should work with diagonals
func TestDiagonals(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})

	pathfinder.SetGrid(grid)

	pathfinder.EnableDiagonals()

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(0, 0, 4, 4)

	assert.Nil(t, err)
	assert.NotNil(t, path)
	assert.Equal(t, 5, len(path))
	assert.Equal(t, 0, path[0].X)
	assert.Equal(t, 0, path[0].Y)
	assert.Equal(t, 1, path[1].X)
	assert.Equal(t, 1, path[1].Y)
	assert.Equal(t, 2, path[2].X)
	assert.Equal(t, 2, path[2].Y)
	assert.Equal(t, 3, path[3].X)
	assert.Equal(t, 3, path[3].Y)
	assert.Equal(t, 4, path[4].X)
	assert.Equal(t, 4, path[4].Y)
}

// It should work in a straight line with diagonals
func TestStraightLineDiagonals(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1, 1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1, 1, 1, 0, 1, 1})

	pathfinder.SetGrid(grid)

	pathfinder.EnableDiagonals()

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(0, 0, 9, 0)

	assert.Nil(t, err)
	assert.NotNil(t, path)

	for i := 0; i < len(path); i++ {
		assert.Equal(t, 0, path[i].Y)
	}
}

// It should return empty path when start and end are the same tile
func TestSameTile(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 0, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})
	grid = append(grid, []int{1, 1, 1, 1, 1})

	pathfinder.SetGrid(grid)

	pathfinder.SetAcceptableTiles([]int{1})

	path, err := pathfinder.FindPath(1, 2, 1, 2)

	assert.Nil(t, err)
	assert.Len(t, path, 0)

}

// It should prefer straight paths when possible
func TestPreferStraightPaths(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{0, 0, 0})
	grid = append(grid, []int{0, 0, 0})
	grid = append(grid, []int{0, 0, 0})

	pathfinder.SetGrid(grid)
	pathfinder.EnableDiagonals()
	pathfinder.SetAcceptableTiles([]int{0})

	path, err := pathfinder.FindPath(0, 1, 2, 1)

	assert.Nil(t, err)
	assert.NotNil(t, path)

	assert.Equal(t, 1, path[1].X)
	assert.Equal(t, 1, path[1].Y)
}

// It should prefer diagonal paths when they are faster
func TestPreferDiagonalPaths(t *testing.T) {
	pathfinder := NewPathfinder()

	var grid [][]int
	grid = append(grid, []int{0, 0, 0, 0, 0})
	grid = append(grid, []int{0, 0, 0, 0, 0})
	grid = append(grid, []int{0, 0, 0, 0, 0})
	grid = append(grid, []int{0, 0, 0, 0, 0})
	grid = append(grid, []int{0, 0, 0, 0, 0})

	pathfinder.SetGrid(grid)
	pathfinder.EnableDiagonals()
	pathfinder.SetAcceptableTiles([]int{0})

	path, err := pathfinder.FindPath(4, 4, 2, 2)

	assert.Nil(t, err)
	assert.NotNil(t, path)
	assert.Equal(t, 3, len(path))
	assert.Equal(t, 3, path[1].X)
	assert.Equal(t, 3, path[1].Y)
}
