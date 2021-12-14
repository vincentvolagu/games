package main

import (
	"math/rand"
	"time"
)

type HorizontalFormation struct {
	board *Board
	ships []Ship
}

func (f HorizontalFormation) Form() {
	for _, ship := range f.ships {
		p := f.findSpotForShip(ship.length)
		for i := 0; i < ship.length; i++ {
			f.board.points[p.X][p.Y+i] = ship.name
		}
	}
}

func (f HorizontalFormation) findSpotForShip(length int) Point {
	var p Point
	for p = f.getRandomPoint(length); !f.fitsHorizontally(p, length); p = f.getRandomPoint(length) {
	}
	return p
}

func (f HorizontalFormation) getRandomPoint(length int) Point {
	rand.Seed(time.Now().UnixNano())
	randX := rand.Intn(f.board.size)
	randY := rand.Intn(f.board.size - length)
	return Point{randX, randY}
}

func (f HorizontalFormation) fitsHorizontally(p Point, length int) bool {
	for i := 0; i < length; i++ {
		if !f.board.IsEmptySpace(p) {
			return false
		}
		p = p.right()
	}
	return true
}
