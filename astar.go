package goeasystar

import (
	"errors"
	"math"
	"sort"

	"github.com/facebookgo/pqueue"
)

const straightCost = 1.0
const diagonalCost = 1.4
const costPrecision = 100

// Pathfinder represents a single instance of a pathfinding configuration
type Pathfinder struct {
	collisionGrid           [][]int
	additionalPointsToAvoid map[string]bool
	costMap                 map[int]float64
	additionalPointsToCost  map[string]float64
	allowCornerCutting      bool
	acceptableTiles         []int
	diagonalsEnabled        bool
	startX                  int
	startY                  int
	endX                    int
	endY                    int
	coordinateToNode        map[string]*node
	openList                pqueue.PriorityQueue
}

// NewPathfinder returns a new instance of a pathfinder
func NewPathfinder() *Pathfinder {
	return &Pathfinder{
		additionalPointsToAvoid: make(map[string]bool),
		costMap:                 make(map[int]float64),
		additionalPointsToCost:  make(map[string]float64),
		allowCornerCutting:      false,
		acceptableTiles:         make([]int, 0),
		diagonalsEnabled:        false,
	}
}

// SetAcceptableTiles sets a list of tiles which are deemed acceptable to
// pass through
func (p *Pathfinder) SetAcceptableTiles(t []int) {
	p.acceptableTiles = t
}

// EnableDiagonals enables diagonals on the Pathfinder
func (p *Pathfinder) EnableDiagonals() {
	p.diagonalsEnabled = true
}

// DisableDiagonals disables diagonals on the Pathfinder
func (p *Pathfinder) DisableDiagonals() {
	p.diagonalsEnabled = false
}

// SetGrid sets the grid of the Pathfinder
func (p *Pathfinder) SetGrid(grid [][]int) {
	p.collisionGrid = grid
}

// SetAdditionalPointCost sets the an additional cost for a particular point
// Overrides the cost from SetTileCost
func (p *Pathfinder) SetAdditionalPointCost(x int, y int, cost float64) {
	key := getHashKeyForPoint(x, y)
	p.additionalPointsToCost[key] = cost
}

// RemoveAdditionalPointCost removes the additional cost for a particular point
func (p *Pathfinder) RemoveAdditionalPointCost(x int, y int) {
	delete(p.additionalPointsToCost, getHashKeyForPoint(x, y))
}

// RemoveAllAdditionalPointCosts removes all additional point costs
func (p *Pathfinder) RemoveAllAdditionalPointCosts() {
	p.additionalPointsToAvoid = make(map[string]bool)
}

// AvoidAdditionalPoint avoids a particular point on the grid
// regardless of whether or not it is an acceptable tile
func (p *Pathfinder) AvoidAdditionalPoint(x int, y int) {
	key := getHashKeyForPoint(x, y)
	p.additionalPointsToAvoid[key] = true
}

// StopAvoidingAdditionalPoint stops avoiding a particular point on the grid
func (p *Pathfinder) StopAvoidingAdditionalPoint(x int, y int) {
	delete(p.additionalPointsToAvoid, getHashKeyForPoint(x, y))
}

// EnableCornerCutting enables corner cutting in diagonal movement
func (p *Pathfinder) EnableCornerCutting() {
	p.allowCornerCutting = true
}

// DisableCornerCutting disables corner cutting in diagonal movement
func (p *Pathfinder) DisableCornerCutting() {
	p.allowCornerCutting = false
}

// StopAvoidingAllAdditionalPoints stops avoiding all additional points on the grid
func (p *Pathfinder) StopAvoidingAllAdditionalPoints() {
	p.additionalPointsToAvoid = make(map[string]bool)
}

// FindPath finds a path
// Returns a slice of points which consists the starting point up to the end point inclusively
func (p *Pathfinder) FindPath(startX int, startY int, endX int, endY int) ([]*Point, error) {
	// No acceptable tiles were set
	if len(p.acceptableTiles) == 0 {
		return nil, errors.New("You can't find a path without first calling SetAcceptableTiles()")
	}

	// No grid was set
	if p.collisionGrid == nil {
		return nil, errors.New("You can't find a path without first calling SetGrid()")
	}

	// Start or endpoint outside of scope
	if startX < 0 || startY < 0 || endX < 0 || endX < 0 ||
		startX > len(p.collisionGrid[0])-1 || startY > len(p.collisionGrid)-1 ||
		endX > len(p.collisionGrid[0])-1 || endY > len(p.collisionGrid)-1 {
		return nil, errors.New("You can't find a path without first calling SetGrid()")
	}

	// Start and end are the same tile
	if startX == endX && startY == endY {
		return make([]*Point, 0), nil
	}

	// End point is not an acceptable tile
	endTile := p.collisionGrid[endY][endX]
	isAcceptable := false

	for i := 0; i < len(p.acceptableTiles); i++ {
		if endTile == p.acceptableTiles[i] {
			isAcceptable = true
			break
		}
	}

	if isAcceptable == false {
		return make([]*Point, 0), nil
	}

	p.startX = startX
	p.startY = startY
	p.endX = endX
	p.endY = endY

	// TODO is there a more memory-efficient way to do this ?
	p.openList = pqueue.New(len(p.collisionGrid[0]) * len(p.collisionGrid))
	p.coordinateToNode = make(map[string]*node)

	n := newNode(p.startX, p.startY, p.endX, p.endY, nil, straightCost)
	p.coordinateToNode[getHashKeyForPoint(startX, startY)] = n

	item := &pqueue.Item{Value: n, Priority: int64(n.bestGuessDistance() * costPrecision)}

	p.openList.Push(item)

	for {
		if p.openList.Len() == 0 {
			return nil, errors.New("path not found")
		}

		item, _ := p.openList.PeekAndShift(math.MaxInt64)
		parent := item.Value.(*node)
		var sn searchNodes
		if parent.y > 0 {
			sn = append(sn, newSearchNode(0, -1, p.endX, p.endY, straightCost*p.getTileCost(parent.x, parent.y-1), parent))
		}
		if parent.x < len(p.collisionGrid[0])-1 {
			sn = append(sn, newSearchNode(1, 0, p.endX, p.endY, straightCost*p.getTileCost(parent.x+1, parent.y), parent))
		}
		if parent.y < len(p.collisionGrid)-1 {
			sn = append(sn, newSearchNode(0, 1, p.endX, p.endY, straightCost*p.getTileCost(parent.x, parent.y+1), parent))
		}
		if parent.x > 0 {
			sn = append(sn, newSearchNode(-1, 0, p.endX, p.endY, straightCost*p.getTileCost(parent.x-1, parent.y), parent))
		}

		if p.diagonalsEnabled {
			if parent.x > 0 && parent.y > 0 {
				if p.allowCornerCutting ||
					(isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x, parent.y-1) &&
						isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x-1, parent.y)) {

					sn = append(sn, newSearchNode(-1, -1, p.endX, p.endY, diagonalCost*p.getTileCost(parent.x-1, parent.y-1), parent))
				}
			}

			if parent.x < len(p.collisionGrid[0])-1 && parent.y < len(p.collisionGrid)-1 {
				if p.allowCornerCutting ||
					(isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x, parent.y+1) &&
						isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x+1, parent.y)) {

					sn = append(sn, newSearchNode(1, 1, p.endX, p.endY, diagonalCost*p.getTileCost(parent.x+1, parent.y+1), parent))
				}
			}

			if parent.x < len(p.collisionGrid[0])-1 && parent.y > 0 {
				if p.allowCornerCutting ||
					(isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x, parent.y-1) &&
						isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x+1, parent.y)) {

					sn = append(sn, newSearchNode(1, -1, p.endX, p.endY, diagonalCost*p.getTileCost(parent.x+1, parent.y-1), parent))
				}
			}

			if parent.x > 0 && parent.y < len(p.collisionGrid)-1 {
				if p.allowCornerCutting ||
					(isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x, parent.y+1) &&
						isTileWalkable(p.collisionGrid, p.acceptableTiles, parent.x-1, parent.y)) {

					sn = append(sn, newSearchNode(-1, 1, p.endX, p.endY, diagonalCost*p.getTileCost(parent.x-1, parent.y+1), parent))
				}
			}
		}

		// First sort all of the potential nodes we could search by their cost + heuristic distance
		sort.Sort(sn)

		// Search all of the adjacent nodes
		for i := 0; i < len(sn); i++ {
			adjacentCoordinateX := sn[i].parent.x + sn[i].x
			adjacentCoordinateY := sn[i].parent.y + sn[i].y
			hashKey := getHashKeyForPoint(adjacentCoordinateX, adjacentCoordinateY)

			_, exists := p.additionalPointsToAvoid[hashKey]
			if !exists {
				if sn[i].endX == adjacentCoordinateX && sn[i].endY == adjacentCoordinateY {
					var path []*Point
					path = append(path, &Point{X: adjacentCoordinateX, Y: adjacentCoordinateY}, &Point{X: sn[i].parent.x, Y: sn[i].parent.y})
					parent := sn[i].parent
					for {
						parent = parent.parent
						if parent == nil {
							break
						}
						path = append(path, &Point{X: parent.x, Y: parent.y})
					}

					// Reverse path slice
					for i := len(path)/2 - 1; i >= 0; i-- {
						opp := len(path) - 1 - i
						path[i], path[opp] = path[opp], path[i]
					}

					return path, nil
				}

				if isTileWalkable(p.collisionGrid, p.acceptableTiles, adjacentCoordinateX, adjacentCoordinateY) {
					existingNode, exists := p.coordinateToNode[hashKey]
					if exists {
						if sn[i].parent.costSoFar+sn[i].cost < existingNode.costSoFar {
							existingNode.costSoFar = sn[i].parent.costSoFar + sn[i].cost
							existingNode.parent = sn[i].parent
						}
					} else {
						n := newNode(adjacentCoordinateX, adjacentCoordinateY, p.endX, p.endY, sn[i].parent, sn[i].cost)
						p.coordinateToNode[hashKey] = n
						item := &pqueue.Item{Value: n, Priority: int64(n.bestGuessDistance() * costPrecision)}
						p.openList.Push(item)
					}
				}
			}
		}
	}

	return nil, errors.New("path not found")
}

func (p *Pathfinder) getTileCost(x int, y int) float64 {
	pos := getHashKeyForPoint(x, y)
	cost, found := p.additionalPointsToCost[pos]
	if !found {
		cost = p.costMap[p.collisionGrid[y][x]]
	}
	return cost
}

func isTileWalkable(grid [][]int, acceptableTiles []int, x int, y int) bool {
	for i := 0; i < len(acceptableTiles); i++ {
		if grid[y][x] == acceptableTiles[i] {
			return true
		}
	}
	return false
}

func getHashKeyForPoint(x int, y int) string {
	return string(x) + "_" + string(y)
}

func getDistance(x1 int, y1 int, x2 int, y2 int) float64 {
	xx := float64(x2 - x1)
	yy := float64(y2 - y1)
	return math.Sqrt(xx*xx + yy*yy)
}
