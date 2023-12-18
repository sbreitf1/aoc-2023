package main

// https://adventofcode.com/2023/day/17

import (
	"aoc/helper"
	"fmt"
	"sort"
	"strings"
)

const (
	MaxInt int = 20000000
)

func main() {
	lines := helper.ReadNonEmptyLines("example-1.txt")

	board := ParseBoard(lines)
	result1 := board.FindMinPathCost(helper.Point2D{X: 0, Y: 0}, helper.Point2D{X: board.Width - 1, Y: board.Height - 1})
	PrintPath(board, result1.Path)
	solution1 := result1.Cost
	//path := board.FindPath(helper.Point2D{X: 0, Y: 0}, helper.Point2D{X: board.Width - 1, Y: board.Height - 1})
	//PrintPath(board, path)
	//solution1 := board.GetPathHeatLoss(path)
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

func (b *Board) FindMinPathCost(from, to helper.Point2D) Result {
	cache := Cache{
		Visiting:     make(map[CacheKey]bool),
		KnownResults: make(map[CacheKey]Result),
	}
	return b.findMinPathCost(&cache, from, to, helper.Point2D{}, 1)
}

type Cache struct {
	Visiting     map[CacheKey]bool
	KnownResults map[CacheKey]Result
}

type CacheKey struct {
	Pos              helper.Point2D
	PreviousDir      helper.Point2D
	PreviousDirCount int
}

type Result struct {
	Cost int
	Path []helper.Point2D
}

func (b *Board) findMinPathCost(cache *Cache, from, to helper.Point2D, previousDir helper.Point2D, previousDirCount int) Result {
	if from.X < 0 || from.Y < 0 || from.X >= b.Width || from.Y >= b.Height {
		return Result{MaxInt, nil}
	}
	if from == to {
		return Result{b.Tiles[from.Y][from.X], []helper.Point2D{to}}
	}

	key := CacheKey{
		Pos:              from,
		PreviousDir:      previousDir,
		PreviousDirCount: previousDirCount,
	}
	if knownResult, ok := cache.KnownResults[key]; ok {
		return knownResult
	}

	if ok := cache.Visiting[key]; ok {
		return Result{MaxInt, nil}
	}
	cache.Visiting[key] = true

	nextDirs := make([]helper.Point2D, 0, 3)
	for _, nextDir := range []helper.Point2D{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
		if nextDir == previousDir && previousDirCount >= 3 {
			continue
		}

		if nextDir == previousDir.Neg() {
			continue
		}

		nextPos := from.Add(nextDir)
		if nextPos.X < 0 || nextPos.Y < 0 || nextPos.X >= b.Width || nextPos.Y >= b.Height {
			continue
		}

		nextDirs = append(nextDirs, nextDir)
	}
	sort.Slice(nextDirs, func(i, j int) bool {
		return b.Tiles[from.Y+nextDirs[i].Y][from.X+nextDirs[i].X] < b.Tiles[from.Y+nextDirs[j].Y][from.X+nextDirs[j].X]
	})

	result := Result{MaxInt, nil}
	for _, nextDir := range nextDirs {
		if nextDir == previousDir && previousDirCount >= 3 {
			continue
		}

		if nextDir == previousDir.Neg() {
			continue
		}

		nextPos := from.Add(nextDir)
		if nextPos.X < 0 || nextPos.Y < 0 || nextPos.X >= b.Width || nextPos.Y >= b.Height {
			continue
		}

		var nextSameDirStepCount int
		if nextDir == previousDir {
			nextSameDirStepCount = previousDirCount + 1
		} else {
			nextSameDirStepCount = 1
		}

		nextResult := b.findMinPathCost(cache, nextPos, to, nextDir, nextSameDirStepCount)
		if nextResult.Cost < result.Cost {
			result = nextResult
		}
	}
	result.Cost += b.Tiles[from.Y][from.X]
	result.Path = append([]helper.Point2D{from}, result.Path...)
	cache.KnownResults[key] = result
	cache.Visiting[key] = false
	return result
}

type PathPoint struct {
	Previous         *PathPoint
	Pos              helper.Point2D
	TotalCost        int
	SameDirStepCount int
}

func (b *Board) FindPath(from, to helper.Point2D) []helper.Point2D {
	nextValues := map[helper.Point2D]*PathPoint{
		{X: 0, Y: 0}: {Pos: helper.Point2D{X: 0, Y: 0}, Previous: nil, TotalCost: 0, SameDirStepCount: 1},
	}
	visited := map[helper.Point2D]*PathPoint{}

	for len(nextValues) > 0 {
		var currentPoint *PathPoint
		for _, p := range nextValues {
			if currentPoint == nil || p.TotalCost < currentPoint.TotalCost {
				currentPoint = p
			}
		}
		delete(nextValues, currentPoint.Pos)

		if v, ok := visited[currentPoint.Pos]; ok {
			if currentPoint.TotalCost < v.TotalCost {
				helper.ExitWithMessage("found better way to %v (%d -> %d)", currentPoint.Pos, v.TotalCost, currentPoint.TotalCost)
			}
			continue
		}
		visited[currentPoint.Pos] = currentPoint

		var previousDir helper.Point2D
		if currentPoint.Previous != nil {
			previousDir = currentPoint.Pos.Sub(currentPoint.Previous.Pos)
		}
		for _, nextDir := range []helper.Point2D{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
			if nextDir == previousDir && currentPoint.SameDirStepCount >= 3 {
				continue
			}

			if nextDir == previousDir.Neg() {
				continue
			}

			nextPos := currentPoint.Pos.Add(nextDir)
			if nextPos.X < 0 || nextPos.Y < 0 || nextPos.X >= b.Width || nextPos.Y >= b.Height {
				continue
			}

			if _, ok := visited[nextPos]; ok {
				continue
			}

			var nextSameDirStepCount int
			if nextDir == previousDir {
				nextSameDirStepCount = currentPoint.SameDirStepCount + 1
			} else {
				nextSameDirStepCount = 1
			}

			nextPoint := PathPoint{
				Pos:              nextPos,
				Previous:         currentPoint,
				TotalCost:        currentPoint.TotalCost + b.Tiles[nextPos.Y][nextPos.X],
				SameDirStepCount: nextSameDirStepCount,
			}
			if alreadyEnqueuedPoint, ok := nextValues[nextPos]; ok {
				if nextPoint.TotalCost < alreadyEnqueuedPoint.TotalCost {
					nextValues[nextPoint.Pos] = &nextPoint
				}
			} else {
				nextValues[nextPoint.Pos] = &nextPoint
			}
		}
	}

	endPos, ok := visited[to]
	if !ok {
		helper.ExitWithMessage("no path from %v to %v found", from, to)
	}

	path := []helper.Point2D{endPos.Pos}
	for {
		pp := visited[path[0]]
		if pp.Previous == nil {
			break
		}
		fmt.Println(pp)
		path = append([]helper.Point2D{pp.Previous.Pos}, path...)
	}
	fmt.Println(path)
	return path
}

func (b *Board) GetPathHeatLoss(path []helper.Point2D) int {
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
