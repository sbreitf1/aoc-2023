package main

// https://adventofcode.com/2023/day/18

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	digInstructions := ParseDigInstructions(lines)
	solution1 := CountInsideTiles(digInstructions)
	fmt.Println("-> part 1:", solution1)

	digInstructions2 := TransformDigInstructions(digInstructions)
	solution2 := CountInsideTiles(digInstructions2)
	fmt.Println("-> part 2:", solution2)
}

type DigInstruction struct {
	Pos helper.Point2D[int]
	Dir helper.Point2D[int]
	Len int
	RGB string
}

func ParseDigInstructions(lines []string) []DigInstruction {
	pattern := regexp.MustCompile(`^([UDLR]+)\s+(\d+)\s+\(#([0-9a-f]{6})\)$`)

	nextPos := helper.Point2D[int]{X: 0, Y: 0}

	digInstructions := make([]DigInstruction, len(lines))
	for i := range lines {
		m := pattern.FindStringSubmatch(lines[i])
		if len(m) == 4 {
			length, _ := strconv.Atoi(m[2])
			var dir helper.Point2D[int]
			switch m[1] {
			case "U":
				dir = helper.Point2D[int]{X: 0, Y: -1}
			case "D":
				dir = helper.Point2D[int]{X: 0, Y: 1}
			case "L":
				dir = helper.Point2D[int]{X: -1, Y: 0}
			case "R":
				dir = helper.Point2D[int]{X: 1, Y: 0}
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

func CountInsideTiles(digInstructions []DigInstruction) int64 {
	// https://www.101computing.net/the-shoelace-algorithm/
	var count int64
	for i := 0; i < len(digInstructions)-1; i++ {
		count += int64(digInstructions[i].Pos.X) * int64(digInstructions[i+1].Pos.Y)
		count -= int64(digInstructions[i+1].Pos.X) * int64(digInstructions[i].Pos.Y)
	}
	count += int64(digInstructions[len(digInstructions)-1].Pos.X) * int64(digInstructions[0].Pos.Y)
	count -= int64(digInstructions[0].Pos.X) * int64(digInstructions[len(digInstructions)-1].Pos.Y)
	if count < 0 {
		return int64(-count)/2 + CountBoundaryTiles(digInstructions)/2 + 1
	}
	return int64(count)/2 + CountBoundaryTiles(digInstructions)/2 + 1
}

func CountBoundaryTiles(digInstructions []DigInstruction) int64 {
	var count int64
	for _, di := range digInstructions {
		count += int64(di.Len)
	}
	return count
}

func TransformDigInstructions(digInstructions []DigInstruction) []DigInstruction {
	diTransformed := make([]DigInstruction, len(digInstructions))
	nextPos := helper.Point2D[int]{X: 0, Y: 0}
	for i := range digInstructions {
		length, err := strconv.ParseInt(digInstructions[i].RGB[:5], 16, 32)
		helper.ExitOnError(err)
		var dir helper.Point2D[int]
		switch rune(digInstructions[i].RGB[5]) {
		case '3':
			dir = helper.Point2D[int]{X: 0, Y: -1}
		case '1':
			dir = helper.Point2D[int]{X: 0, Y: 1}
		case '2':
			dir = helper.Point2D[int]{X: -1, Y: 0}
		case '0':
			dir = helper.Point2D[int]{X: 1, Y: 0}
		}
		diTransformed[i] = DigInstruction{
			Pos: nextPos,
			Dir: dir,
			Len: int(length),
		}
		nextPos = nextPos.Add(dir.Mul(int(length)))
	}
	return diTransformed
}
