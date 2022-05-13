package main

import "fmt"
import "math"

// Problem: for 1, 2, 3...n integer sequence, find a particular ordering that the sum of any 2 adjacent numbers are square number
//
// Thought: if we take a graph approach, this is to find the longest path of the graph that connects all vertexes
// TODO: plot the tuple on a x-y plane shows special navigation pattern, that may help improve performance

type SquareSumSorter interface {
	SquareSumSort(n int) [][]int
}

func main() {
	n := 16
	graph := NewSquareNumberGraph(n, false)
	squareInts := graph.SquareSumSort(n)
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

// Runs in O(n)
func naturalOrderedList(n int) []int {
	list := make([]int, n)
	for i := 1; i <= n; i++ {
		list[i-1] = i
	}
	return list
}

func getSquareNumbersUpTo(n int) []int {
	maxSum := n + n - 1
	var squareNums []int
	for i := 1; i <= maxSum; i++ {
		if isSquareNumber(float64(i)) {
			squareNums = append(squareNums, i)
		}
	}
	return squareNums
}

func bruteForce(n int) {
	list := make([]int, n)
	for i := 1; i <= n; i++ {
		list[i-1] = i
	}
	combo := factorialCombination(list)
	for _, row := range combo {
		// fmt.Println(row)
		if verify(row) {
			fmt.Println("found list that are sequence of adjacent square numbers", row)
			break
		}
	}
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}
	return factorial(n-1) * n
}

func factorialCombination(list []int) [][]int {
	if len(list) == 1 {
		return [][]int{list}
	}
	// newCombo := make([][]int, factorial(len(list)))
	var newCombo [][]int
	for i, v := range list {
		head := v
		// cheap way to remove a value from slice, changing the ordering though
		tail := make([]int, len(list))
		copy(tail, list)
		tail[i] = tail[len(list)-1]
		tail = tail[:len(tail)-1]
		tailCombo := factorialCombination(tail)
		// fmt.Println("create combo for ", head, " by ", tail, " with tail ", tailCombo)
		for _, combo := range tailCombo {
			newCombo = append(newCombo, append([]int{head}, combo...))
			// fmt.Println("new combo ", newCombo[i])
		}
	}

	return newCombo
}

func verify(orderedList []int) bool {
	for i, v := range orderedList {
		if !isSquareNumber((float64)(v + orderedList[i+1])) {
			return false
		}
	}
	return true
}

func isSquareNumber(n float64) bool {
	sqrt := math.Sqrt(n)
	return math.Floor(sqrt) == sqrt
}
