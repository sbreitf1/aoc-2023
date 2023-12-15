package main

// https://adventofcode.com/2023/day/15

import (
	"aoc/helper"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	initSequences := ParseInitSequences(lines)
	solution1 := ComputeSumOfHashes(initSequences)
	fmt.Println("-> part 1:", solution1)

	boxes := ComputeBoxes(initSequences)
	solution2 := ComputePart2(boxes)
	fmt.Println("-> part 2:", solution2)
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

type Lens struct {
	Label       string
	FocalLength int
}

func ComputeBoxes(initSequences []string) map[int][]Lens {
	boxes := make(map[int][]Lens)
	for _, s := range initSequences {
		label, action := ParseAction(s)
		hash := HashString(label)
		lenses := boxes[hash]
		if action == "-" {
			for i := len(lenses) - 1; i >= 0; i-- {
				if lenses[i].Label == label {
					lenses = append(lenses[:i], lenses[i+1:]...)
				}
			}
		} else {
			focalLength, err := strconv.Atoi(action)
			helper.ExitOnError(err)
			found := false
			for i := range lenses {
				if lenses[i].Label == label {
					lenses[i].FocalLength = focalLength
					found = true
				}
			}
			if !found {
				lenses = append(lenses, Lens{Label: label, FocalLength: focalLength})
			}
		}
		boxes[hash] = lenses
	}
	return boxes
}

func ParseAction(str string) (string, string) {
	if strings.HasSuffix(str, "-") {
		return str[:len(str)-1], "-"
	}
	parts := strings.Split(str, "=")
	if len(parts) > 3 {
		helper.ExitWithMessage("invalid sequence %q", str)
	}
	return parts[0], parts[1]
}

func ComputePart2(boxes map[int][]Lens) int {
	var sum int
	for hash, lenses := range boxes {
		for i, l := range lenses {
			sum += (hash + 1) * (i + 1) * l.FocalLength
		}
	}
	return sum
}
