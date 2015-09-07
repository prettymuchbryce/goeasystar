package goeasystar

type node struct {
	parent                 *node
	x                      int
	y                      int
	costSoFar              float64
	simpleDistanceToTarget float64
	list                   int
}

// Point represents a point in the completed path
type Point struct {
	X int
	Y int
}

type searchNode struct {
	x      int
	y      int
	endX   int
	endY   int
	cost   float64
	parent *node
}

type searchNodes []*searchNode

func (slice searchNodes) Len() int {
	return len(slice)
}

func (slice searchNodes) Less(i, j int) bool {
	iCost := slice[i].cost + getDistance(slice[i].parent.x+slice[i].x,
		slice[i].parent.y+slice[i].y, slice[i].endX, slice[i].endY)*costPrecision
	jCost := slice[j].cost + getDistance(slice[j].parent.x+slice[j].x,
		slice[j].parent.y+slice[j].y, slice[j].endX, slice[j].endY)*costPrecision

	if iCost > jCost {
		return true
	}

	return false
}

func (slice searchNodes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newSearchNode(x int, y int, endX int, endY int, cost float64, parent *node) *searchNode {
	sn := &searchNode{}

	sn.x = x
	sn.y = y
	sn.endX = endX
	sn.endY = endY
	sn.cost = cost
	sn.parent = parent

	return sn
}

func newNode(x int, y int, endX int, endY int, parent *node, cost float64) *node {
	n := &node{}

	simpleDistanceToTarget := getDistance(x, y, endX, endY)

	costSoFar := float64(0)
	if parent != nil {
		costSoFar = parent.costSoFar + cost
	} else {
		costSoFar = simpleDistanceToTarget
	}

	n.x = x
	n.y = y
	n.parent = parent
	n.costSoFar = costSoFar
	n.simpleDistanceToTarget = simpleDistanceToTarget

	return n
}

func (n *node) bestGuessDistance() float64 {
	return n.costSoFar + n.simpleDistanceToTarget
}
