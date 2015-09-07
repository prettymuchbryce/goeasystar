# goeasystar

A port of [EasyStar.js](http://www.easystarjs.com) to golang.

##Installation
* `go get github.com/prettymuchbryce/goeasystar`

## Description

goeasystar is an A* pathfinding API written in Go.

## API

#### Main Methods

`easystar := goeasystar.NewPathfinder()``

`easystar.SetGrid(grid [][]int)`

`easystar.SetAcceptableTiles(t []int);`

`easystar.FindPath(startX, startY, endX, endY) ([]*goeasystar.Point, error);`

#### Additional Features

* Avoiding additional points (outside of acceptable tiles)

* Diagonals

* Corner cutting

* Setting costs per-tile

* Setting additional point costs (outside of tile costs)

## Usage

First create an instance of the Pathfinder.

	easystar := NewPathfinder()

Create a grid, or tilemap. You may have made this with a level editor, or procedurally. Let's keep it simple for this example.

    var grid [][]int
    grid = append(grid, []int{1, 0, 0, 0, 0})
    grid = append(grid, []int{0, 1, 0, 0, 0})
    grid = append(grid, []int{0, 0, 1, 0, 0})
    grid = append(grid, []int{0, 0, 0, 1, 0})
    grid = append(grid, []int{0, 0, 0, 0, 1})

Set our grid.

	easystar.SetGrid(grid)

Set tiles which are "walkable".

	easystar.SetAcceptableTiles([]int{1})

Find a path.

	path, err := easystar.FindPath(0, 0, 4, 4)
    fmt.Println(err == nil) // true

Oops. We didn't have diagonals enabled so there is no valid path. Lets try again.

    easystar.EnableDiagonals()
    path, err := easystar.FindPath(0, 0, 4, 4)
    fmt.Println(err) // nil
    len(path) // 4

## License

goeasystar is licensed under the MIT license.

goeasystar uses facebookgo's (pqueue)[http://www.github.com/facebookgo/pqueue] which falls under the Apache 2 license.

I would be happy to eventually remove it in favor of an MIT-licensed priority queue implementation.

## Support

If you have any questions, comments, or suggestions please open an issue.
