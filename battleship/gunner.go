package main

type Gunner interface {
	Target() Point
	Hit(p Point)
	Miss(p Point)
}

type GunnerFactory interface {
	MakeGunner(board *Board, shipSizes []int) Gunner
}

type LuckyGunnerFactory struct{}

func (f LuckyGunnerFactory) MakeGunner(board *Board, shipSizes []int) Gunner {
	return NewLuckyGunner(board)
}

func NewLuckyGunner(board *Board) Gunner {
	return &luckyGunner{
		board,
		Points{},
	}
}

// luckyGunner hits random targets without any memory of hit/miss history
type luckyGunner struct {
	board        *Board
	triedTargets Points
}

func (g *luckyGunner) Target() Point {
	var p Point
	for {
		p = g.board.PickRandomPoint()
		if !g.triedTargets.Contains(p) {
			break
		}
	}
	return p
}
func (g *luckyGunner) Hit(p Point) {
	g.triedTargets = g.triedTargets.Add(p)
}
func (g *luckyGunner) Miss(p Point) {
	g.triedTargets = g.triedTargets.Add(p)
}

type LinearGunnerFactory struct{}

func (f LinearGunnerFactory) MakeGunner(board *Board, shipSizes []int) Gunner {
	return NewLinearGunner(board)
}

// linearGunner hits targets by linear scanning (left to right, top to bottom)
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

// ClusterGunner take the strategry of targeting all the adjacent points to the last hit
type ClusterGunner struct {
	board             *Board
	triedTargets      Points
	candidates        []Point
	defaultCandidates []Point
	triedIndex        int
}

type ClusterGunnerFactory struct{}

func (f ClusterGunnerFactory) MakeGunner(board *Board, shipSizes []int) Gunner {
	return NewClusterGunner(board, shipSizes)
}

func NewClusterGunner(board *Board, shipSizes []int) Gunner {
	g := &ClusterGunner{
		board,
		Points{},
		[]Point{},
		[]Point{},
		0,
	}
	g.markCandidates(shipSizes)
	return g
}

// this take the strategry of targeting all points on a diagonal line
// based on a certain ship size, dividing the entire grid into small boxes and
// check points on the diagonal lines of each small box
// this could be better than linear scanning or picking completely random points
func (g *ClusterGunner) markCandidates(shipSizes []int) {
	// ASSUME the ship sizes are already sorted in descending order without duplicates
	for _, shipSize := range shipSizes {
		for row := 0; row < g.board.getSize(); row++ {
			repetition := g.board.getSize() / shipSize
			for i := 1; i <= repetition; i++ {
				col := shipSize*i - 1 - (row % shipSize)
				g.defaultCandidates = append(g.defaultCandidates, Point{row, col})
			}
		}
	}
}
func (g *ClusterGunner) Target() Point {
	// target points on the candidate list
	for {
		n := len(g.candidates)
		if n == 0 {
			break
		}

		p := g.candidates[n-1]
		g.candidates = g.candidates[:n-1]
		if !g.triedTargets.Contains(p) {
			return p
		}
	}

	// no cluster target, try a defaut strategy

	// go through all candidates if any remains
	var p Point
	for {
		if g.triedIndex >= len(g.defaultCandidates) {
			break
		}
		p = g.defaultCandidates[g.triedIndex]
		g.triedIndex++
		if !g.triedTargets.Contains(p) {
			return p
		}
	}

	for {
		p = g.board.PickRandomPoint()
		if !g.triedTargets.Contains(p) {
			break
		}
	}

	return p
}

func (g *ClusterGunner) Hit(p Point) {
	// record the hit
	g.triedTargets = g.triedTargets.Add(p)

	// record the cluster points next to this point
	cluster := []Point{
		p.left(),
		p.right(),
		p.up(),
		p.down(),
	}

	// filter out the outOfBound or already hit ones
	for _, v := range cluster {
		if !g.board.IsOutOfBound(v) && !g.triedTargets.Contains(v) {
			g.candidates = append(g.candidates, v)
		}
	}
}

func (g *ClusterGunner) Miss(p Point) {
	g.triedTargets = g.triedTargets.Add(p)
}
