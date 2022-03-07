package main

type Gunner interface {
	Target() Point
	Hit(p Point)
	Miss(p Point)
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
