package main

import (
	"fmt"
	"math/rand"
	"time"
)

const EMPTY_SPACE = "."
const HIT = "x"

type Point struct {
	X, Y int
}

func (p Point) right() Point {
	return Point{p.X, p.Y + 1}
}
func (p Point) left() Point {
	return Point{p.X, p.Y - 1}
}
func (p Point) up() Point {
	return Point{p.X - 1, p.Y}
}
func (p Point) down() Point {
	return Point{p.X + 1, p.Y}
}

type Ship struct {
	name   string
	length int
}

type Board struct {
	size   int
	points [][]string
}

func (b *Board) getSize() int {
	return b.size
}

func (b *Board) Init(size int) {
	b.size = size
	b.points = make([][]string, size)
	for i := 0; i < size; i++ {
		b.points[i] = make([]string, size)
		for j := 0; j < size; j++ {
			b.points[i][j] = EMPTY_SPACE
		}
	}
}

func (b *Board) PlaceShipAt(p Point, name string) {
	if b.IsOutOfBound(p) {
		panic("point is out of bound")
	}
	b.points[p.X][p.Y] = name
}

func (b *Board) PickRandomPoint() Point {
	rand.Seed(time.Now().UnixNano())
	return Point{rand.Intn(b.size), rand.Intn(b.size)}
}

func (b Board) Print() {
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			fmt.Print(b.points[i][j])
		}
		fmt.Println()
	}
}

func (b Board) IsGameOver() bool {
	// game over when all left is empty space or sunking ships
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			p := Point{i, j}
			if !b.IsEmptySpace(p) && !b.IsSunkShip(p) {
				return false
			}
		}
	}
	return true
}

func (b Board) IsOutOfBound(p Point) bool {
	return p.X >= b.size || p.X < 0 || p.Y >= b.size || p.Y < 0
}

func (b Board) IsEmptySpace(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	return b.points[p.X][p.Y] == EMPTY_SPACE
}

func (b Board) IsSunkShip(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	return b.points[p.X][p.Y] == HIT
}

func (b Board) Hit(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	if b.IsEmptySpace(p) {
		return false
	}
	b.points[p.X][p.Y] = HIT
	return true
}

func (b Board) AreEmptySpaces(ps []Point) bool {
	for _, p := range ps {
		if !b.IsEmptySpace(p) {
			return false
		}
	}
	return true
}
