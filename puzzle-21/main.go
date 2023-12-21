package main

// https://adventofcode.com/2023/day/21

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	garden := ParseGarden(lines)
	solution1 := garden.CountPossiblePositions(64, false)
	fmt.Println("-> part 1:", solution1)

	solution2 := garden.CountPossiblePositions(100, true)
	// 100 -> 8864
	// 200 -> 35083
	// 333 -> 96725
	// 500 -> 217446
	// 600 -> 313463
	// 709 -> 436047
	// 811 -> 570424
	fmt.Println("-> part 2:", solution2)
}

func ParseGarden(lines []string) Garden {
	tiles := make([][]rune, len(lines))
	var startPos helper.Point2D
	for y := range lines {
		tiles[y] = []rune(lines[y])
		for x := range tiles[y] {
			if tiles[y][x] == 'S' {
				startPos = helper.Point2D{X: x, Y: y}
			}
		}
	}
	return Garden{
		Width:    len(tiles[0]),
		Height:   len(tiles),
		Tiles:    tiles,
		StartPos: startPos,
	}
}

type Garden struct {
	Width, Height int
	Tiles         [][]rune
	StartPos      helper.Point2D
}

func (g Garden) CountPossiblePositions(steps int, repeat bool) int64 {
	type VisitKey struct {
		Pos            helper.Point2D
		RemainingSteps int
	}
	nextSteps := []VisitKey{{helper.Point2D{X: 0, Y: 0}, steps}}
	visited := make(map[VisitKey]bool)
	for len(nextSteps) > 0 {
		p := nextSteps[len(nextSteps)-1]
		nextSteps = nextSteps[:len(nextSteps)-1]

		if _, ok := visited[p]; ok {
			continue
		}
		visited[p] = true

		if p.RemainingSteps <= 0 {
			continue
		}

		for _, dir := range []helper.Point2D{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
			nextPos := p.Pos.Add(dir)
			absPos := nextPos.Add(g.StartPos)
			if !repeat && (absPos.X < 0 || absPos.Y < 0 || absPos.X >= g.Width || absPos.Y >= g.Height) {
				continue
			}

			if g.Tiles[modLikePython(absPos.Y, g.Height)][modLikePython(absPos.X, g.Width)] == '#' {
				continue
			}

			nextSteps = append(nextSteps, VisitKey{Pos: nextPos, RemainingSteps: p.RemainingSteps - 1})
		}
	}
	var count int64
	for v := range visited {
		if v.RemainingSteps == 0 {
			count++
		}
	}
	return count
}

func modLikePython(d, m int) int {
	var res int = d % m
	if res < 0 && m > 0 {
		return res + m
	}
	return res
}
