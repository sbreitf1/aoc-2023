package main

// https://adventofcode.com/2023/day/9

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	sequences := ReadSequences(lines)
	solution1 := SumExtrapolations(sequences)
	solution2 := SumReverseExtrapolations(sequences)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Sequence []int

func ReadSequences(lines []string) []Sequence {
	sequences := make([]Sequence, 0, len(lines))
	for _, line := range lines {
		sequence := helper.SplitAndParseInts(line, " ")
		sequences = append(sequences, sequence)
	}
	return sequences
}

func (s Sequence) Reverse() Sequence {
	rs := make([]int, len(s))
	for i := range s {
		rs[len(s)-i-1] = s[i]
	}
	return rs
}

func (s Sequence) Extrapolate() int {
	if s.IsAllZeroes() {
		return 0
	}

	diff := s.GetDiffSequence()
	extDiff := diff.Extrapolate()
	return s[len(s)-1] + extDiff
}

func (s Sequence) GetDiffSequence() Sequence {
	diff := make([]int, len(s)-1)
	for i := range diff {
		diff[i] = s[i+1] - s[i]
	}
	return diff
}

func (s Sequence) IsAllZeroes() bool {
	for _, num := range s {
		if num != 0 {
			return false
		}
	}
	return true
}

func SumExtrapolations(sequences []Sequence) int {
	sum := 0
	for _, s := range sequences {
		sum += s.Extrapolate()
	}
	return sum
}

func SumReverseExtrapolations(sequences []Sequence) int {
	sum := 0
	for _, s := range sequences {
		sum += s.Reverse().Extrapolate()
	}
	return sum
}
