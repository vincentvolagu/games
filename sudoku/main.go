package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	count := 0
	for {
		grid := newGrid(5)
		grid.Fill()
		found := grid.Validate()

		if found {
			fmt.Println("found a valid sudoku after", count, "attempts")
			grid.Print()
			break
		}
		count++
	}
}

type Grid struct {
	size     int
	rowBased [][]int
	colBased [][]int
}

func newGrid(size int) Grid {
	return Grid{
		size,
		emptyGrid(size),
		emptyGrid(size),
	}
}

func (g Grid) Validate() bool {
	// do a simple version first, check no repeated number in any given row or column
	for i := 0; i < g.size; i++ {
		seen := make(map[int]bool)
		for j := 0; j < g.size; j++ {
			v := g.rowBased[i][j]
			if v <= 0 || v > g.size {
				return false
			}
			if _, exists := seen[v]; exists {
				return false
			}
			seen[v] = true
		}
	}
	for i := 0; i < g.size; i++ {
		seen := make(map[int]bool)
		for j := 0; j < g.size; j++ {
			v := g.rowBased[j][i]
			if v <= 0 || v > g.size {
				return false
			}
			if _, exists := seen[v]; exists {
				return false
			}
			seen[v] = true
		}
	}
	return true
}

func (g Grid) validRange() []int {
	nr := make([]int, g.size)
	for i := 0; i < g.size; i++ {
		nr[i] = i + 1
	}
	return nr
}

func (g Grid) randOrderedList() []int {
	nr := g.validRange()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nr), func(i, j int) { nr[i], nr[j] = nr[j], nr[i] })

	return nr
}

func (g Grid) Fill() {
	// start with a randomised row 0
	l := g.randOrderedList()
	for j := 0; j < g.size; j++ {
		g.rowBased[0][j] = l[j]
		g.colBased[j][0] = l[j]
	}

	// for the following rows, check column to find valid number to put
	for i := 1; i < g.size; i++ {
		l = g.randOrderedList()
		for j := 0; j < g.size; j++ {
			next := l[j]
			if !g.hasValue(next, g.colBased[j]) {
				g.rowBased[i][j] = next
				g.colBased[j][i] = next
			}
		}
	}
}

func (g Grid) hasValue(needle int, haystack []int) bool {
	for k := 0; k < g.size; k++ {
		if needle == haystack[k] {
			return true
		}
	}
	return false
}

func (g Grid) Print() {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			fmt.Print("|", g.rowBased[i][j])
		}
		fmt.Println("|")
	}
}

func initialiseGridFromInput() Grid {
	return Grid{}
}

func emptyGrid(size int) [][]int {
	g := make([][]int, size)
	for i := 0; i < size; i++ {
		g[i] = make([]int, size)
	}
	return g
}
