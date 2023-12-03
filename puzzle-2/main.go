package main

// https://adventofcode.com/2023/day/2

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadLines("input.txt")

	games := parseGames(lines)
	solution1 := sumPossibleGameIDs(games, Bag{
		ColorRed:   12,
		ColorGreen: 13,
		ColorBlue:  14,
	})
	solution2 := sumPowerOfMinBags(games)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Game struct {
	ID   int
	Sets []Set
}

type Set map[Color]int
type Bag map[Color]int

type Color string

var ColorRed Color = "red"
var ColorGreen Color = "green"
var ColorBlue Color = "blue"

func parseGames(lines []string) []Game {
	games := make([]Game, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			game := parseGame(line)
			games = append(games, game)
		}
	}
	return games
}

var patternGameHeader = regexp.MustCompile(`^Game\s+(\d+):(.*)$`)

func parseGame(line string) Game {
	m := patternGameHeader.FindStringSubmatch(line)
	if len(m) != 3 {
		helper.ExitWithMessage("game string %q is malformed", line)
	}

	id, _ := strconv.Atoi(m[1])
	setStrings := strings.Split(m[2], ";")
	sets := make([]Set, 0, len(setStrings))
	for _, str := range setStrings {
		set := parseSet(str)
		sets = append(sets, set)
	}
	return Game{ID: id, Sets: sets}
}

var patternCube = regexp.MustCompile(`(\d+)\s+(red|green|blue)`)

func parseSet(str string) Set {
	matches := patternCube.FindAllStringSubmatch(str, -1)
	set := make(Set)
	for _, m := range matches {
		num, _ := strconv.Atoi(m[1])
		color := Color(m[2])
		if _, ok := set[color]; ok {
			helper.ExitWithMessage("color %q appeared twice in set", color)
		}
		set[color] = num
	}
	if len(set) == 0 {
		helper.ExitWithMessage("empty set %q", str)
	}
	return set
}

func sumPossibleGameIDs(games []Game, bag Bag) int {
	var sum int
	for _, g := range games {
		if g.IsPossible(bag) {
			sum += g.ID
		}
	}
	return sum
}

func (game Game) IsPossible(bag Bag) bool {
	for _, set := range game.Sets {
		if !set.IsPossible(bag) {
			return false
		}
	}
	return true
}

func (game Game) GetMinBag() Bag {
	minBag := make(Bag)
	for _, set := range game.Sets {
		for c, v := range set {
			currentVal := minBag[c]
			if currentVal < v {
				minBag[c] = v
			}
		}
	}
	return minBag
}

func (set Set) IsPossible(bag Bag) bool {
	for c, v := range set {
		vMax := bag[c]
		if v > vMax {
			return false
		}
	}
	return true
}

func sumPowerOfMinBags(games []Game) int {
	var sum int
	for _, g := range games {
		minBag := g.GetMinBag()
		power := minBag.Power()
		sum += power
	}
	return sum
}

func (bag Bag) Power() int {
	power := 1
	for _, v := range bag {
		power *= v
	}
	return power
}
