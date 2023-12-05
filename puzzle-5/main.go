package main

// https://adventofcode.com/2023/day/5

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	seeds, mapChain := ParseInput(lines)
	locations := mapChain.MapManySeedsToLocations(seeds)
	solution1 := GetLowestValue(locations)

	fmt.Println("-> part 1:", solution1)
}

type MapChain struct {
	MappingGroups []MappingGroup
}

type MappingGroup struct {
	SrcName, DstName string
	Mappings         []Mapping
}

type Mapping struct {
	SrcStart, DstStart int
	Range              int
}

func ParseInput(lines []string) ([]int, MapChain) {
	seeds := ParseSeeds(lines)
	mapChain := ParseMapChain(lines)
	return seeds, mapChain
}

func ParseSeeds(lines []string) []int {
	if !strings.HasPrefix(lines[0], "seeds:") {
		helper.ExitWithMessage("no seeds in first line")
	}

	parts := strings.Split(lines[0][6:], " ")
	seeds := make([]int, 0)
	for _, p := range parts {
		if len(p) > 0 {
			seed, err := strconv.Atoi(p)
			helper.ExitOnError(err, "invalid seed value %q", p)
			seeds = append(seeds, seed)
		}
	}
	return seeds
}

var patternMappingHeader = regexp.MustCompile(`^(.*)-to-(.*)\s+map:$`)

func ParseMapChain(lines []string) MapChain {
	mappingGroups := make([]MappingGroup, 0)
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) > 0 {
			if m := patternMappingHeader.FindStringSubmatch(lines[i]); len(m) == 3 {
				mappingGroups = append(mappingGroups, MappingGroup{
					SrcName: m[1],
					DstName: m[2],
				})
			} else {
				parts := strings.Split(lines[i], " ")
				if len(parts) != 3 {
					helper.ExitWithMessage("invalid range mapping %q", lines[i])
				}
				var m Mapping
				m.DstStart, _ = strconv.Atoi(parts[0])
				m.SrcStart, _ = strconv.Atoi(parts[1])
				m.Range, _ = strconv.Atoi(parts[2])
				mappingGroups[len(mappingGroups)-1].Mappings = append(mappingGroups[len(mappingGroups)-1].Mappings, m)
			}
		}
	}
	return MapChain{MappingGroups: mappingGroups}
}

func (mc MapChain) MapManySeedsToLocations(seeds []int) []int {
	locations := make([]int, 0, len(seeds))
	for _, seed := range seeds {
		location := mc.MapSeedToLocation(seed)
		locations = append(locations, location)
	}
	return locations
}

func (mc MapChain) MapSeedToLocation(seed int) int {
	val := seed
	for _, m := range mc.MappingGroups {
		val = m.MapValue(val)
	}
	return val
}

func (mg MappingGroup) MapValue(src int) int {
	for _, m := range mg.Mappings {
		if src >= m.SrcStart && src < m.SrcStart+m.Range {
			return src - m.SrcStart + m.DstStart
		}
	}
	return src
}

func GetLowestValue(values []int) int {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}
