package main

import "fmt"

// Problem: for 1, 2, 3...n integer sequence, find a particular ordering that the sum of any 2 adjacent numbers are square number
//
// Thought: if we take a graph approach, this is to find the longest path of the graph that connects all vertexes
// maybe it's a NP-hard problem https://en.wikipedia.org/wiki/Longest_path_problem
// TODO: plot the tuple on a x-y plane shows special navigation pattern, that may help improve performance

type SquareSumSorter interface {
	SquareSumSort(n int) [][]int
}

func main() {
	n := 32
	// graph := NewSquareNumberGraph(n, false)
	// squareInts := graph.SquareSumSort(n)
	r := NewSquareNumberGraph(n, false)
	squareInts := r.SquareSumSort(n)
	for _, v := range squareInts {
		fmt.Println(v)
	}

	// if we plot those nodes on a x-y plane
	// then they all sit on diagnol lines where x+y=sq
	// eg. x+y=4, x+y=9 ...
	// TODO: figure out the formula f(n) = # of possible vertex
	//
	// normally with n vertex, we have have n(n-1)/2 ~ O(n^2) edges
	// TODO: figure out the formula to possible # of edges
}
