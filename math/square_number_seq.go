package main

import "fmt"
import "math"

func main() {
	n := 16
	// for n := 16; n < 40; n++ {
	// list := naturalOrderedList(n)
	// sn := squareNumbers(n)
	// fmt.Println("list is ", list)
	// fmt.Println("square numbers can be ", sn)
	graph := newSquareNumberGraph(n, true)
	graph.initVertexes()
	graph.initPaths()
	fmt.Println("total recorded path", len(graph.paths))
	// fullPaths := graph.findFullLengthPath()
	// for _, v := range fullPaths {
	// fmt.Println("found path", v)
	// }
	// if we plot those nodes on a x-y plane
	// then they all sit on diagnol lines where x+y=sq
	// eg. x+y=4, x+y=9 ...
	// TODO: figure out the formula f(n) = # of possible vertex
	//
	// normally with n vertex, we have have n(n-1)/2 ~ O(n^2) edges
	// TODO: figure out the formula to possible # of edges
	// totalNodes := 0
	// for _, v := range graph {
	// totalNodes = totalNodes + len(v)
	// }
	// fmt.Println(n, "numbers has total of", totalNodes, "nodes")
	// }

	// searchGraph(graph, n)
}

// navigate the graph until we can cover all points
func searchGraph(graph map[int][]squarePoint, n int) {
	for _, row := range graph {
		for _, col := range row {
			// this has some repetition as it discover the same paths repeately
			path := newPath()
			path.add(col.source)
			if depthFirstSearch(graph, col, path, n) {
				fmt.Println("found path", path.nodes)
				return
			}
		}
	}
}

func depthFirstSearch(graph map[int][]squarePoint, p squarePoint, currentPath *path, depth int) bool {
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
	nodes []int
}

func newPath() *path {
	return &path{}
}

func (p *path) clone() *path {
	dest := make([]int, len(p.nodes))
	copy(dest, p.nodes)
	return &path{
		nodes: dest,
	}
}

func (p *path) has(next int) bool {
	for _, v := range p.nodes {
		if next == v {
			return true
		}
	}
	return false
}

func (p *path) add(next int) bool {
	// avoid acyclic path
	if !p.has(next) {
		p.nodes = append(p.nodes, next)
		return true
	}
	return false
}

func (p *path) remove(prev int) {
	p.nodes = p.nodes[:len(p.nodes)-1]
}

func (p path) len() int {
	return len(p.nodes)
}

func (p path) asString() string {
	return fmt.Sprint(p.nodes)
}

type squarePoint struct {
	source, dest int
}

func (p squarePoint) connects(other squarePoint) bool {
	return p.source == other.dest && p.dest == other.source
}
func (p squarePoint) asKey() string {
	return fmt.Sprint(p.source, "_", p.dest)
}

type squareNumberGraph struct {
	n int
	// map of possible vertex at x
	// 1 -> [8, 3, ...]
	vertexMap map[int][]int
	// edges that already been traversed / linked
	paths []*path
	debug bool
}

func newSquareNumberGraph(n int, debug bool) squareNumberGraph {
	return squareNumberGraph{
		n:         n,
		vertexMap: make(map[int][]int),
		debug:     debug,
	}
}

// Runs in O(n^2)
// TODO: can this be improved ?
func (g *squareNumberGraph) initVertexes() {
	list := naturalOrderedList(g.n)
	for _, v := range list {
		g.vertexMap[v] = g.findSquarePoints(v, list)
	}
}
func (g *squareNumberGraph) findSquarePoints(n int, list []int) []int {
	var points []int
	for _, v := range list {
		if v == n {
			continue
		}
		if isSquareNumber(float64(n + v)) {
			points = append(points, v)
		}
	}
	return points
}

func (g *squareNumberGraph) findFullLengthPath() []*path {
	var fullPath []*path
	for _, path := range g.paths {
		if path.len() == g.n {
			fullPath = append(fullPath, path)
		}
	}
	return fullPath
}

func (g *squareNumberGraph) initPaths() {
	for k, list := range g.vertexMap {
		path := newPath()
		path.add(k)
		for _, v := range list {
			g.nextDepth(v, path.clone())
		}
	}
}

func (g *squareNumberGraph) nextDepth(source int, path *path) {
	// got cyclic path, the path ends there, record it
	if !path.add(source) {
		if g.debug {
			fmt.Println("found path end", path)
		}
		g.paths = append(g.paths, path)
		return
	}

	// fmt.Println("continue path search", path)
	for _, dest := range g.vertexMap[source] {
		if g.debug {
			fmt.Println("looking path", path, "next", dest)
		}
		g.nextDepth(dest, path.clone())
	}
}

// Runs in O(n)

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
