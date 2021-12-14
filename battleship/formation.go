package main

import (
	"math/rand"
	"time"
)

type EdgeFormation struct {
	board *Board
	ships []Ship
}

func (f EdgeFormation) PlaceShips() {
	rand.Seed(time.Now().UnixNano())
	ps := f.initEdgePoints()
	var shipPoints []Point
	for _, ship := range f.ships {
		for {
			randEdge := rand.Intn(4)
			randIndex := rand.Intn(f.board.size - ship.length)
			shipPoints = ps[randEdge][randIndex : randIndex+ship.length]
			if f.board.AreEmptySpaces(shipPoints) {
				break
			}
		}
		for _, p := range shipPoints {
			f.board.points[p.X][p.Y] = ship.name
		}
	}
}

func (f EdgeFormation) initEdgePoints() [][]Point {
	ps := make([][]Point, 4)
	ps[0] = make([]Point, f.board.size)
	ps[1] = make([]Point, f.board.size)
	ps[2] = make([]Point, f.board.size)
	ps[3] = make([]Point, f.board.size)
	// top edge
	for i := 0; i < f.board.size; i++ {
		ps[0][i] = Point{0, i}
	}
	// bottom edge
	for i := 0; i < f.board.size; i++ {
		ps[1][i] = Point{f.board.size - 1, i}
	}
	// left edge
	for i := 0; i < f.board.size; i++ {
		ps[2][i] = Point{i, 0}
	}
	// right edge
	for i := 0; i < f.board.size; i++ {
		ps[3][i] = Point{i, f.board.size - 1}
	}
	return ps
}

func makeEdgeFormation(board *Board, ships []Ship) EdgeFormation {
	return EdgeFormation{board, ships}
}

type Formation struct {
	board     *Board
	ships     []Ship
	direction DirectionStrategy
}

func makeHorizontalFormation(board *Board, ships []Ship) Formation {
	return Formation{
		board,
		ships,
		horizontal{},
	}
}
func makeVerticalFormation(board *Board, ships []Ship) Formation {
	return Formation{
		board,
		ships,
		vertical{},
	}
}
func makeRandomFormation(board *Board, ships []Ship) Formation {
	return Formation{
		board,
		ships,
		&randomDirection{},
	}
}

func (f Formation) PlaceShips() {
	for _, ship := range f.ships {
		f.direction.Reset()
		ps := f.findSpotForShip(ship.length)
		for _, p := range ps {
			f.board.points[p.X][p.Y] = ship.name
		}
	}
}

func (f Formation) findSpotForShip(length int) []Point {
	var ps []Point
	var err error
	p := f.getRandomPoint(length)
	for ps, err = f.placeShipAtPoint(p, length); err != nil; {
		p = f.getRandomPoint(length)
		ps, err = f.placeShipAtPoint(p, length)
	}
	return ps
}

func (f Formation) getRandomPoint(length int) Point {
	rand.Seed(time.Now().UnixNano())
	randX := rand.Intn(f.board.size)
	randY := rand.Intn(f.board.size)
	return Point{randX, randY}
}

func (f Formation) placeShipAtPoint(p Point, length int) ([]Point, error) {
	ps := make([]Point, length)
	for i := 0; i < length; i++ {
		if !f.board.IsEmptySpace(p) {
			return []Point{}, NoSpaceForShipError{}
		}
		ps[i] = p
		p = f.direction.NextPoint(p)
	}
	return ps, nil
}

type NoSpaceForShipError struct{}

func (e NoSpaceForShipError) Error() string {
	return "no space for ship"
}

type DirectionStrategy interface {
	Reset()
	NextPoint(p Point) Point
}

type horizontal struct {
}

func (f horizontal) NextPoint(p Point) Point {
	return p.right()
}
func (f horizontal) Reset() {
}

type vertical struct {
}

func (f vertical) NextPoint(p Point) Point {
	return p.down()
}
func (f vertical) Reset() {
}

type randomDirection struct {
	ds DirectionStrategy
}

func (f *randomDirection) Reset() {
	randNum := rand.Intn(100) % 2
	if randNum == 0 {
		f.ds = horizontal{}
	} else {
		f.ds = vertical{}
	}
}
func (f *randomDirection) NextPoint(p Point) Point {
	if f.ds == nil {
		f.Reset()
	}
	return f.ds.NextPoint(p)
}
