package main

// https://adventofcode.com/2023/day/21

import (
	"aoc/helper"
	"fmt"
	"os"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	garden := ParseGarden(lines)
	/*solution1 := garden.CountPossiblePositionsFromStartPos(64, false, false)
	fmt.Println("-> part 1:", solution1)*/

	//solution2 := garden.CountPossiblePositionsFromStartPos(327, true)
	// 64  -> 3697
	// 65  -> 3762
	// 66  -> 3961
	// 80  -> 5763
	// 100 -> 8864
	// 129 -> 14624
	// 130 -> 14838
	// 131 -> 15055
	// 132 -> 15273
	// 150 -> 19644
	// 196 -> 33547 #
	// 200 -> 35083
	// 262 -> 59829
	// 327 -> 93052 #
	// 333 -> 96725
	// 458 -> 182277 #
	// 500 -> 217446
	// 589 -> 301222 #
	// 600 -> 313463
	// 709 -> 436047
	// 720 -> 449887 #
	// 811 -> 570424
	solution2 := garden.CountPossiblePositionsWithRepeat(589)
	fmt.Println("-> part 2:", solution2)
	fmt.Println("expected is 301222")
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

func (g Garden) CountPossiblePositionsFromStartPos(steps int64, repeatX, repeatY bool) int64 {
	return g.CountPossiblePositionsFromPos(g.StartPos, steps, repeatX, repeatY)
}

func (g Garden) CountPossiblePositionsFromPos(startPos helper.Point2D, steps int64, repeatX, repeatY bool) int64 {
	type VisitKey struct {
		Pos            helper.Point2D
		RemainingSteps int64
	}
	nextSteps := []VisitKey{{startPos, steps}}
	visited := make(map[VisitKey]bool)
	//visitedPoints := make(map[helper.Point2D]int64)
	for len(nextSteps) > 0 {
		p := nextSteps[len(nextSteps)-1]
		nextSteps = nextSteps[:len(nextSteps)-1]

		if _, ok := visited[p]; ok {
			continue
		}
		visited[p] = true

		/*if lastRemaining, ok := visitedPoints[p.Pos]; ok {
			if lastRemaining <= p.RemainingSteps {
				continue
			} else {
				visitedPoints[p.Pos] = p.RemainingSteps
			}
		} else {
			visitedPoints[p.Pos] = p.RemainingSteps
		}*/

		if p.RemainingSteps <= 0 {
			continue
		}

		for _, dir := range []helper.Point2D{{X: 0, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}} {
			nextPos := p.Pos.Add(dir)
			if !repeatX && (nextPos.X < 0 || nextPos.X >= g.Width) {
				continue
			}
			if !repeatY && (nextPos.Y < 0 || nextPos.Y >= g.Height) {
				continue
			}

			if g.Tiles[helper.Mod(nextPos.Y, g.Height)][helper.Mod(nextPos.X, g.Width)] == '#' {
				continue
			}

			nextSteps = append(nextSteps, VisitKey{Pos: nextPos, RemainingSteps: p.RemainingSteps - 1})
		}
	}
	tiles := make([][]rune, len(g.Tiles))
	for y := 0; y < g.Height; y++ {
		tiles[y] = make([]rune, len(g.Tiles[y]))
		copy(tiles[y], g.Tiles[y])
	}
	var count int64
	for v := range visited {
		if v.RemainingSteps == 0 {
			count++
			if v.Pos.X >= 0 && v.Pos.Y >= 0 && v.Pos.X < g.Width && v.Pos.Y < g.Height {
				tiles[v.Pos.Y][v.Pos.X] = 'O'
			}
		}
	}
	buf := make([]string, len(tiles))
	for y := 0; y < g.Height; y++ {
		buf[y] = string(tiles[y])
	}
	os.WriteFile("out.txt", []byte(strings.Join(buf, "\n")), os.ModePerm)
	return count
}

func (g Garden) CountPossiblePositionsWithRepeat(steps int64) int64 {
	// check optimized conditions
	if g.Width != (2*g.StartPos.X)+1 || g.Height != (2*g.StartPos.Y)+1 {
		fmt.Println("start pos not in center of field, optimization not available")
		return g.CountPossiblePositionsFromStartPos(steps, true, true)
	}
	/*if g.Width != g.Height {
		fmt.Println("field is not square, optimization not available")
		return g.CountPossiblePositionsFromStartPos(steps, true, true)
	}*/
	for y := 0; y < g.Height; y++ {
		if g.Tiles[y][g.StartPos.X] == '#' {
			fmt.Println("no direct path to border, optimization not available")
			return g.CountPossiblePositionsFromStartPos(steps, true, true)
		}
		if g.Tiles[y][0] != '.' || g.Tiles[y][g.Width-1] != '.' {
			fmt.Println("border not free, optimization not available")
			return g.CountPossiblePositionsFromStartPos(steps, true, true)
		}
	}
	for x := 0; x < g.Width; x++ {
		if g.Tiles[g.StartPos.Y][x] == '#' {
			fmt.Println("no direct path to border, optimization not available")
			return g.CountPossiblePositionsFromStartPos(steps, true, true)
		}
		if g.Tiles[0][x] != '.' || g.Tiles[g.Height-1][x] != '.' {
			fmt.Println("border not free, optimization not available")
			return g.CountPossiblePositionsFromStartPos(steps, true, true)
		}
	}
	if steps < int64(g.StartPos.X+g.StartPos.Y+2) {
		fmt.Println("step count too small, optimization not available")
		return g.CountPossiblePositionsFromStartPos(steps, true, true)
	}
	/*if steps%int64(g.Width) != int64(g.StartPos.X) {
		fmt.Println("step count does not match, optimization not available")
		return g.CountPossiblePositionsFromStartPos(steps, true, true)
	}*/

	// conditions for optimized computation are met

	fullFieldCountEven := g.CountPossiblePositionsFromStartPos(int64(2*g.StartPos.X+2*g.StartPos.Y)+(steps%2), false, false)
	fullFieldCountOdd := g.CountPossiblePositionsFromStartPos(int64(2*g.StartPos.X+2*g.StartPos.Y)+(steps%2)+1, false, false)
	/*fieldCountLeft := g.CountPossiblePositionsFromPos(helper.Point2D{X: g.Width - 1, Y: g.StartPos.Y}, int64(g.Width-1), false)
	fieldCountRight := g.CountPossiblePositionsFromPos(helper.Point2D{X: 0, Y: g.StartPos.Y}, int64(g.Width-1), false)
	fieldCountTop := g.CountPossiblePositionsFromPos(helper.Point2D{X: g.StartPos.X, Y: g.Height - 1}, int64(g.Height-1), false)
	fieldCountBottom := g.CountPossiblePositionsFromPos(helper.Point2D{X: g.StartPos.X, Y: 0}, int64(g.Height-1), false)
	fmt.Println(fullFieldCountOdd, fullFieldCountEven, fieldCountLeft, fieldCountRight, fieldCountTop, fieldCountBottom)*/

	var sum int64

	//fmt.Println(g.CountPossiblePositionsFromStartPos(steps, true, false))
	// = 63202
	fullXCountOneDir := (steps - int64(g.StartPos.X) - int64(g.StartPos.Y)) / int64(g.Width)
	sum += fullFieldCountEven + 2*((fullXCountOneDir/2+fullXCountOneDir%2)*fullFieldCountOdd+(fullXCountOneDir/2)*fullFieldCountEven)
	fmt.Println(helper.Mod(int64(g.StartPos.X)-fullXCountOneDir*int64(g.Width), int64(g.Width)))
	sum += g.CountPossiblePositionsFromPos(helper.Point2D{}, steps, false, false)

	return sum
}
