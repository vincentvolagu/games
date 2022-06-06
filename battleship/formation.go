package main

import (
	"math/rand"
	"time"
)

type ShipCoordinator interface {
	PlaceShips(board *Board, ships []Ship)
}

type RandomCoordinator struct {
	coordinators []ShipCoordinator
}

func MakeRandomCoordinator() ShipCoordinator {
	return RandomCoordinator{
		[]ShipCoordinator{
			LuckyDraw{&randomLine{}},
			EdgeLover{},
		},
	}
}

func (c RandomCoordinator) PlaceShips(board *Board, ships []Ship) {
	rand.Seed(time.Now().UnixNano())
	randPick := rand.Intn(len(c.coordinators))
	c.coordinators[randPick].PlaceShips(board, ships)
}

type EdgeLover struct {
}

func (f EdgeLover) PlaceShips(board *Board, ships []Ship) {
	// there is a potential infinite loop here if we place the longest ship last
	// because by then all 4 edges might have already been occupied by other shorter
	// ships so none are left to fit the longest ship
	// In order to avoid that:
	// 1. mandate or sort the input ships by length and place the longer ones first
	// 2. fallback to 2nd out most edges when the out most ones are all used up
	// TODO: implement option 2 given that's more scalable
	ps := f.initEdgePoints(board.getSize())
	var shipPoints []Point
	for _, ship := range ships {
		for {
			rand.Seed(time.Now().UnixNano())
			randEdge := rand.Intn(4)
			randIndex := rand.Intn(board.getSize() - ship.length)
			shipPoints = ps[randEdge][randIndex : randIndex+ship.length]
			if board.AreEmptySpaces(shipPoints) {
				break
			}
		}
		for _, p := range shipPoints {
			board.PlaceShipAt(p, ship.name)
		}
	}
}

func (f EdgeLover) initEdgePoints(boardSize int) [][]Point {
	ps := make([][]Point, 4)
	ps[0] = make([]Point, boardSize)
	ps[1] = make([]Point, boardSize)
	ps[2] = make([]Point, boardSize)
	ps[3] = make([]Point, boardSize)
	// top edge
	for i := 0; i < boardSize; i++ {
		ps[0][i] = Point{0, i}
	}
	// bottom edge
	for i := 0; i < boardSize; i++ {
		ps[1][i] = Point{boardSize - 1, i}
	}
	// left edge
	for i := 0; i < boardSize; i++ {
		ps[2][i] = Point{i, 0}
	}
	// right edge
	for i := 0; i < boardSize; i++ {
		ps[3][i] = Point{i, boardSize - 1}
	}
	return ps
}

type LuckyDraw struct {
	ln LineNavigator
}

func MakeLuckyDraw() LuckyDraw {
	return LuckyDraw{&randomLine{}}
}

func (f LuckyDraw) PlaceShips(board *Board, ships []Ship) {
	var randPoint Point
	var ps []Point
	var err error
	for _, ship := range ships {
		f.ln.Reset()
		for {
			randPoint = board.PickRandomPoint()
			ps, err = f.placeShipAtPoint(board, randPoint, ship.length)
			if err == nil {
				break
			}
		}
		for _, p := range ps {
			board.PlaceShipAt(p, ship.name)
		}
	}
}

func (f LuckyDraw) placeShipAtPoint(board *Board, p Point, length int) ([]Point, error) {
	ps := make([]Point, length)
	for i := 0; i < length; i++ {
		if !board.IsEmptySpace(p) {
			return []Point{}, NoSpaceForShipError{}
		}
		ps[i] = p
		p = f.ln.NextPoint(p)
	}
	return ps, nil
}

type NoSpaceForShipError struct{}

func (e NoSpaceForShipError) Error() string {
	return "no space for ship"
}
