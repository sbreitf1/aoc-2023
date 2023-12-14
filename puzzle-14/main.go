package main

// https://adventofcode.com/2023/day/14

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	panel := ParsePanel(lines)
	panel.TiltNorth()
	solution1 := panel.ComputeNorthWeight()

	fmt.Println("-> part 1:", solution1)
}

type Panel struct {
	Rows [][]rune
}

func (p *Panel) String() string {
	lines := make([]string, len(p.Rows))
	for y := range lines {
		lines[y] = string(p.Rows[y])
	}
	return strings.Join(lines, "\n")
}

func ParsePanel(lines []string) Panel {
	rows := make([][]rune, len(lines))
	for y := range lines {
		rows[y] = []rune(lines[y])
	}
	return Panel{Rows: rows}
}

func (p *Panel) TiltNorth() {
	for y := range p.Rows {
		for x := range p.Rows[y] {
			if p.Rows[y][x] == 'O' {
				p.MoveRockNorth(x, y)
			}
		}
	}
}

func (p *Panel) MoveRockNorth(x, y int) {
	for i := y - 1; i >= 0; i-- {
		if p.Rows[i][x] != '.' {
			break
		}

		p.Rows[i][x] = 'O'
		p.Rows[i+1][x] = '.'
	}
}

func (p *Panel) ComputeNorthWeight() int {
	var weight int
	for y := range p.Rows {
		for x := range p.Rows[y] {
			if p.Rows[y][x] == 'O' {
				weight += len(p.Rows) - y
			}
		}
	}
	return weight
}
