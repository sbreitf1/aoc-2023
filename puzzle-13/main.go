package main

// https://adventofcode.com/2023/day/13

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadLines("input.txt")

	patterns := ParsePatterns(lines)
	solution1 := SummarizeReflectionsPart1(patterns)
	solution2 := SummarizeReflectionsPart2(patterns)
	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Pattern struct {
	Rows []string
	Cols []string
}

func ParsePatterns(lines []string) []Pattern {
	patterns := make([]Pattern, 0)
	requireNewPattern := true
	for _, line := range lines {
		if len(line) == 0 {
			requireNewPattern = true

		} else {
			if requireNewPattern {
				patterns = append(patterns, Pattern{})
				requireNewPattern = false
			}
			patterns[len(patterns)-1].Rows = append(patterns[len(patterns)-1].Rows, line)
		}
	}
	for i := range patterns {
		patterns[i].computeCols()
	}
	return patterns
}

func (p *Pattern) computeCols() {
	p.Cols = make([]string, len(p.Rows[0]))
	for y := range p.Rows {
		for x, r := range p.Rows[y] {
			p.Cols[x] += string(r)
		}
	}
}

func findReflection(arr []string) int {
	for i := 0; i < len(arr)-1; i++ {
		if countReflectionSmudges(arr, i) == 0 {
			return i
		}
	}
	return -1
}

func findReflectionWithSingleSmudge(arr []string) int {
	for i := 0; i < len(arr)-1; i++ {
		if countReflectionSmudges(arr, i) == 1 {
			return i
		}
	}
	return -1
}

func countReflectionSmudges(arr []string, pos int) int {
	checkSize := pos
	if checkSize > (len(arr) - pos - 2) {
		checkSize = len(arr) - pos - 2
	}
	var count int
	for i := 0; i <= checkSize; i++ {
		for j := range arr[pos-i] {
			if arr[pos-i][j] != arr[pos+i+1][j] {
				count++
			}
		}
	}
	return count
}

func SummarizeReflectionsPart1(patterns []Pattern) int {
	var sum int
	for _, p := range patterns {
		reflectionRow := findReflection(p.Rows)
		reflectionCol := findReflection(p.Cols)
		if reflectionRow >= 0 && reflectionCol >= 0 {
			helper.ExitWithMessage("both row and col reflection detected")
		}
		if reflectionRow < 0 && reflectionCol < 0 {
			helper.ExitWithMessage("no reflection detected")
		}

		if reflectionCol >= 0 {
			sum += reflectionCol + 1
		}
		if reflectionRow >= 0 {
			sum += 100 * (reflectionRow + 1)
		}
	}
	return sum
}

func SummarizeReflectionsPart2(patterns []Pattern) int {
	var sum int
	for _, p := range patterns {
		reflectionRow := findReflectionWithSingleSmudge(p.Rows)
		reflectionCol := findReflectionWithSingleSmudge(p.Cols)
		if reflectionRow >= 0 && reflectionCol >= 0 {
			helper.ExitWithMessage("both row and col reflection detected")
		}
		if reflectionRow < 0 && reflectionCol < 0 {
			helper.ExitWithMessage("no reflection detected")
		}

		if reflectionCol >= 0 {
			sum += reflectionCol + 1
		}
		if reflectionRow >= 0 {
			sum += 100 * (reflectionRow + 1)
		}
	}
	return sum
}
