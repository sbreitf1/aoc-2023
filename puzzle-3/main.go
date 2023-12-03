package main

// https://adventofcode.com/2023/day/2

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	partNumbers := ExtractPartNumbers(lines)
	solution1 := SumPartNumbers(partNumbers)

	fmt.Println("-> part 1:", solution1)
}

func ExtractPartNumbers(lines []string) []int {
	lines = ExtendEmptyBorder(lines)
	partNumbers := make([]int, 0)
	for i := 1; i < (len(lines) - 1); i++ {
		if len(lines[i]) > 0 {
			linePartNumbers := ScanPartNumbers(lines, i)
			partNumbers = append(partNumbers, linePartNumbers...)
		}
	}
	return partNumbers
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

func ScanPartNumbers(lines []string, lineIndex int) []int {
	partNumbers := make([]int, 0)
	for i := 1; i < (len(lines[lineIndex]) - 1); i++ {
		pos, numStr, ok := ScanNextNumberInLine(lines[lineIndex], i)
		if !ok {
			break
		}

		if IsPartNumber(lines, lineIndex, pos, numStr) {
			num, _ := strconv.Atoi(numStr)
			partNumbers = append(partNumbers, num)
		}

		i = pos + len(numStr)
	}
	return partNumbers
}

var patternNum = regexp.MustCompile(`\d+`)

func ScanNextNumberInLine(line string, pos int) (int, string, bool) {
	m := patternNum.FindStringIndex(line[pos:])
	if len(m) != 2 {
		return 0, "", false
	}

	return pos + m[0], line[pos+m[0] : pos+m[1]], true
}

func IsPartNumber(lines []string, lineIndex int, pos int, numStr string) bool {
	line1 := lines[lineIndex-1][pos-1 : pos+len(numStr)+1]
	line2 := lines[lineIndex][pos-1 : pos+len(numStr)+1]
	line3 := lines[lineIndex+1][pos-1 : pos+len(numStr)+1]
	allSurroundingChars := line1 + string(line2[0]) + string(line2[len(line2)-1]) + line3
	return allSurroundingChars != strings.Repeat(".", 2*(len(numStr)+2)+2)
}

func SumPartNumbers(partNumbers []int) int {
	var sum int
	for _, num := range partNumbers {
		sum += num
	}
	return sum
}
