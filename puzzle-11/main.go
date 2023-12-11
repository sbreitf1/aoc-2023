package main

// https://adventofcode.com/2023/day/11

import (
	"aoc/helper"
	"fmt"
	"strings"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	universe1 := ParseUniverse(lines)
	universe2 := universe1.Clone()

	universe1.Expand(1)
	solution1 := universe1.GetAllShortestPathPairSum()

	universe2.Expand(999999)
	solution2 := universe2.GetAllShortestPathPairSum()

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Universe struct {
	MaxX, MaxY int
	Galaxies   []Galaxy
}

type Galaxy struct {
	X, Y int
}

func ParseUniverse(lines []string) Universe {
	galaxies := make([]Galaxy, 0)
	var maxX, maxY int
	for y := 0; y < len(lines); y++ {
		for x, r := range lines[y] {
			if r == '#' {
				galaxies = append(galaxies, Galaxy{X: x, Y: y})
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	return Universe{
		MaxX:     maxX,
		MaxY:     maxY,
		Galaxies: galaxies,
	}
}

func (u *Universe) String() string {
	lines := make([][]rune, u.MaxY+1)
	for y := range lines {
		lines[y] = []rune(strings.Repeat(".", u.MaxX+1))
	}
	for _, g := range u.Galaxies {
		lines[g.Y][g.X] = '#'
	}
	linesStr := make([]string, len(lines))
	for i := range lines {
		linesStr[i] = string(lines[i])
	}
	return strings.Join(linesStr, "\n")
}

func (u *Universe) Clone() *Universe {
	galaxies := make([]Galaxy, len(u.Galaxies))
	copy(galaxies, u.Galaxies)
	return &Universe{
		MaxX:     u.MaxX,
		MaxY:     u.MaxY,
		Galaxies: galaxies,
	}
}

func (u *Universe) Expand(size int) {
	for y := 0; y <= u.MaxY; y++ {
		if u.IsEmptyRow(y) {
			u.ExpandRow(y, size)
			y += size
		}
	}
	for x := 0; x <= u.MaxX; x++ {
		if u.IsEmptyCol(x) {
			u.ExpandCol(x, size)
			x += size
		}
	}
}

func (u *Universe) IsEmptyRow(y int) bool {
	for _, g := range u.Galaxies {
		if g.Y == y {
			return false
		}
	}
	return true
}

func (u *Universe) ExpandRow(y, size int) bool {
	u.MaxY += size
	for i := range u.Galaxies {
		if u.Galaxies[i].Y > y {
			u.Galaxies[i].Y += size
		}
	}
	return true
}

func (u *Universe) IsEmptyCol(x int) bool {
	for _, g := range u.Galaxies {
		if g.X == x {
			return false
		}
	}
	return true
}

func (u *Universe) ExpandCol(x, size int) bool {
	u.MaxX += size
	for i := range u.Galaxies {
		if u.Galaxies[i].X > x {
			u.Galaxies[i].X += size
		}
	}
	return true
}

func (u *Universe) GetAllShortestPathPairSum() int {
	var sum int
	for i := range u.Galaxies {
		for j := i + 1; j < len(u.Galaxies); j++ {
			length := u.ShortestPathLengthFromTo(u.Galaxies[i], u.Galaxies[j])
			sum += length
		}
	}
	return sum
}

func (u *Universe) ShortestPathLengthFromTo(g1, g2 Galaxy) int {
	dx := g2.X - g1.X
	if dx < 0 {
		dx = -dx
	}
	dy := g2.Y - g1.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
