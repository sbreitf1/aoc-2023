package main

// https://adventofcode.com/2023/day/9

import (
	"aoc/helper"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	world := ParseWorld(lines)
	solution1 := world.FindMaxPathToAnimal()
	solution2 := world.CountEmptyFieldsWithNonZeroWindingNumber()

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Point struct {
	X, Y int
}

type World struct {
	Width, Height int
	Tiles         [][]Tile
	Animal        Point
}

type Tile struct {
	Rune          rune
	StepsToAnimal int
	PartOfLoop    bool
	Enclosed      bool
}

func (t Tile) ConnectsToWest() bool {
	return t.Rune == 'S' || t.Rune == '-' || t.Rune == 'J' || t.Rune == '7'
}
func (t Tile) ConnectsToNorth() bool {
	return t.Rune == 'S' || t.Rune == '|' || t.Rune == 'L' || t.Rune == 'J'
}
func (t Tile) ConnectsToEast() bool {
	return t.Rune == 'S' || t.Rune == '-' || t.Rune == 'L' || t.Rune == 'F'
}
func (t Tile) ConnectsToSouth() bool {
	return t.Rune == 'S' || t.Rune == '|' || t.Rune == '7' || t.Rune == 'F'
}

func ParseWorld(lines []string) World {
	var world World
	world.Tiles = make([][]Tile, len(lines))
	for y := 0; y < len(lines); y++ {
		world.Tiles[y] = make([]Tile, len(lines[y]))
		if len(world.Tiles[y]) != len(world.Tiles[0]) {
			helper.ExitWithMessage("mismatching line length in line %d", (y + 1))
		}
		for x, r := range lines[y] {
			world.Tiles[y][x] = Tile{
				Rune:          r,
				StepsToAnimal: -1,
			}
			if r == 'S' {
				world.Animal = Point{X: x, Y: y}
			}
		}
	}
	world.Height = len(world.Tiles)
	world.Width = len(world.Tiles[0])
	return world
}

func (w *World) FindMaxPathToAnimal() int {
	nextVisit := []Point{w.Animal}
	w.Tiles[w.Animal.Y][w.Animal.X].StepsToAnimal = 0
	visited := map[Point]bool{
		w.Animal: true,
	}
	var maxSteps int
	for len(nextVisit) > 0 {
		t := nextVisit[0]
		nextVisit = nextVisit[1:]
		if w.Tiles[t.Y][t.X].StepsToAnimal > maxSteps {
			maxSteps = w.Tiles[t.Y][t.X].StepsToAnimal
		}

		west := Point{t.X - 1, t.Y}
		if w.CanMoveWest(t.X, t.Y) && !visited[west] {
			w.Tiles[west.Y][west.X].StepsToAnimal = w.Tiles[t.Y][t.X].StepsToAnimal + 1
			nextVisit = append(nextVisit, west)
			visited[west] = true
		}

		north := Point{t.X, t.Y - 1}
		if _, ok := visited[north]; !ok && w.CanMoveNorth(t.X, t.Y) {
			w.Tiles[north.Y][north.X].StepsToAnimal = w.Tiles[t.Y][t.X].StepsToAnimal + 1
			nextVisit = append(nextVisit, north)
			visited[north] = true
		}

		east := Point{t.X + 1, t.Y}
		if _, ok := visited[east]; !ok && w.CanMoveEast(t.X, t.Y) {
			w.Tiles[east.Y][east.X].StepsToAnimal = w.Tiles[t.Y][t.X].StepsToAnimal + 1
			nextVisit = append(nextVisit, east)
			visited[east] = true
		}

		south := Point{t.X, t.Y + 1}
		if _, ok := visited[south]; !ok && w.CanMoveSouth(t.X, t.Y) {
			w.Tiles[south.Y][south.X].StepsToAnimal = w.Tiles[t.Y][t.X].StepsToAnimal + 1
			nextVisit = append(nextVisit, south)
			visited[south] = true
		}
	}
	//fmt.Println(w.VisitedWorldString())
	return maxSteps
}

func (w *World) CanMoveWest(x, y int) bool {
	return x > 0 && w.Tiles[y][x].ConnectsToWest() && w.Tiles[y][x-1].ConnectsToEast()
}
func (w *World) CanMoveNorth(x, y int) bool {
	return y > 0 && w.Tiles[y][x].ConnectsToNorth() && w.Tiles[y-1][x].ConnectsToSouth()
}
func (w *World) CanMoveEast(x, y int) bool {
	return x < (w.Width-1) && w.Tiles[y][x].ConnectsToEast() && w.Tiles[y][x+1].ConnectsToWest()
}
func (w *World) CanMoveSouth(x, y int) bool {
	return y < (w.Height-1) && w.Tiles[y][x].ConnectsToSouth() && w.Tiles[y+1][x].ConnectsToNorth()
}

func (w *World) ExtractLoop() []Point {
	if w.Tiles[w.Animal.Y][w.Animal.X].StepsToAnimal != 0 {
		helper.ExitWithMessage("use World.FindMaxPathToAnimal before World.ExtractLoop")
	}

	visited := map[Point]bool{}

	loop := []Point{w.Animal}
	for {
		t := loop[len(loop)-1]
		w.Tiles[t.Y][t.X].PartOfLoop = true
		visited[t] = true

		west := Point{t.X - 1, t.Y}
		if w.CanMoveWest(t.X, t.Y) && w.Tiles[west.Y][west.X].StepsToAnimal >= 0 {
			if len(loop) > 2 && west == w.Animal {
				break
			}
			if !visited[west] {
				loop = append(loop, west)
				continue
			}
		}

		north := Point{t.X, t.Y - 1}
		if w.CanMoveNorth(t.X, t.Y) && w.Tiles[north.Y][north.X].StepsToAnimal >= 0 {
			if len(loop) > 2 && north == w.Animal {
				break
			}
			if !visited[north] {
				loop = append(loop, north)
				continue
			}
		}

		east := Point{t.X + 1, t.Y}
		if w.CanMoveEast(t.X, t.Y) && w.Tiles[east.Y][east.X].StepsToAnimal >= 0 {
			if len(loop) > 2 && east == w.Animal {
				break
			}
			if !visited[east] {
				loop = append(loop, east)
				continue
			}
		}

		south := Point{t.X, t.Y + 1}
		if w.CanMoveSouth(t.X, t.Y) && w.Tiles[south.Y][south.X].StepsToAnimal >= 0 {
			if len(loop) > 2 && south == w.Animal {
				break
			}
			if !visited[south] {
				loop = append(loop, south)
				continue
			}
		}
	}
	return loop
}

func (w *World) String() string {
	lines := []string{}
	for y := 0; y < w.Height; y++ {
		var line string
		for x := 0; x < w.Width; x++ {
			t := w.Tiles[y][x]
			if t.Enclosed {
				line += "I"
			} else if t.PartOfLoop {
				line += string(t.Rune)
			} else {
				line += "."
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (w *World) CountEmptyFieldsWithNonZeroWindingNumber() int {
	loop := w.ExtractLoop()
	candidates := make([]Point, 0)
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Tiles[y][x].Rune == '.' || !w.Tiles[y][x].PartOfLoop {
				candidates = append(candidates, Point{x, y})
			}
		}
	}

	var count int32
	var wg sync.WaitGroup
	for _, p := range candidates {
		wg.Add(1)
		go func(p Point) {
			defer wg.Done()

			windingNumber := ComputeWindingNumber(p, loop)
			if windingNumber > 0.1 || windingNumber < -0.1 {
				w.Tiles[p.Y][p.X].Enclosed = true
				atomic.AddInt32(&count, 1)
			}
		}(p)
	}
	wg.Wait()
	return int(count)
}

func ComputeWindingNumber(p Point, loop []Point) float64 {
	var windingNumber float64
	for i := 0; i < len(loop); i++ {
		l1 := loop[i]
		l2 := loop[(i+1)%len(loop)]
		w := ComputeWindingNumberOfLine(p, l1, l2)
		windingNumber += w
	}
	return windingNumber
}

func ComputeWindingNumberOfLine(p Point, l1, l2 Point) float64 {
	if p.X == l1.X && p.X == l2.X {
		return 0
	}
	if p.Y == l1.Y && p.Y == l2.Y {
		return 0
	}
	a1 := math.Atan2(float64(p.Y)-float64(l1.Y), float64(p.X)-float64(l1.X))
	a2 := math.Atan2(float64(p.Y)-float64(l2.Y), float64(p.X)-float64(l2.X))
	diff := a2 - a1
	if diff >= math.Pi {
		diff -= 2 * math.Pi
	}
	if diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	return diff
}
