package main

// https://adventofcode.com/2023/day/12

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	groups := ParseHotSpringGroups(lines)
	solution1 := CountArrangements(groups)

	fmt.Println("-> part 1:", solution1)
}

type HotSpringGroup struct {
	HotSprings    string
	DamagedGroups []int
}

func ParseHotSpringGroups(lines []string) []HotSpringGroup {
	pattern := regexp.MustCompile(`^\s*([.#?]+)\s+([0-9,]+)\s*$`)
	groups := make([]HotSpringGroup, 0, len(lines))
	for i, line := range lines {
		m := pattern.FindStringSubmatch(line)
		if len(m) != 3 {
			fmt.Println("line", i+1, "did not match")
			continue
		}

		group := HotSpringGroup{
			HotSprings:    m[1],
			DamagedGroups: helper.SplitAndParseInts(m[2], ","),
		}
		groups = append(groups, group)
	}
	return groups
}

func (g HotSpringGroup) FindArrangements() []string {
	return g.FindArrangementsRecursive(g.HotSprings, 0)
}

func (g HotSpringGroup) FindArrangementsRecursive(str string, i int) []string {
	if i >= len(g.HotSprings) {
		dg := GetDamagedGroups(str)
		if DamagedGroupsAreEqual(dg, g.DamagedGroups) {
			return []string{str}
		}
		return nil
	}
	if str[i] == '?' {
		result := make([]string, 0)
		result = append(result, g.FindArrangementsRecursive(str[:i]+"."+str[i+1:], i+1)...)
		result = append(result, g.FindArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1)...)
		return result
	}
	return g.FindArrangementsRecursive(str, i+1)
}

func GetDamagedGroups(str string) []int {
	parts := strings.Split(str, ".")
	dg := make([]int, 0)
	for _, p := range parts {
		if len(p) > 0 {
			dg = append(dg, len(p))
		}
	}
	return dg
}

func DamagedGroupsAreEqual(dg1, dg2 []int) bool {
	if len(dg1) != len(dg2) {
		return false
	}
	for i := range dg1 {
		if dg1[i] != dg2[i] {
			return false
		}
	}
	return true
}

func CountArrangements(groups []HotSpringGroup) int {
	count := 0
	for _, g := range groups {
		arrangements := g.FindArrangements()
		count += len(arrangements)
	}
	return count
}
