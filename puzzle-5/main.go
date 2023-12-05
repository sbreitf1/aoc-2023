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

	seedRanges2, mapChain := ParseInput(lines)
	seedRanges1 := GetSeedRangesPart1(seedRanges2)
	locationRanges1 := mapChain.MapSeedRangesToLocationRanges(seedRanges1)
	locationRanges2 := mapChain.MapSeedRangesToLocationRanges(seedRanges2)
	solution1 := GetLowestRangeValue(locationRanges1)
	solution2 := GetLowestRangeValue(locationRanges2)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Range struct {
	First, Last int
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

func ParseInput(lines []string) ([]Range, MapChain) {
	seedRanges := ParseSeedRanges(lines)
	mapChain := ParseMapChain(lines)
	return seedRanges, mapChain
}

func ParseSeedRanges(lines []string) []Range {
	if !strings.HasPrefix(lines[0], "seeds:") {
		helper.ExitWithMessage("no seeds in first line")
	}

	parts := strings.Split(lines[0][6:], " ")
	ints := make([]int, 0)
	for _, p := range parts {
		if len(p) > 0 {
			seed, err := strconv.Atoi(p)
			helper.ExitOnError(err, "invalid seed value %q", p)
			ints = append(ints, seed)
		}
	}
	seedRanges := make([]Range, 0, len(ints)/2)
	for i := 0; i < len(ints); i += 2 {
		seedRanges = append(seedRanges, Range{
			First: ints[i],
			Last:  ints[i] + ints[i+1] - 1,
		})
	}
	return seedRanges
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

func (mc MapChain) MapSeedRangesToLocationRanges(seedRanges []Range) []Range {
	ranges := seedRanges
	for _, mg := range mc.MappingGroups {
		ranges = mg.MapRanges(ranges)
	}
	return ranges
}

func (mg MappingGroup) MapRanges(srcRanges []Range) []Range {
	ranges := srcRanges
	dstRanges := make([]Range, 0)
	for _, m := range mg.Mappings {
		if len(ranges) == 0 {
			break
		}

		mapped, remainder := m.MapRanges(ranges)
		dstRanges = append(dstRanges, mapped...)
		ranges = remainder
	}
	dstRanges = append(dstRanges, ranges...)
	return dstRanges
}

func (m Mapping) MapRanges(srcRanges []Range) (dstRanges []Range, remainderRanges []Range) {
	for _, src := range srcRanges {
		d, r := m.MapRange(src)
		dstRanges = append(dstRanges, d...)
		remainderRanges = append(remainderRanges, r...)
	}
	return
}

func (m Mapping) MapRange(srcRange Range) (dstRanges []Range, remainderRanges []Range) {
	if srcRange.First > (m.SrcStart + m.Range - 1) {
		return nil, []Range{srcRange}
	}
	if srcRange.Last < m.SrcStart {
		return nil, []Range{srcRange}
	}
	if srcRange.First < m.SrcStart {
		remainderRanges = append(remainderRanges, Range{First: srcRange.First, Last: m.SrcStart - 1})
		srcRange.First = m.SrcStart
	}
	if srcRange.Last > m.SrcStart+m.Range-1 {
		remainderRanges = append(remainderRanges, Range{First: m.SrcStart + m.Range, Last: srcRange.Last})
		srcRange.Last = m.SrcStart + m.Range - 1
	}
	return []Range{{
		First: srcRange.First - m.SrcStart + m.DstStart,
		Last:  srcRange.Last - m.SrcStart + m.DstStart,
	}}, remainderRanges
}

func GetSeedRangesPart1(seedRanges []Range) []Range {
	seedRanges1 := make([]Range, 0, 2*len(seedRanges))
	for _, r := range seedRanges {
		seedRanges1 = append(seedRanges1, Range{First: r.First, Last: r.First})
		seedRanges1 = append(seedRanges1, Range{First: r.Last - r.First + 1, Last: r.Last - r.First + 1})
	}
	return seedRanges1
}

func GetLowestRangeValue(ranges []Range) int {
	min := ranges[0].First
	for _, r := range ranges {
		if r.First < min {
			min = r.First
		}
	}
	return min
}
