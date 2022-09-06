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
	size    int
	entries [][]int
}

func newGrid(size int) Grid {
	g := make([][]int, size)
	for i := 0; i < size; i++ {
		g[i] = make([]int, size)
	}
	return Grid{
		size,
		g,
	}
}

func (g Grid) Validate() bool {
	// do a simple version first, check no repeated number in any given row or column
	for i := 0; i < g.size; i++ {
		seen := make(map[int]bool)
		for j := 0; j < g.size; j++ {
			v := g.entries[i][j]
			if _, exists := seen[v]; exists {
				return false
			}
			seen[v] = true
		}
	}
	for i := 0; i < g.size; i++ {
		seen := make(map[int]bool)
		for j := 0; j < g.size; j++ {
			v := g.entries[j][i]
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
	for i := 0; i < g.size; i++ {
		l := g.randOrderedList()
		for j := 0; j < g.size; j++ {
			g.entries[i][j] = l[j]
		}
	}
}

func (g Grid) Print() {
	for i := 0; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			fmt.Print("|", g.entries[i][j])
		}
		fmt.Println("|")
	}
}

func initialiseGridFromInput() Grid {
	return Grid{}
}
