package main

// https://adventofcode.com/2023/day/14

import (
	"aoc/helper"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	panel := ParsePanel(lines)
	panel2 := panel.Clone()
	panel.TiltNorth()
	solution1 := panel.ComputeNorthWeight()

	panel2.TiltCycles(1000000000)
	solution2 := panel2.ComputeNorthWeight()

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Panel struct {
	Width, Height int
	Rows          [][]rune
}

func (p *Panel) Clone() Panel {
	rows := make([][]rune, len(p.Rows))
	for y := range rows {
		rows[y] = make([]rune, len(p.Rows[y]))
		copy(rows[y], p.Rows[y])
	}
	return Panel{
		Height: len(rows),
		Width:  len(rows[0]),
		Rows:   rows,
	}
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
	return Panel{
		Height: len(rows),
		Width:  len(rows[0]),
		Rows:   rows,
	}
}

func (p *Panel) TiltCycles(count int) {
	lastPrintTime := time.Now()
	knownIterations := make(map[uint32]int)
	knownConfigs := make([]Panel, 0)
	for i := 0; i < count; i++ {
		if time.Since(lastPrintTime) > time.Second {
			lastPrintTime = time.Now()
			fmt.Println(i, "of", count, fmt.Sprintf("(%d %%)", int(100*i/count)))
		}

		if knownIndex, ok := knownIterations[p.Hash()]; ok {
			dstIndex := knownIndex + (count-knownIndex)%(knownIndex-i)

			for y := range p.Rows {
				copy(p.Rows[y], knownConfigs[dstIndex].Rows[y])
			}
			break
		}
		knownIterations[p.Hash()] = i
		knownConfigs = append(knownConfigs, p.Clone())

		p.TiltNorth()
		p.TiltWest()
		p.TiltSouth()
		p.TiltEast()
	}
}

func (p *Panel) Hash() uint32 {
	h := fnv.New32a()
	for y := range p.Rows {
		h.Write([]byte(string(p.Rows[y])))
	}
	return h.Sum32()
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

func (p *Panel) TiltWest() {
	for y := range p.Rows {
		for x := range p.Rows[y] {
			if p.Rows[y][x] == 'O' {
				p.MoveRockWest(x, y)
			}
		}
	}
}

func (p *Panel) MoveRockWest(x, y int) {
	for i := x - 1; i >= 0; i-- {
		if p.Rows[y][i] != '.' {
			break
		}

		p.Rows[y][i] = 'O'
		p.Rows[y][i+1] = '.'
	}
}

func (p *Panel) TiltSouth() {
	for y := len(p.Rows) - 1; y >= 0; y-- {
		for x := range p.Rows[y] {
			if p.Rows[y][x] == 'O' {
				p.MoveRockSouth(x, y)
			}
		}
	}
}

func (p *Panel) MoveRockSouth(x, y int) {
	for i := y + 1; i < p.Height; i++ {
		if p.Rows[i][x] != '.' {
			break
		}

		p.Rows[i][x] = 'O'
		p.Rows[i-1][x] = '.'
	}
}

func (p *Panel) TiltEast() {
	for y := range p.Rows {
		for x := p.Width - 1; x >= 0; x-- {
			if p.Rows[y][x] == 'O' {
				p.MoveRockEast(x, y)
			}
		}
	}
}

func (p *Panel) MoveRockEast(x, y int) {
	for i := x + 1; i < p.Height; i++ {
		if p.Rows[y][i] != '.' {
			break
		}

		p.Rows[y][i] = 'O'
		p.Rows[y][i-1] = '.'
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
