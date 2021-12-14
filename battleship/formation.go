package main

import (
	"math/rand"
	"time"
)

type EdgeFormation struct {
	board     *Board
	ships     []Ship
	direction DirectionStrategy
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
