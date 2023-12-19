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

	solution2 := system.CountAcceptedValues(PartRange{Categories: map[rune]ValRange{
		'x': {Min: 1, Max: 4000},
		'm': {Min: 1, Max: 4000},
		'a': {Min: 1, Max: 4000},
		's': {Min: 1, Max: 4000},
	}})
	fmt.Println("-> part 2:", solution2)
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
	Value        int64
	NextWorkflow string
}

type PartRating struct {
	Categories map[rune]int64
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
						Value:        int64(val),
						NextWorkflow: m[4],
					})
				} else {
					rules = append(rules, Rule{NextWorkflow: p})
				}
			}
			workflows[m[1]] = Workflow{Rules: rules}

		} else if m := patternPart.FindStringSubmatch(line); len(m) == 2 {
			parts := strings.Split(m[1], ",")
			categories := make(map[rune]int64)
			for _, p := range parts {
				if m := patternPartRating.FindStringSubmatch(p); len(m) == 3 {
					val, _ := strconv.Atoi(m[2])
					categories[rune(m[1][0])] = int64(val)
				}
			}
			partRatings = append(partRatings, PartRating{Categories: categories})
		}
	}
	return System{Workflows: workflows}, partRatings
}

func SumCategoryValues(parts []PartRating) int64 {
	var sum int64
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

type PartRange struct {
	Categories map[rune]ValRange
}

func (p PartRange) Size() int64 {
	if p.Categories == nil || len(p.Categories) == 0 {
		return 0
	}
	size := int64(1)
	for _, r := range p.Categories {
		if r.Max < r.Min {
			return 0
		}
		size *= (r.Max - r.Min + 1)
	}
	return size
}

type RangeMapping struct {
	Range        PartRange
	NextWorkflow string
}

type ValRange struct {
	Min, Max int64
}

func (s System) CountAcceptedValues(partRange PartRange) int64 {
	return s.CountAcceptedValuesOfWorkflow(partRange, "in")
}

func (s System) CountAcceptedValuesOfWorkflow(partRange PartRange, workflow string) int64 {
	if workflow == "A" {
		return partRange.Size()
	}
	if workflow == "R" {
		return 0
	}

	var acceptedCount int64
	w := s.Workflows[workflow]
	for _, r := range w.Rules {
		matching, remainder := r.CutRange(partRange)
		if matching.Size() > 0 {
			acceptedCount += s.CountAcceptedValuesOfWorkflow(matching, r.NextWorkflow)
		}
		if remainder.Size() < 0 {
			break
		}
		partRange = remainder
	}
	return acceptedCount
}

func (r Rule) CutRange(partRange PartRange) (PartRange, PartRange) {
	if r.Category == 0 || r.Operator == 0 {
		return partRange, PartRange{}
	}

	matchingCategories := helper.CloneMap(partRange.Categories)
	remainderCategories := helper.CloneMap(partRange.Categories)
	if r.Operator == '<' {
		matchingCategories[r.Category] = ValRange{partRange.Categories[r.Category].Min, r.Value - 1}
		remainderCategories[r.Category] = ValRange{r.Value, partRange.Categories[r.Category].Max}
		return PartRange{matchingCategories}, PartRange{remainderCategories}
	}
	if r.Operator == '>' {
		matchingCategories[r.Category] = ValRange{r.Value + 1, partRange.Categories[r.Category].Max}
		remainderCategories[r.Category] = ValRange{partRange.Categories[r.Category].Min, r.Value}
		return PartRange{matchingCategories}, PartRange{remainderCategories}
	}
	helper.ExitWithMessage("operator %q not supported", r.Operator)
	return PartRange{}, PartRange{}
}
