package main

import "fmt"

type Gunner interface {
	Target() Point
	Hit(p Point)
	Miss(p Point)
}

type randomGunner struct {
	board *Board
}

func NewRandomGunner(board *Board) Gunner {
	return &randomGunner{
		board,
	}
}
func (g *randomGunner) Target() Point {
	return g.board.PickRandomPoint()
}
func (g *randomGunner) Hit(p Point) {
}
func (g *randomGunner) Miss(p Point) {
}

type linearGunner struct {
	board      *Board
	lastTarget Point
}

func NewLinearGunner(board *Board) Gunner {
	return &linearGunner{
		board,
		Point{0, -1},
	}
}
func (g *linearGunner) Target() Point {
	// scan row by row, left to right
	p := g.lastTarget.right()
	if g.board.IsOutOfBound(p) {
		p = Point{g.lastTarget.X + 1, 0}
	}
	g.lastTarget = p

	return p
}
func (g *linearGunner) Hit(p Point) {
}
func (g *linearGunner) Miss(p Point) {
}

// diagonalGunner take the strategry of targeting all points on a diagonal line
// based on a certain ship size, dividing the entire grid into small boxes and
// check points on the diagonal lines of each small box
// this could be better than linear scanning or picking completely random points
type diagonalGunner struct {
	board      *Board
	candidates []Point
	triedIndex int
}

func NewDiagonalGunner(board *Board, shipSizes []int) *diagonalGunner {
	g := &diagonalGunner{
		board,
		[]Point{},
		0,
	}
	g.markCandidates(shipSizes)

	return g
}

func (g *diagonalGunner) markCandidates(shipSizes []int) {
	// ASSUME the ship sizes are already sorted in descending order without duplicates
	for _, shipSize := range shipSizes {
		for row := 0; row < g.board.getSize(); row++ {
			repetition := g.board.getSize() / shipSize
			for i := 1; i <= repetition; i++ {
				col := shipSize*i - 1 - (row % shipSize)
				g.candidates = append(g.candidates, Point{row, col})
			}
		}
	}
}

func (g *diagonalGunner) print() {
	for _, v := range g.candidates {
		fmt.Println(v)
	}
}

func (g *diagonalGunner) Target() Point {
	// go through all candidates if any remains
	var p Point
	if g.triedIndex < len(g.candidates) {
		p = g.candidates[g.triedIndex]
		g.triedIndex++
		return p
	}

	return g.board.PickRandomPoint()
}
func (g *diagonalGunner) Hit(p Point) {
}
func (g *diagonalGunner) Miss(p Point) {
}

// ClusterGunner take the strategry of targeting all the adjacent points to the last hit
type ClusterGunner struct {
	board        *Board
	triedTargets Points
	candidates   []Point
}

func NewClusterGunner(board *Board) Gunner {
	return &ClusterGunner{
		board,
		Points{},
		[]Point{},
	}
}
func (cg *ClusterGunner) Target() Point {
	// target points on the candidate list
	for {
		n := len(cg.candidates)
		if n == 0 {
			break
		}

		p := cg.candidates[n-1]
		cg.candidates = cg.candidates[:n-1]
		if !cg.triedTargets.Contains(p) {
			return p
		}
	}

	// no cluster target, try a random one that hasn't tried before
	var p Point
	for {
		p = cg.board.PickRandomPoint()
		if !cg.triedTargets.Contains(p) {
			break
		}
	}
	return p
}

func (cg *ClusterGunner) Hit(p Point) {
	// record the hit
	cg.triedTargets = cg.triedTargets.Add(p)

	// record the cluster points next to this point
	cluster := []Point{
		p.left(),
		p.right(),
		p.up(),
		p.down(),
	}

	// filter out the outOfBound or already hit ones
	for _, v := range cluster {
		if !cg.board.IsOutOfBound(v) && !cg.triedTargets.Contains(v) {
			cg.candidates = append(cg.candidates, v)
		}
	}
}

func (cg *ClusterGunner) Miss(p Point) {
	cg.triedTargets = cg.triedTargets.Add(p)
}
