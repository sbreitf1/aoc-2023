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
	path1 := board.FindPath(helper.Point2D[int]{X: 0, Y: 0}, helper.Point2D[int]{X: board.Width - 1, Y: board.Height - 1}, 1, 3)
	//PrintPath(board, path1)
	solution1 := board.GetPathHeatLoss(path1)
	fmt.Println("-> part 1:", solution1)

	path2 := board.FindPath(helper.Point2D[int]{X: 0, Y: 0}, helper.Point2D[int]{X: board.Width - 1, Y: board.Height - 1}, 4, 10)
	//PrintPath(board, path2)
	solution2 := board.GetPathHeatLoss(path2)
	fmt.Println("-> part 2:", solution2)
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
	Pos              helper.Point2D[int]
	TotalCost        int
	SameDirStepCount int
}

func (b *Board) FindPath(from, to helper.Point2D[int], minDist, maxDist int) []helper.Point2D[int] {
	type VisitKey struct {
		Pos              helper.Point2D[int]
		InDir            helper.Point2D[int]
		SameDirStepCount int
	}
	nextValues := helper.MakePriorityQueue[int, *PathPoint]()
	nextValues.Push(0, &PathPoint{Pos: from, Previous: nil, TotalCost: 0, SameDirStepCount: 1})
	visited := map[VisitKey]*PathPoint{}

	var bestEndPos *PathPoint

	for nextValues.Len() > 0 {
		currentPoint, _ := nextValues.Pop()

		var inDir helper.Point2D[int]
		if currentPoint.Previous != nil {
			inDir = currentPoint.Pos.Sub(currentPoint.Previous.Pos)
		}

		vkey := VisitKey{Pos: currentPoint.Pos, InDir: inDir, SameDirStepCount: currentPoint.SameDirStepCount}

		if v, ok := visited[vkey]; ok {
			if currentPoint.TotalCost < v.TotalCost {
				helper.ExitWithMessage("found better way to %v (%d -> %d)", currentPoint.Pos, v.TotalCost, currentPoint.TotalCost)
			}
			continue
		}
		visited[vkey] = currentPoint

		if currentPoint.Pos == to && currentPoint.SameDirStepCount >= minDist {
			if bestEndPos == nil || currentPoint.TotalCost < bestEndPos.TotalCost {
				bestEndPos = currentPoint
			}
		}

		for _, nextDir := range []helper.Point2D[int]{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
			if nextDir == inDir {
				if currentPoint.SameDirStepCount >= maxDist {
					continue
				}
			} else if nextDir == inDir.Neg() {
				continue
			} else if (inDir == helper.Point2D[int]{}) {
				// accept
			} else {
				if currentPoint.SameDirStepCount < minDist {
					continue
				}
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
			nextValues.Push(nextPoint.TotalCost, &nextPoint)
		}
	}

	if bestEndPos == nil {
		helper.ExitWithMessage("no path from %v to %v found", from, to)
	}

	path := []helper.Point2D[int]{}
	for current := bestEndPos; current != nil; current = current.Previous {
		path = append([]helper.Point2D[int]{current.Pos}, path...)
	}
	return path
}

func (b *Board) GetPathHeatLoss(path []helper.Point2D[int]) int {
	var heatLoss int
	for i := 1; i < len(path); i++ {
		heatLoss += b.Tiles[path[i].Y][path[i].X]
	}
	return heatLoss
}

func PrintPath(board *Board, path []helper.Point2D[int]) {
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
