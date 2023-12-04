package main

// https://adventofcode.com/2023/day/3

import (
	"aoc/helper"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	partNumbers := ExtractPartNumbers(lines)
	solution1 := SumPartNumbers(partNumbers)

	gears := FindGears(lines, partNumbers)
	solution2 := SumGears(gears)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type PartNumber struct {
	Val   int
	X, Y  int
	Width int
}

func ExtractPartNumbers(lines []string) []PartNumber {
	linesWithBorder := ExtendEmptyBorder(lines)
	numbers := FindNumbers(linesWithBorder)
	partNumbers := FilterPartNumbers(linesWithBorder, numbers)
	return partNumbers
}

func FindNumbers(linesWithBorder []string) []PartNumber {
	numbers := make([]PartNumber, 0)
	for i := 1; i < (len(linesWithBorder) - 1); i++ {
		numbersInLine := FindNumbersInLine(linesWithBorder, i)
		numbers = append(numbers, numbersInLine...)
	}
	return numbers
}

func FindNumbersInLine(linesWithBorder []string, lineIndex int) []PartNumber {
	numbersInLine := make([]PartNumber, 0)
	var currentNumberStr string
	for i := 1; i < len(linesWithBorder[lineIndex]); i++ {
		r := linesWithBorder[lineIndex][i]
		if r >= '0' && r <= '9' {
			currentNumberStr += string(r)
		} else if len(currentNumberStr) > 0 {
			val, _ := strconv.Atoi(currentNumberStr)
			numbersInLine = append(numbersInLine, PartNumber{
				Val:   val,
				X:     i - len(currentNumberStr) - 1,
				Y:     lineIndex - 1,
				Width: len(currentNumberStr),
			})
			currentNumberStr = ""
		}
	}
	return numbersInLine
}

func FilterPartNumbers(linesWithBorder []string, numbers []PartNumber) []PartNumber {
	partNumbers := make([]PartNumber, 0)
	for _, num := range numbers {
		if IsPartNumber(linesWithBorder, num) {
			partNumbers = append(partNumbers, num)
		}
	}
	return partNumbers
}

func IsPartNumber(linesWithBorder []string, num PartNumber) bool {
	for x := (num.X - 1); x < (num.X + num.Width + 1); x++ {
		if linesWithBorder[num.Y][x+1] != '.' {
			return true
		}
		if linesWithBorder[num.Y+2][x+1] != '.' {
			return true
		}
	}
	return linesWithBorder[num.Y+1][num.X] != '.' || linesWithBorder[num.Y+1][num.X+num.Width+1] != '.'
}

func SumPartNumbers(partNumbers []PartNumber) int {
	var sum int
	for _, num := range partNumbers {
		sum += num.Val
	}
	return sum
}

func FindGears(lines []string, partNumbers []PartNumber) map[string][]int {
	gears := make(map[string][]int, 0)
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			if lines[y][x] == '*' {
				gear := FindGearNumbers(partNumbers, x, y)
				if len(gear) >= 2 {
					gears[fmt.Sprintf("%d,%d", x, y)] = gear
				}
			}
		}
	}
	return gears
}

func FindGearNumbers(partNumbers []PartNumber, x, y int) []int {
	gear := make([]int, 0)
	for _, num := range partNumbers {
		if BoxesOverlap(x-1, y-1, 3, 3, num.X, num.Y, num.Width, 1) {
			gear = append(gear, num.Val)
		}
	}
	return gear
}

func BoxesOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	if x1+w1-1 < x2 || x1 > x2+w2-1 {
		return false
	}
	if y1+h1-1 < y2 || y1 > y2+h2-1 {
		return false
	}
	return true
}

func SumGears(gears map[string][]int) int {
	var sum int
	for _, g := range gears {
		gearRatio := 1
		for _, v := range g {
			gearRatio *= v
		}
		sum += gearRatio
	}
	return sum
}

func ExtendEmptyBorder(lines []string) []string {
	linesWithBorder := make([]string, 0, len(lines)+2)
	linesWithBorder = append(linesWithBorder, strings.Repeat(".", len(lines[0])+2))
	for _, l := range lines {
		linesWithBorder = append(linesWithBorder, "."+l+".")
	}
	linesWithBorder = append(linesWithBorder, strings.Repeat(".", len(lines[0])+2))
	return linesWithBorder
}
