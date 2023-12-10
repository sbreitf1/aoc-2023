package main

// https://adventofcode.com/2023/day/9

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	world := ParseWorld(lines)
	solution1 := world.FindMaxPathToAnimal()

	fmt.Println("-> part 1:", solution1)
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

func (w *World) VisitedWorldString() string {
	lines := make([]string, 0)
	for y := 0; y < w.Height; y++ {
		var line string
		for x := 0; x < w.Width; x++ {
			if w.Tiles[y][x].StepsToAnimal >= 0 {
				line += string(w.Tiles[y][x].Rune)
			} else {
				line += "."
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
