package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

func main() {
	// grid := newGrid(4)
	// grid = grid.Fill()
	// fmt.Println("end result")
	// grid.Print()
	count := 0
	for {
		grid := newGrid(16)
		grid = grid.Fill()
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

type subGrid struct {
	mainGrid Grid

	startRow, startCol, endRow, endCol int
}

func newGrid(size int) Grid {
	return Grid{
		size,
		emptyGrid(size),
		emptyGrid(size),
	}
}

// the row/column is valid only it contains all the numbers from 1...n only once
func (g Grid) isValidLine(line []int) bool {
	expected := g.validRange()
	actual := make([]int, g.size)
	copy(actual, line)
	sort.Ints(actual)
	return reflect.DeepEqual(expected, actual)
}

func (g Grid) Validate() bool {
	// do a simple version first, check no repeated number in any given row or column
	for i := 0; i < g.size; i++ {
		if !g.isValidLine(g.rowBased[i]) && !g.isValidLine(g.colBased[i]) {
			return false
		}
	}

	// find smaller sub-grid and check numbers in that grid
	if subGridSize, hasSubGrid := g.subGridSize(); hasSubGrid {
		for i := 0; i < subGridSize; i++ {
			for j := 0; j < subGridSize; j++ {
				var sgList []int
				for k := 0; k < g.size; k++ {
					row := (i * subGridSize) + (k / subGridSize)
					col := (j * subGridSize) + (k % subGridSize)
					sgList = append(sgList, g.rowBased[row][col])
				}
				fmt.Println("checking list", sgList)
				if !g.isValidLine(sgList) {
					return false
				}
			}
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

func (g Grid) Fill() Grid {
	// start with a randomised row 0
	startingRow := g.randOrderedList()
	for j := 0; j < g.size; j++ {
		g.rowBased[0][j] = startingRow[j]
		g.colBased[j][0] = startingRow[j]
	}

	// shift the 1st randomised row to populate the rest rows
	for i := 1; i < g.size; i++ {
		for j := 0; j < g.size; j++ {
			g.rowBased[i][j] = startingRow[(i+j)%g.size]
			g.colBased[j][i] = startingRow[(i+j)%g.size]
		}
	}

	fmt.Println("after initial fill")
	g.Print()

	// swap rows so sub grids also match sequence
	if subGridSize, hasSubGrid := g.subGridSize(); hasSubGrid {
		clone := newGrid(g.size)
		ci := 0
		for i := 0; i < subGridSize; i++ {
			for j := 0; j < subGridSize; j++ {
				copy(clone.rowBased[ci], g.rowBased[i+j*subGridSize])
				ci++
			}
		}
		fmt.Println("after row swapping")
		clone.Print()
		return clone
	}

	return g
}

func (g Grid) subGridSize() (int, bool) {
	sqrt := math.Sqrt(float64(g.size))
	if math.Floor(sqrt) == sqrt {
		return int(sqrt), true
	}
	return 0, false
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
			fmt.Print("|", g.rowBased[i][j], "\t")
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
