package main

// https://adventofcode.com/2023/day/6

import (
	"aoc/helper"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	races := ParseRaces(lines)
	solution1 := ComputeSolution(races)
	race2 := AccountForBadKerning(races)
	solution2 := ComputeSolution([]Race{race2})

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Race struct {
	Time     int
	Distance int
}

func ParseRaces(lines []string) []Race {
	times := ParseInts(lines[0])
	distances := ParseInts(lines[1])
	if len(times) != len(distances) {
		helper.ExitWithMessage("mismatching times and distances count")
	}
	races := make([]Race, len(times))
	for i := range times {
		races[i].Time = times[i]
		races[i].Distance = distances[i]
	}
	return races
}

func ParseInts(line string) []int {
	pos := strings.IndexRune(line, ':')
	line = line[pos+1:]
	parts := strings.Split(line, " ")
	ints := make([]int, 0)
	for _, p := range parts {
		if len(p) > 0 {
			num, err := strconv.Atoi(p)
			helper.ExitOnError(err, "invalid int value %q", p)
			ints = append(ints, num)
		}
	}
	return ints
}

func (r Race) Simulate(holdTime int) int {
	if holdTime > r.Time {
		helper.ExitWithMessage("cannot simulate a holdTime longer than the actual race time")
	}
	dv := holdTime
	moveTime := r.Time - holdTime
	return dv * moveTime
}

func (r Race) NumberOfWinningHoldTimes() int {
	count := 0
	for i := 0; i <= r.Time; i++ {
		if r.Simulate(i) > r.Distance {
			count++
		}
	}
	return count
}

func ComputeSolution(races []Race) int {
	product := 1
	for _, r := range races {
		v := r.NumberOfWinningHoldTimes()
		product *= v
	}
	return product
}

func AccountForBadKerning(races []Race) Race {
	var timeStr, distanceStr string
	for _, r := range races {
		timeStr += strconv.Itoa(r.Time)
		distanceStr += strconv.Itoa(r.Distance)
	}
	var race Race
	race.Time, _ = strconv.Atoi(timeStr)
	race.Distance, _ = strconv.Atoi(distanceStr)
	return race
}
