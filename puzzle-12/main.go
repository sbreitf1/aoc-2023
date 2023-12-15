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
	if solution1 != 7163 {
		helper.ExitWithMessage("solution 1 is wrong")
	}

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
	var damageSumToPlace int
	for i := range g.DamagedGroups {
		damageSumToPlace += g.DamagedGroups[i]
	}
	okSumToPlace := len(g.DamagedGroups)
	str := g.HotSprings + "."
	positions := make([]Pos, len(str))
	for i := len(str) - 1; i >= 0; i-- {
		positions[i].Rune = rune(str[i])
		if i < (len(str) - 1) {
			positions[i].FollowingDamageCount = positions[i+1].FollowingDamageCount
			positions[i].FollowingOKCount = positions[i+1].FollowingOKCount
			positions[i].FollowingAnyCount = positions[i+1].FollowingAnyCount
		}
		if str[i] == '#' {
			positions[i].FollowingDamageCount++
		}
		if str[i] == '.' {
			positions[i].FollowingOKCount++
		}
		if str[i] == '?' {
			positions[i].FollowingAnyCount++
		}
	}
	return findArrangementsRecursive(positions, g.DamagedGroups, damageSumToPlace, okSumToPlace)
}

type Pos struct {
	Rune                 rune
	FollowingAnyCount    int
	FollowingOKCount     int
	FollowingDamageCount int
}

func findArrangementsRecursive(positions []Pos, dgPlace []int, damageSumToPlace, okSumToPlace int) uint64 {
	if len(dgPlace) == 0 {
		if len(positions) == 0 || positions[0].FollowingDamageCount == 0 {
			return 1
		} else {
			return 0
		}
	}
	if len(positions) == 0 {
		return 0
	}
	if (positions[0].FollowingDamageCount + positions[0].FollowingAnyCount) < damageSumToPlace {
		return 0
	}
	if (positions[0].FollowingOKCount + positions[0].FollowingAnyCount) < okSumToPlace {
		return 0
	}
	if (positions[0].FollowingOKCount + positions[0].FollowingDamageCount + positions[0].FollowingAnyCount) < (okSumToPlace + damageSumToPlace) {
		return 0
	}

	var result uint64
	for i := 0; i < len(positions)-(damageSumToPlace+okSumToPlace)+1; i++ {
		// check placement possible
		canPlace := true
		for j := 0; j < dgPlace[0]; j++ {
			if positions[i+j].Rune == '.' {
				canPlace = false
				break
			}
		}
		if canPlace && positions[i+dgPlace[0]].Rune != '#' {
			result += findArrangementsRecursive(positions[i+dgPlace[0]+1:], dgPlace[1:], damageSumToPlace-dgPlace[0], okSumToPlace-1)
		}

		if positions[i].Rune == '#' {
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
