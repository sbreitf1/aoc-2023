package main

// https://adventofcode.com/2023/day/23

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadNonEmptyLines("example-1.txt")

	world := ParseWorld(lines)
	solution1 := world.FindLongestPathLengthFromTo(helper.Point2D{X: 1, Y: 0}, helper.Point2D{X: world.Width - 2, Y: world.Height - 1}, false)
	fmt.Println("-> part 1:", solution1)
	//solution2 := world.FindLongestPathLengthFromTo(helper.Point2D{X: 1, Y: 0}, helper.Point2D{X: world.Width - 2, Y: world.Height - 1}, true)
	//fmt.Println("-> part 2:", solution2)
}

func ParseWorld(lines []string) *World {
	return &World{
		Width:  len(lines[0]),
		Height: len(lines),
		Tiles:  helper.LinesToRunes(lines),
	}
}

type World struct {
	Width, Height int
	Tiles         [][]rune
}

func (w *World) FindLongestPathLengthFromTo(from, to helper.Point2D, part2 bool) int64 {
	visited := make([][]bool, len(w.Tiles))
	for y := 0; y < len(visited); y++ {
		visited[y] = make([]bool, len(w.Tiles[y]))
	}

	maxPathLength, ok := w.findLongestPathLengthFromToRecursive(visited, from, to, part2)
	if !ok {
		helper.ExitWithMessage("no path found!")
	}
	return maxPathLength
}

func (w *World) findLongestPathLengthFromToRecursive(visited [][]bool, from, to helper.Point2D, part2 bool) (int64, bool) {
	if visited[from.Y][from.X] {
		return 0, false
	}

	if from == to {
		return 0, true
	}

	visited[from.Y][from.X] = true
	defer func() {
		visited[from.Y][from.X] = false
	}()

	walkableChars := map[helper.Point2D]map[rune]bool{
		{X: 0, Y: -1}: {'.': true, '^': true},
		{X: 0, Y: 1}:  {'.': true, 'v': true},
		{X: 1, Y: 0}:  {'.': true, '>': true},
		{X: -1, Y: 0}: {'.': true, '<': true},
	}

	var maxPathLength int64
	var ok bool
	for _, dir := range []helper.Point2D{{X: 0, Y: -1}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: -1, Y: 0}} {
		newPos := from.Add(dir)
		if newPos.X < 0 || newPos.Y < 0 || newPos.X >= w.Width || newPos.Y >= w.Height {
			continue
		}

		if part2 {
			if w.Tiles[newPos.Y][newPos.X] != '#' {
				continue
			}

		} else {
			if !walkableChars[dir][w.Tiles[newPos.Y][newPos.X]] {
				continue
			}
		}

		if nextMaxPathLength, nextOK := w.findLongestPathLengthFromToRecursive(visited, newPos, to, part2); nextOK {
			if nextMaxPathLength >= maxPathLength {
				ok = true
				maxPathLength = nextMaxPathLength + 1
			}
		}
	}

	if !ok {
		return 0, false
	}
	return maxPathLength, true
}
