package main

// https://adventofcode.com/2022/day/1

import (
	"aoc/helper"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	elfCounts := getElfCounts(lines)

	sort.Ints(elfCounts)

	fmt.Println("-> part 1:", elfCounts[len(elfCounts)-1])
	fmt.Println("-> part 2:", elfCounts[len(elfCounts)-1]+elfCounts[len(elfCounts)-2]+elfCounts[len(elfCounts)-3])
}

func getElfCounts(lines []string) []int {
	elfCounts := make([]int, 1)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			elfCounts = append(elfCounts, 0)
		} else {
			d, err := strconv.Atoi(l)
			helper.ExitOnError(err)
			elfCounts[len(elfCounts)-1] += d
		}
	}
	return elfCounts
}
