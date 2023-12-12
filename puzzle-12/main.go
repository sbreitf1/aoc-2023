package main

// https://adventofcode.com/2023/day/12

import (
	"aoc/helper"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	//cacheDisabled = true
	file := "input.txt"
	lines := helper.ReadNonEmptyLines(file)

	groups := ParseHotSpringGroups(lines)
	solution1 := CountArrangements(file+"_part1", groups)
	fmt.Println("-> part 1:", solution1)

	unfoldedGroups := UnfoldGroups(groups, 5)
	solution2 := CountArrangements(file+"_part2", unfoldedGroups)
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

func (g HotSpringGroup) CountArrangements() int {
	var minRemainingStringLen int
	for i := range g.DamagedGroups {
		minRemainingStringLen += g.DamagedGroups[i]
	}
	minRemainingStringLen += len(g.DamagedGroups) - 1
	return g.findArrangementsRecursive(g.HotSprings, 0, -1, 0, minRemainingStringLen)
}

func (g HotSpringGroup) findArrangementsRecursive(str string, i int, dgPos, damageCount int, minRemainingStringLen int) int {
	if dgPos >= len(g.DamagedGroups) {
		// not a solution, more damaged groups than expected
		return 0
	}
	if dgPos >= 0 && damageCount > g.DamagedGroups[dgPos] {
		// not a solution, current damage group does not match expected one
		return 0
	}

	if len(str)-i < minRemainingStringLen {
		// not enough chars in string remaining to fulfill all damage groups
		return 0
	}

	if i >= len(g.HotSprings) {
		if dgPos == (len(g.DamagedGroups)-1) && damageCount == g.DamagedGroups[dgPos] {
			return 1
		}
		return 0
	}
	if str[i] == '.' {
		if dgPos >= 0 && str[i-1] == '#' && damageCount != g.DamagedGroups[dgPos] {
			return 0
		}
		return g.findArrangementsRecursive(str[:i]+"."+str[i+1:], i+1, dgPos, damageCount, minRemainingStringLen)
	}
	if str[i] == '#' {
		if i == 0 {
			return g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos+1, 1, minRemainingStringLen-2)
		} else if str[i-1] == '.' {
			return g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos+1, 1, minRemainingStringLen-2)
		} else {
			return g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos, damageCount+1, minRemainingStringLen-1)
		}
	}

	result := 0

	if dgPos < 0 || (str[i-1] == '.') || (str[i-1] == '#' && damageCount == g.DamagedGroups[dgPos]) {
		result += g.findArrangementsRecursive(str[:i]+"."+str[i+1:], i+1, dgPos, damageCount, minRemainingStringLen)
	}

	if i == 0 {
		result += g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos+1, 1, minRemainingStringLen-2)
	} else if str[i-1] == '.' {
		result += g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos+1, 1, minRemainingStringLen-2)
	} else {
		result += g.findArrangementsRecursive(str[:i]+"#"+str[i+1:], i+1, dgPos, damageCount+1, minRemainingStringLen-1)
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

func CountArrangements(id string, groups []HotSpringGroup) int {
	var count, doneCount int32
	var m sync.Mutex
	var wg sync.WaitGroup
	for i := range groups {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var arragementCount int
			if val, ok := ReadFromCache(id, i); ok {
				arragementCount = val
			} else {
				arragementCount = groups[i].CountArrangements()
				SetCache(id, i, arragementCount)
			}

			m.Lock()
			defer m.Unlock()

			atomic.AddInt32(&count, int32(arragementCount))
			atomic.AddInt32(&doneCount, 1)
			fmt.Println(len(groups)-int(doneCount), "remaining")
		}(i)
	}
	wg.Wait()
	return int(count)
}

var cacheMutex sync.RWMutex
var cacheDisabled bool

func ReadFromCache(id string, num int) (int, bool) {
	if cacheDisabled {
		return 0, false
	}

	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	data, err := os.ReadFile("cache.json")
	if err != nil {
		if !os.IsNotExist(err) {
			helper.ExitOnError(err)
		}
		data = []byte("{}")
	}

	var cache map[string]int
	helper.ExitOnError(json.Unmarshal(data, &cache))

	key := fmt.Sprintf("%s[%d]", id, num)
	val, ok := cache[key]
	return val, ok
}

func SetCache(id string, num, val int) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	data, err := os.ReadFile("cache.json")
	if err != nil {
		if !os.IsNotExist(err) {
			helper.ExitOnError(err)
		}
		data = []byte("{}")
	}

	var cache map[string]int
	helper.ExitOnError(json.Unmarshal(data, &cache))

	key := fmt.Sprintf("%s[%d]", id, num)
	cache[key] = val

	data, err = json.MarshalIndent(&cache, "", "  ")
	helper.ExitOnError(err)

	helper.ExitOnError(os.WriteFile("cache.json", data, os.ModePerm))
}
