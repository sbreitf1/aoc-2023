package main

// https://adventofcode.com/2023/day/4

import (
	"aoc/helper"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	scratchCards := ParseScratchCards(lines)
	solution1 := SumScratchCardPoints(scratchCards)

	EvalScratchCardRules(scratchCards)
	solution2 := CountScratchCards(scratchCards)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type ScratchCard struct {
	Count          int
	WinningNumbers map[int]struct{}
	HavingNumbers  []int
}

var patternScratchCard = regexp.MustCompile(`^Card\s+(\d+):\s*([\d\s]+)\s*\|\s*([\d\s]+)\s*$`)

func ParseScratchCards(lines []string) []ScratchCard {
	scratchCards := make([]ScratchCard, 0)
	for _, line := range lines {
		m := patternScratchCard.FindStringSubmatch(line)
		if len(m) == 4 {
			winningNumbers := ParseSpaceSeparatedInts(m[2])
			winningNumbersMap := make(map[int]struct{}, len(winningNumbers))
			for _, n := range winningNumbers {
				winningNumbersMap[n] = struct{}{}
			}
			havingNumbers := ParseSpaceSeparatedInts(m[3])
			scratchCards = append(scratchCards, ScratchCard{
				Count:          1,
				WinningNumbers: winningNumbersMap,
				HavingNumbers:  havingNumbers,
			})
		}
	}
	return scratchCards
}

func ParseSpaceSeparatedInts(str string) []int {
	parts := strings.Split(str, " ")
	ints := make([]int, 0, len(parts))
	for _, p := range parts {
		if len(p) > 0 {
			n, _ := strconv.Atoi(p)
			ints = append(ints, n)
		}
	}
	return ints
}

func (sc ScratchCard) Points() int {
	matchCount := sc.MatchCount()
	if matchCount == 0 {
		return 0
	}
	return int(math.Pow(2, float64(matchCount-1)))
}

func (sc ScratchCard) MatchCount() int {
	count := 0
	for _, v := range sc.HavingNumbers {
		if _, ok := sc.WinningNumbers[v]; ok {
			count++
		}
	}
	return count
}

func SumScratchCardPoints(scratchCards []ScratchCard) int {
	sum := 0
	for _, sc := range scratchCards {
		sum += sc.Points()
	}
	return sum
}

func EvalScratchCardRules(scratchCards []ScratchCard) {
	for i := range scratchCards {
		matchCount := scratchCards[i].MatchCount()
		for j := 0; j < matchCount; j++ {
			index := i + 1 + j
			if index < len(scratchCards) {
				scratchCards[index].Count += scratchCards[i].Count
			}
		}
	}
}

func CountScratchCards(scratchCards []ScratchCard) int {
	sum := 0
	for _, sc := range scratchCards {
		sum += sc.Count
	}
	return sum
}
