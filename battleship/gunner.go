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
		Point{0, 0},
	}
}
func (g *linearGunner) Target() Point {
	// scan row by row, left to right
	p := g.lastTarget.right()
	for {
		if !g.board.IsOutOfBound(p) {
			break
		}
		p = Point{g.lastTarget.X + 1, 0}
	}
	g.lastTarget = p
	return p
}
func (g *linearGunner) Hit(p Point) {
}
func (g *linearGunner) Miss(p Point) {
}
