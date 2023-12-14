package main

// https://adventofcode.com/2023/day/12

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	groups := ParseHotSpringGroups(lines)
	solution1 := CountArrangements(groups)
	fmt.Println("-> part 1:", solution1)

	unfoldedGroups := UnfoldGroups(groups, 5)
	solution2 := CountArrangements(unfoldedGroups)
	fmt.Println("-> part 2:", solution2)
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

func (g HotSpringGroup) CountArrangements() uint64 {
	var minRequiredLen int
	for i := range g.DamagedGroups {
		minRequiredLen += g.DamagedGroups[i]
	}
	minRequiredLen += len(g.DamagedGroups)
	str := g.HotSprings + "."
	followingDamageCounts := make([]int, len(str))
	followingAnyCounts := make([]int, len(str))
	for i := len(str) - 2; i >= 0; i-- {
		if str[i] == '#' {
			followingDamageCounts[i] = followingDamageCounts[i+1] + 1
		} else {
			followingDamageCounts[i] = followingDamageCounts[i+1]
		}
		if str[i] == '?' {
			followingAnyCounts[i] = followingAnyCounts[i+1] + 1
		} else {
			followingAnyCounts[i] = followingAnyCounts[i+1]
		}
	}
	return findArrangementsRecursive(str, followingDamageCounts, followingAnyCounts, g.DamagedGroups, minRequiredLen)
}

func findArrangementsRecursive(str string, followingDamageCounts, followingAnyCounts []int, dgPlace []int, minRequiredLen int) uint64 {
	if len(dgPlace) == 0 {
		if len(followingDamageCounts) == 0 || followingDamageCounts[0] == 0 {
			return 1
		} else {
			return 0
		}
	}
	if len(str) < minRequiredLen {
		return 0
	}

	var result uint64
	for i := 0; i < len(str)-minRequiredLen+1; i++ {
		// check placement possible
		canPlace := true
		for j := 0; j < dgPlace[0]; j++ {
			if str[i+j] == '.' {
				canPlace = false
				break
			}
		}
		if canPlace && str[i+dgPlace[0]] != '#' {
			result += findArrangementsRecursive(str[i+dgPlace[0]+1:], followingDamageCounts[i+dgPlace[0]+1:], followingAnyCounts[i+dgPlace[0]+1:], dgPlace[1:], minRequiredLen-dgPlace[0]-1)
		}

		if str[i] == '#' {
			break
		}
	}
	return result
}

func UnfoldGroups(groups []HotSpringGroup, count int) []HotSpringGroup {
	unfoldedGroups := make([]HotSpringGroup, len(groups))
	for i := range groups {
		unfoldedGroups[i] = groups[i].Unfold(count)
	}
	return unfoldedGroups
}

func (g HotSpringGroup) Unfold(count int) HotSpringGroup {
	damagedGroups := make([]int, count*len(g.DamagedGroups))
	for i := 0; i < count; i++ {
		copy(damagedGroups[i*len(g.DamagedGroups):(i+1)*len(g.DamagedGroups)], g.DamagedGroups)
	}
	return HotSpringGroup{
		HotSprings:    strings.Repeat("?"+g.HotSprings, count)[1:],
		DamagedGroups: damagedGroups,
	}
}

func CountArrangements(groups []HotSpringGroup) uint64 {
	var count uint64
	var doneCount int32
	var m sync.Mutex
	var wg sync.WaitGroup
	for i := range groups {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			arragementCount := groups[i].CountArrangements()

			m.Lock()
			defer m.Unlock()

			atomic.AddUint64(&count, arragementCount)
			atomic.AddInt32(&doneCount, 1)
			fmt.Println(len(groups)-int(doneCount), "remaining")
		}(i)
	}
	wg.Wait()
	return count
}
