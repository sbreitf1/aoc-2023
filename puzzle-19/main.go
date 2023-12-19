package main

// https://adventofcode.com/2023/day/19

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	system, partRatings := ParseInput(lines)
	acceptedParts := system.GetAcceptedParts(partRatings)
	solution1 := SumCategoryValues(acceptedParts)
	fmt.Println("-> part 1:", solution1)

}

type System struct {
	Workflows map[string]Workflow
}

type Workflow struct {
	Rules []Rule
}

type Rule struct {
	Category     rune
	Operator     rune
	Value        int
	NextWorkflow string
}

type PartRating struct {
	Categories map[rune]int
}

func ParseInput(lines []string) (System, []PartRating) {
	patternWorkflow := regexp.MustCompile(`^(.+)\{(.*)\}$`)
	patternRule := regexp.MustCompile(`^([xmas])([<>])(\d+):(.+)$`)
	patternPart := regexp.MustCompile(`^\{(.*)\}$`)
	patternPartRating := regexp.MustCompile(`^([xmas])=(\d+)$`)

	workflows := make(map[string]Workflow)
	partRatings := make([]PartRating, 0)
	for _, line := range lines {
		if m := patternWorkflow.FindStringSubmatch(line); len(m) == 3 {
			parts := strings.Split(m[2], ",")
			rules := make([]Rule, 0, len(parts))
			for _, p := range parts {
				if m := patternRule.FindStringSubmatch(p); len(m) == 5 {
					val, _ := strconv.Atoi(m[3])
					rules = append(rules, Rule{
						Category:     rune(m[1][0]),
						Operator:     rune(m[2][0]),
						Value:        val,
						NextWorkflow: m[4],
					})
				} else {
					rules = append(rules, Rule{NextWorkflow: p})
				}
			}
			workflows[m[1]] = Workflow{Rules: rules}

		} else if m := patternPart.FindStringSubmatch(line); len(m) == 2 {
			parts := strings.Split(m[1], ",")
			categories := make(map[rune]int)
			for _, p := range parts {
				if m := patternPartRating.FindStringSubmatch(p); len(m) == 3 {
					val, _ := strconv.Atoi(m[2])
					categories[rune(m[1][0])] = val
				}
			}
			partRatings = append(partRatings, PartRating{Categories: categories})
		}
	}
	return System{Workflows: workflows}, partRatings
}

func SumCategoryValues(parts []PartRating) int {
	var sum int
	for _, p := range parts {
		for _, val := range p.Categories {
			sum += val
		}
	}
	return sum
}

func (s System) GetAcceptedParts(parts []PartRating) []PartRating {
	acceptedParts := make([]PartRating, 0)
	for _, p := range parts {
		if s.Accepts(p) {
			acceptedParts = append(acceptedParts, p)
		}
	}
	return acceptedParts
}

func (s System) Accepts(p PartRating) bool {
	currentWorkflow := "in"
	for {
		w := s.Workflows[currentWorkflow]
		currentWorkflow = w.GetNextWorkflow(p)
		if currentWorkflow == "R" {
			return false
		}
		if currentWorkflow == "A" {
			return true
		}
	}
}

func (w Workflow) GetNextWorkflow(p PartRating) string {
	for _, r := range w.Rules {
		if r.Matches(p) {
			return r.NextWorkflow
		}
	}
	helper.ExitWithMessage("no next workflow found after workflow %v", w)
	return ""
}

func (r Rule) Matches(p PartRating) bool {
	if r.Category == 0 || r.Operator == 0 {
		return true
	}
	val := p.Categories[r.Category]
	if r.Operator == '<' {
		return val < r.Value
	}
	if r.Operator == '>' {
		return val > r.Value
	}
	helper.ExitWithMessage("operator %q not supported", r.Operator)
	return false
}
