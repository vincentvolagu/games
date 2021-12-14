package main

import (
	"math/rand"
	"time"
)

type LineNavigator interface {
	Reset()
	NextPoint(p Point) Point
}

type horizontal struct {
}

func (f horizontal) NextPoint(p Point) Point {
	return p.right()
}
func (f horizontal) Reset() {
}

type vertical struct {
}

func (f vertical) NextPoint(p Point) Point {
	return p.down()
}
func (f vertical) Reset() {
}

type randomLine struct {
	nav LineNavigator
}

func (f *randomLine) Reset() {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100) % 2
	if randNum == 0 {
		f.nav = horizontal{}
	} else {
		f.nav = vertical{}
	}
}
func (f *randomLine) NextPoint(p Point) Point {
	if f.nav == nil {
		f.Reset()
	}
	return f.nav.NextPoint(p)
}
