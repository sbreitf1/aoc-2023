package main

// https://adventofcode.com/2023/day/25

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadNonEmptyLines("example-1.txt")

	network := ParseNetwork(lines)
	fmt.Println(network)
	solution1 := 0
	fmt.Println("-> part 1:", solution1)
}

func ParseNetwork(lines []string) *Network {
	components := make(map[string]*map[string]bool)

	insertLink := func(from, to string) {
		if mFrom, ok := components[from]; ok {
			(*mFrom)[to] = true
		} else {
			components[from] = &map[string]bool{to: true}
		}
	}

	for _, line := range lines {
		parts := helper.SplitAndTrim(line, ":")
		if len(parts) != 2 {
			helper.ExitWithMessage("malformed line %q", line)
		}
		from := parts[0]
		parts = helper.SplitAndTrim(parts[1], " ")
		for _, to := range parts {
			insertLink(from, to)
		}
	}
	return &Network{Components: components}
}

type Network struct {
	Components map[string]*map[string]bool
}
