package main

import "fmt"
import "math"

func main() {
	fmt.Println("vim-go")
	n := 16
	list := naturalOrderedList(n)
	sn := squareNumbers(n)
	fmt.Println("list is ", list)
	fmt.Println("square numbers can be ", sn)
	pairMap := make(map[int][]squarePair)
	for _, v := range list {
		pairMap[v] = findSquarePairs(v, list)
	}
	fmt.Println("map: ")
	for _, v := range pairMap {
		fmt.Println(v)
	}

	// navigate the graph until we can cover all points

	for _, row := range pairMap {
		for _, col := range row {
			path := newPath()
			path.add(col.source)
			if depthFirstSearch(pairMap, col, path, n) {
				fmt.Println("found path", path.nodes)
				break
			}
		}
	}
}

func depthFirstSearch(graph map[int][]squarePair, p squarePair, currentPath *path, depth int) bool {
	edges, _ := graph[p.dest]
	if !currentPath.add(p.dest) {
		return false
	}
	if currentPath.len() == depth {
		return true
	}
	// fmt.Println("looking at ", p)
	// fmt.Println("looking ahead, path = ", currentPath.asString())
	for _, v := range edges {
		// fmt.Println("checking edge", v)
		if v.connects(p) {
			continue
		}
		// depth first search, when fail, backtrack and try next node
		if depthFirstSearch(graph, v, currentPath, depth) {
			return true
		}
	}
	// fmt.Println("no path there, removing ", p.dest)
	currentPath.remove(p.dest)
	return false
}

type path struct {
	nodes   []int
	visited map[int]int
}

func newPath() *path {
	return &path{visited: make(map[int]int)}
}

func (p *path) add(next int) bool {
	// avoid acyclic path
	if _, ok := p.visited[next]; !ok {
		p.nodes = append(p.nodes, next)
		p.visited[next] = len(p.nodes) - 1
		return true
	}
	return false
}

func (p *path) remove(prev int) {
	p.nodes = p.nodes[:len(p.nodes)-1]
	delete(p.visited, prev)
}

func (p path) len() int {
	return len(p.nodes)
}

func (p path) asString() string {
	return fmt.Sprint(p.nodes)
}

type squarePair struct {
	source, dest int
}

func (p squarePair) connects(other squarePair) bool {
	return p.source == other.dest && p.dest == other.source
}
func (p squarePair) asKey() string {
	return fmt.Sprint(p.source, "_", p.dest)
}

func findSquarePairs(n int, list []int) []squarePair {
	var pairs []squarePair
	for _, v := range list {
		if v == n {
			continue
		}
		if isSquareNumber(float64(n + v)) {
			pairs = append(pairs, squarePair{n, v})
		}
	}
	return pairs
}

func naturalOrderedList(n int) []int {
	list := make([]int, n)
	for i := 1; i <= n; i++ {
		list[i-1] = i
	}
	return list
}

func squareNumbers(n int) []int {
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
