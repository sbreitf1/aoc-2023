package main

// https://adventofcode.com/2023/day/17

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	board := ParseBoard(lines)
	path := board.FindPath(helper.Point2D{X: 0, Y: 0}, helper.Point2D{X: board.Width - 1, Y: board.Height - 1})
	PrintPath(board, path)
	solution1 := board.GetPathHeatLoss(path)
	fmt.Println("-> part 1:", solution1)
}

type Board struct {
	Width, Height int
	Tiles         [][]int
}

func ParseBoard(lines []string) *Board {
	tiles := make([][]int, len(lines))
	for y := range lines {
		tiles[y] = make([]int, len(lines[y]))
		for x := range lines[y] {
			tiles[y][x] = int(lines[y][x] - '0')
		}
	}
	return &Board{
		Width:  len(tiles[0]),
		Height: len(tiles),
		Tiles:  tiles,
	}
}

type PathPoint struct {
	Previous         *PathPoint
	Pos              helper.Point2D
	TotalCost        int
	SameDirStepCount int
}

func (b *Board) FindPath(from, to helper.Point2D) []helper.Point2D {
	type VisitKey struct {
		Pos              helper.Point2D
		InDir            helper.Point2D
		SameDirStepCount int
	}
	nextValues := map[VisitKey]*PathPoint{
		{Pos: from, SameDirStepCount: 1}: {Pos: from, Previous: nil, TotalCost: 0, SameDirStepCount: 1},
	}
	visited := map[VisitKey]*PathPoint{}

	var bestEndPos *PathPoint

	for len(nextValues) > 0 {
		var currentPoint *PathPoint
		var delKey VisitKey
		for k, p := range nextValues {
			if currentPoint == nil || p.TotalCost < currentPoint.TotalCost {
				currentPoint = p
				delKey = k
			}
		}
		delete(nextValues, delKey)

		var inDir helper.Point2D
		if currentPoint.Previous != nil {
			inDir = currentPoint.Pos.Sub(currentPoint.Previous.Pos)
		}

		vkey := VisitKey{Pos: currentPoint.Pos, InDir: inDir, SameDirStepCount: currentPoint.SameDirStepCount}
		if delKey != vkey {
			fmt.Println(delKey, vkey)
		}

		if v, ok := visited[vkey]; ok {
			if currentPoint.TotalCost < v.TotalCost {
				helper.ExitWithMessage("found better way to %v (%d -> %d)", currentPoint.Pos, v.TotalCost, currentPoint.TotalCost)
			}
			continue
		}
		visited[vkey] = currentPoint

		if currentPoint.Pos == to {
			if bestEndPos == nil || currentPoint.TotalCost < bestEndPos.TotalCost {
				bestEndPos = currentPoint
			}
		}

		for _, nextDir := range []helper.Point2D{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
			if nextDir == inDir && currentPoint.SameDirStepCount >= 3 {
				continue
			}

			if nextDir == inDir.Neg() {
				continue
			}

			nextPos := currentPoint.Pos.Add(nextDir)
			if nextPos.X < 0 || nextPos.Y < 0 || nextPos.X >= b.Width || nextPos.Y >= b.Height {
				continue
			}

			var nextSameDirStepCount int
			if nextDir == inDir {
				nextSameDirStepCount = currentPoint.SameDirStepCount + 1
			} else {
				nextSameDirStepCount = 1
			}

			vkeyNext := VisitKey{Pos: nextPos, InDir: nextDir, SameDirStepCount: nextSameDirStepCount}
			if _, ok := visited[vkeyNext]; ok {
				continue
			}

			nextPoint := PathPoint{
				Pos:              nextPos,
				Previous:         currentPoint,
				TotalCost:        currentPoint.TotalCost + b.Tiles[nextPos.Y][nextPos.X],
				SameDirStepCount: nextSameDirStepCount,
			}
			if alreadyEnqueuedPoint, ok := nextValues[vkeyNext]; ok {
				if nextPoint.TotalCost < alreadyEnqueuedPoint.TotalCost {
					nextValues[vkeyNext] = &nextPoint
				}
			} else {
				nextValues[vkeyNext] = &nextPoint
			}
		}
	}

	if bestEndPos == nil {
		helper.ExitWithMessage("no path from %v to %v found", from, to)
	}

	path := []helper.Point2D{}
	for current := bestEndPos; current != nil; current = current.Previous {
		path = append([]helper.Point2D{current.Pos}, path...)
	}
	return path
}

func (b *Board) GetPathHeatLoss(path []helper.Point2D) int {
	fmt.Println(path)
	var heatLoss int
	for i := 1; i < len(path); i++ {
		heatLoss += b.Tiles[path[i].Y][path[i].X]
	}
	return heatLoss
}

func PrintPath(board *Board, path []helper.Point2D) {
	runeLines := make([][]rune, board.Height)
	for y := 0; y < board.Height; y++ {
		runeLines[y] = make([]rune, len(board.Tiles[y]))
		for x := 0; x < board.Width; x++ {
			runeLines[y][x] = '0' + rune(board.Tiles[y][x])
		}
	}
	for _, p := range path {
		runeLines[p.Y][p.X] = '#'
	}
	lines := make([]string, len(runeLines))
	for y := 0; y < board.Height; y++ {
		lines[y] = string(runeLines[y])
	}
	fmt.Println(strings.Join(lines, "\n"))
}
