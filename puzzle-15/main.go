package main

// https://adventofcode.com/2023/day/15

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	initSequences := ParseInitSequences(lines)
	solution1 := ComputeSumOfHashes(initSequences)

	fmt.Println("-> part 1:", solution1)
}

func ParseInitSequences(lines []string) []string {
	initSequences := make([]string, 0)
	for _, line := range lines {
		initSequences = append(initSequences, strings.Split(line, ",")...)
	}
	return initSequences
}

func ComputeSumOfHashes(initSequences []string) int {
	var sum int
	for _, s := range initSequences {
		sum += HashString(s)
	}
	return sum
}

func HashString(str string) int {
	var hash int
	for _, r := range str {
		hash += int(r)
		hash = (hash * 17) % 256
	}
	return hash
}
