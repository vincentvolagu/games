package main

import (
	"math/rand"
	"time"
)

type LineNavigator interface {
	NextPoint(p Point) Point
	Alternate() LineNavigator
	Reverse() LineNavigator
	IsHorizontal() bool
}

type goLeft struct {
}

func (nav goLeft) NextPoint(p Point) Point {
	return p.left()
}
func (nav goLeft) Alternate() LineNavigator {
	return pickRandom([]LineNavigator{goDown{}, goUp{}})
}
func (nav goLeft) Reverse() LineNavigator {
	return goRight{}
}
func (nav goLeft) IsHorizontal() bool {
	return true
}

type goRight struct {
}

func (nav goRight) NextPoint(p Point) Point {
	return p.right()
}
func (nav goRight) Alternate() LineNavigator {
	return pickRandom([]LineNavigator{goDown{}, goUp{}})
}
func (nav goRight) Reverse() LineNavigator {
	return goLeft{}
}
func (nav goRight) IsHorizontal() bool {
	return true
}

type goUp struct {
}

func (nav goUp) NextPoint(p Point) Point {
	return p.up()
}
func (nav goUp) Alternate() LineNavigator {
	return pickRandom([]LineNavigator{goLeft{}, goRight{}})
}
func (nav goUp) Reverse() LineNavigator {
	return goDown{}
}
func (nav goUp) IsHorizontal() bool {
	return false
}

type goDown struct {
}

func (nav goDown) NextPoint(p Point) Point {
	return p.down()
}
func (nav goDown) Alternate() LineNavigator {
	return pickRandom([]LineNavigator{goLeft{}, goRight{}})
}
func (nav goDown) Reverse() LineNavigator {
	return goUp{}
}
func (nav goDown) IsHorizontal() bool {
	return false
}

func pickRandom(navs []LineNavigator) LineNavigator {
	rand.Seed(time.Now().UnixNano())
	return navs[rand.Intn(len(navs))]
}

func MakeRandomLineNavigator() LineNavigator {
	navs := []LineNavigator{
		goRight{},
		goDown{},
		goLeft{},
		goUp{},
	}
	return pickRandom(navs)
}
