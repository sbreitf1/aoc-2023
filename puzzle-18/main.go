package main

// https://adventofcode.com/2023/day/18

import (
	"aoc/helper"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"sync"
)

const (
	MaxInt int = 20000000
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	digInstructions := ParseDigInstructions(lines)
	solution1 := CountInsideTiles(digInstructions)
	fmt.Println("-> part 1:", solution1)
}

type DigInstruction struct {
	Pos helper.Point2D
	Dir helper.Point2D
	Len int
	RGB string
}

func ParseDigInstructions(lines []string) []DigInstruction {
	pattern := regexp.MustCompile(`^([UDLR]+)\s+(\d+)\s+(\(#[0-9a-f]{6})\)$`)

	nextPos := helper.Point2D{X: 0, Y: 0}

	digInstructions := make([]DigInstruction, len(lines))
	for i := range lines {
		m := pattern.FindStringSubmatch(lines[i])
		if len(m) == 4 {
			length, _ := strconv.Atoi(m[2])
			var dir helper.Point2D
			switch m[1] {
			case "U":
				dir = helper.Point2D{X: 0, Y: -1}
			case "D":
				dir = helper.Point2D{X: 0, Y: 1}
			case "L":
				dir = helper.Point2D{X: -1, Y: 0}
			case "R":
				dir = helper.Point2D{X: 1, Y: 0}
			}
			digInstructions[i] = DigInstruction{
				Pos: nextPos,
				Dir: dir,
				Len: length,
				RGB: m[3],
			}
			nextPos = nextPos.Add(dir.Mul(length))
		}
	}
	return digInstructions
}

func CountInsideTiles(digInstructions []DigInstruction) int {
	min, max := GetBoundary(digInstructions)
	points := make([]helper.Point2D, 0)
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			points = append(points, helper.Point2D{X: x, Y: y})
		}
	}

	var count int
	var m sync.Mutex
	var wg sync.WaitGroup
	for _, p := range points {
		wg.Add(1)
		go func(p helper.Point2D) {
			defer wg.Done()
			if IsInside(p, digInstructions) {
				m.Lock()
				defer m.Unlock()
				count++
			}
		}(p)
	}
	wg.Wait()
	return count
}

func GetBoundary(digInstructions []DigInstruction) (helper.Point2D, helper.Point2D) {
	var min, max helper.Point2D
	for _, di := range digInstructions {
		if di.Pos.X < min.X {
			min.X = di.Pos.X
		}
		if di.Pos.Y < min.Y {
			min.Y = di.Pos.Y
		}
		if di.Pos.X > max.X {
			max.X = di.Pos.X
		}
		if di.Pos.Y > max.Y {
			max.Y = di.Pos.Y
		}
	}
	return min, max
}

func IsInside(p helper.Point2D, digInstructions []DigInstruction) bool {
	wn := ComputeWindingNumber(p, digInstructions)
	return wn <= -0.9 || wn >= 0.9
}

func ComputeWindingNumber(p helper.Point2D, digInstructions []DigInstruction) float64 {
	var windingNumber float64
	for i := 0; i < len(digInstructions); i++ {
		l1 := digInstructions[i].Pos
		l2 := digInstructions[(i+1)%len(digInstructions)].Pos
		w := ComputeWindingNumberOfLine(p, l1, l2)
		windingNumber += w
	}
	return windingNumber
}

func ComputeWindingNumberOfLine(p helper.Point2D, l1, l2 helper.Point2D) float64 {
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
