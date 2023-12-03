package main

// https://adventofcode.com/2023/day/1

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	calibrationValues1 := getCalibrationValues(lines, map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9})
	solution1 := sumAllCalibrationValues(calibrationValues1)

	calibrationValues2 := getCalibrationValues(lines, map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	})
	solution2 := sumAllCalibrationValues(calibrationValues2)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

func getCalibrationValues(lines []string, tokenMappings map[string]int) []int {
	calibrationValues := make([]int, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			tokens := tokenize(line, tokenMappings)
			if len(tokens) > 0 {
				val := 10*tokens[0] + tokens[len(tokens)-1]
				calibrationValues = append(calibrationValues, val)
			}
		}
	}
	return calibrationValues
}

func tokenize(line string, tokenMappings map[string]int) []int {
	runes := []rune(line)
	tokens := make([]int, 0)
	for i := 0; i < len(runes); i++ {
		remainder := string(runes[i:])
		for t, v := range tokenMappings {
			if strings.HasPrefix(remainder, t) {
				tokens = append(tokens, v)
				break
			}
		}
	}
	return tokens
}

func sumAllCalibrationValues(calibrationValues []int) int {
	var sum int
	for _, val := range calibrationValues {
		sum += val
	}
	return sum
}
