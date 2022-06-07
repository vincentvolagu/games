package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// const EMPTY_SPACE = "."
// const HIT = "#"
// const MISS = "x"
const EMPTY_SPACE = "🌊"
const MISS = "👻"
const HIT = "😁"

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

type Points []Point

func (ps Points) Contains(p Point) bool {
	for _, v := range ps {
		if v == p {
			return true
		}
	}
	return false
}

func (ps Points) Add(p Point) Points {
	return append(ps, p)
}

type Ship struct {
	name   string
	length int
}

type Board struct {
	size   int
	points [][]string

	usesEmoji bool
}

func NewBoard(size int, usesEmoji bool) *Board {
	b := &Board{
		size:      size,
		usesEmoji: usesEmoji,
	}
	b.init(size)
	return b
}

func (b *Board) getSize() int {
	return b.size
}
func (b *Board) printEmptySpace() string {
	if b.usesEmoji {
		return "🌊"
	}
	return "."
}
func (b *Board) printHit() string {
	if b.usesEmoji {
		return "😁"
	}
	return "#"
}
func (b *Board) printMiss() string {
	if b.usesEmoji {
		return "👻"
	}
	return "x"
}
func (b *Board) printShip(name string) string {
	if b.usesEmoji {
		return "🚢"
	}
	return name
}
func (b *Board) init(size int) {
	b.points = make([][]string, size)
	for i := 0; i < size; i++ {
		b.points[i] = make([]string, size)
		for j := 0; j < size; j++ {
			b.points[i][j] = b.printEmptySpace()
		}
	}
}
func (b *Board) PlaceShipAt(p Point, name string) {
	if b.IsOutOfBound(p) {
		panic("point is out of bound")
	}
	b.points[p.X][p.Y] = b.printShip(name)
}

func (b *Board) PickRandomPoint() Point {
	rand.Seed(time.Now().UnixNano())
	return Point{rand.Intn(b.size), rand.Intn(b.size)}
}

func (b Board) columnHeading() []string {
	if b.usesEmoji {
		return []string{" ", "A ", "B ", "C ", "D ", "E ", "F ", "G ", "H ", "I ", "J "}
	}
	return []string{" ", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
}
func (b Board) Print() {
	for _, v := range b.columnHeading() {
		fmt.Print(v)
	}
	fmt.Println()
	for i := 0; i < b.size; i++ {
		fmt.Print(i)
		for j := 0; j < b.size; j++ {
			fmt.Print(b.points[i][j])
		}
		fmt.Println()
	}
}
func (b Board) PrintForHumanPlayer() {
	for _, v := range b.columnHeading() {
		fmt.Print(v)
	}
	fmt.Println()
	for i := 0; i < b.size; i++ {
		fmt.Print(i)
		for j := 0; j < b.size; j++ {
			if b.IsFloatingShip(Point{i, j}) {
				fmt.Print(b.printEmptySpace())
			} else {
				fmt.Print(b.points[i][j])
			}
		}
		fmt.Println()
	}
}
func (b Board) TransformHumanPointInput(row int, col string) Point {
	rowIndex := row
	colIndex := -1
	for i, v := range b.columnHeading() {
		if v == strings.ToUpper(col) {
			colIndex = i - 1 // there's an empty left padding
		}
	}

	return Point{rowIndex, colIndex}
}

func (b Board) IsGameOver() bool {
	// game over when all left is empty space or sunking ships
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			p := Point{i, j}
			if b.IsFloatingShip(p) {
				return false
			}
		}
	}
	return true
}

func (b Board) IsOutOfBound(p Point) bool {
	return p.X >= b.size || p.X < 0 || p.Y >= b.size || p.Y < 0
}
func (b Board) IsFloatingShip(p Point) bool {
	return !b.IsEmptySpace(p) && !b.IsSunkShip(p) && !b.IsMiss(p)
}
func (b Board) IsEmptySpace(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	return b.points[p.X][p.Y] == b.printEmptySpace()
}
func (b Board) IsSunkShip(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	return b.points[p.X][p.Y] == b.printHit()
}
func (b Board) IsMiss(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	return b.points[p.X][p.Y] == b.printMiss()
}
func (b Board) Hit(p Point) bool {
	if b.IsOutOfBound(p) {
		return false
	}
	if b.IsFloatingShip(p) {
		b.points[p.X][p.Y] = b.printHit()
		return true
	}
	b.points[p.X][p.Y] = b.printMiss()
	return false
}
func (b Board) AreEmptySpaces(ps []Point) bool {
	for _, p := range ps {
		if !b.IsEmptySpace(p) {
			return false
		}
	}
	return true
}
