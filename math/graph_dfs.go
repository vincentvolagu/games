package main

import "fmt"
import "math"
import "sort"

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

func (p *path) head() int {
	return p.nodes[0]
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
	if p.has(next) {
		return false
	}
	p.nodes = append(p.nodes, next)
	return true
}

func (p *path) remove(prev int) {
	p.nodes = p.nodes[:len(p.nodes)-1]
}

func (p path) len() int {
	return len(p.nodes)
}

func (p path) runningSum() string {
	sum := make([]int, p.len()-1)
	for i := 1; i < p.len(); i++ {
		sum[i-1] = p.nodes[i] + p.nodes[i-1]
	}
	return fmt.Sprint(sum)
}

func (p path) asString() string {
	return fmt.Sprint(p.nodes)
}

func (p path) asInts() []int {
	return p.nodes
}

type squareNumberGraph struct {
	n int

	// map of adjacent edge from vertex x
	// 1 -> [8, 3, ...]
	vertexMap map[int][]int

	// edges that already been traversed / linked
	paths map[string]*path

	debug         bool
	nodeTraversed int
}

func NewSquareNumberGraph(n int, debug bool) squareNumberGraph {
	return squareNumberGraph{
		n:         n,
		vertexMap: make(map[int][]int),
		paths:     make(map[string]*path),
		debug:     debug,
	}
}

func (g *squareNumberGraph) SquareSumSort(n int) [][]int {
	g.initVertexes()
	g.initPaths()

	if g.debug {
		fmt.Println("n =", n)
		fmt.Println("total vertex", g.vertexCount())
		fmt.Println("total recorded path", len(g.paths))
		fmt.Println("total node traversed", g.nodeTraversed)
		fullPaths, _ := g.findFullLengthPath()
		fmt.Println("total full path found", len(fullPaths))
		sort.Strings(fullPaths)
		for _, v := range fullPaths {
			fmt.Println("found path", v)
		}
		// for _, v := range runningSums {
		// fmt.Println("running sum", v)
		// }
	}
	return g.getFullLengthPathAsInts()
}

// Runs in O(n^2)
func (g *squareNumberGraph) initVertexes() {
	list := make([]int, g.n)
	for i := 1; i <= g.n; i++ {
		list[i-1] = i
	}
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
		if g.isSquareNumber(float64(n + v)) {
			points = append(points, v)
		}
	}
	return points
}
func (g *squareNumberGraph) vertexCount() int {
	sum := 0
	for _, v := range g.vertexMap {
		sum = sum + len(v)
	}
	return sum
}

func (g *squareNumberGraph) getFullLengthPathAsInts() [][]int {
	var fullPath [][]int
	for _, path := range g.paths {
		if path.len() == g.n {
			fullPath = append(fullPath, path.asInts())
		}
	}
	return fullPath
}

func (g *squareNumberGraph) findFullLengthPath() ([]string, []string) {
	var fullPath []string
	var runningSum []string
	for _, path := range g.paths {
		if path.len() == g.n {
			fullPath = append(fullPath, path.asString())
			runningSum = append(runningSum, path.runningSum())
		}
	}
	return fullPath, runningSum
}

// This seems to be in some from of O(n!), which is really bad
// TODO: find a way to improve this
// - just throw goroutines at it to increase parallelism
func (g *squareNumberGraph) initPaths() {
	ch := make(chan *path, 100)
	done := make(chan int, len(g.vertexMap))
	for k, list := range g.vertexMap {
		path := newPath()
		path.add(k)
		go g.findPathsFrom(list, path, ch, done)
	}

	doneCount := 0
	for {
		select {
		case v := <-ch:
			if v.len() == g.n {
				g.paths[v.asString()] = v
			}

		case <-done:
			doneCount++
			if doneCount == g.n {
				return
			}
		}
	}
}

func (g *squareNumberGraph) findPathsFrom(adjacentPoints []int, currentPath *path, ch chan *path, done chan int) {
	for _, v := range adjacentPoints {
		g.nextDepth(v, currentPath.clone(), ch)
	}
	done <- 1
}

func (g *squareNumberGraph) nextDepth(source int, currentPath *path, ch chan *path) {
	g.nodeTraversed++

	// got cyclic path, the path ends there, record it
	if !currentPath.add(source) {
		// fmt.Println("found path end", currentPath)
		// g.paths[currentPath.asString()] = currentPath
		ch <- currentPath
		return
	}

	for _, dest := range g.vertexMap[source] {
		if g.debug {
			fmt.Println("looking path", currentPath, "next", dest)
		}
		g.nextDepth(dest, currentPath.clone(), ch)
	}
}

func (g *squareNumberGraph) isSquareNumber(n float64) bool {
	sqrt := math.Sqrt(n)
	return math.Floor(sqrt) == sqrt
}
