package main

import "fmt"

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

func (b *Board) placeShip(x, y int, ship Ship) {

}

func (b Board) Print() {
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			fmt.Print(b.points[i][j])
		}
		fmt.Println()
	}
}

func (b Board) IsEmptySpace(p Point) bool {
	if p.X >= b.size || p.X < 0 || p.Y >= b.size || p.Y < 0 {
		return false
	}
	return b.points[p.X][p.Y] == EMPTY_SPACE
}

func (b Board) AreEmptySpaces(ps []Point) bool {
	for _, p := range ps {
		if !b.IsEmptySpace(p) {
			return false
		}
	}
	return true
}
