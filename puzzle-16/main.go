package main

// https://adventofcode.com/2023/day/16

import (
	"aoc/helper"
	"fmt"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	board := ParseBoard(lines)
	board.FollowBeam(Point{0, 0}, Point{1, 0})
	solution1 := board.CountEnergizedTiles()
	fmt.Println("-> part 1:", solution1)

	board.ResetEnergizedTiles()
	solution2 := board.FindMaxEnergizedTiles()
	fmt.Println("-> part 2:", solution2)
}

type Board struct {
	Width, Height int
	Tiles         [][]Tile
}

type Tile struct {
	Rune      rune
	Energized bool
}

func ParseBoard(lines []string) *Board {
	tiles := make([][]Tile, len(lines))
	for y := range lines {
		tiles[y] = make([]Tile, len(lines[y]))
		for x := range lines[y] {
			tiles[y][x] = Tile{Rune: rune(lines[y][x])}
		}
	}
	return &Board{
		Width:  len(tiles[0]),
		Height: len(tiles),
		Tiles:  tiles,
	}
}

type Point struct {
	X, Y int
}

func (p Point) Add(p2 Point) Point {
	return Point{X: p.X + p2.X, Y: p.Y + p2.Y}
}

func (b *Board) FollowBeam(pos, dir Point) {
	cache := Cache{Visited: make(map[CacheKey]bool)}
	b.followBeam(&cache, pos, dir)
}

type Cache struct {
	Visited map[CacheKey]bool
}

type CacheKey struct {
	Pos, Dir Point
}

func (b *Board) followBeam(cache *Cache, pos, dir Point) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= b.Width || pos.Y >= b.Height {
		return
	}

	key := CacheKey{Pos: pos, Dir: dir}
	if _, ok := cache.Visited[key]; ok {
		return
	}
	cache.Visited[key] = true
	b.Tiles[pos.Y][pos.X].Energized = true

	r := b.Tiles[pos.Y][pos.X].Rune
	if r == '.' {
		b.followBeam(cache, pos.Add(dir), dir)
		return
	}

	if r == '/' {
		if dir.X == 1 && dir.Y == 0 {
			newDir := Point{X: 0, Y: -1}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == 0 && dir.Y == 1 {
			newDir := Point{X: -1, Y: 0}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == -1 && dir.Y == 0 {
			newDir := Point{X: 0, Y: 1}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == 0 && dir.Y == -1 {
			newDir := Point{X: 1, Y: 0}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
	}

	if r == '\\' {
		if dir.X == 1 && dir.Y == 0 {
			newDir := Point{X: 0, Y: 1}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == 0 && dir.Y == 1 {
			newDir := Point{X: 1, Y: 0}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == -1 && dir.Y == 0 {
			newDir := Point{X: 0, Y: -1}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
		if dir.X == 0 && dir.Y == -1 {
			newDir := Point{X: -1, Y: 0}
			b.followBeam(cache, pos.Add(newDir), newDir)
			return
		}
	}

	if r == '-' {
		if dir.Y == 0 {
			b.followBeam(cache, pos.Add(dir), dir)
			return
		}
		newDir1 := Point{X: -1, Y: 0}
		newDir2 := Point{X: 1, Y: 0}
		b.followBeam(cache, pos.Add(newDir1), newDir1)
		b.followBeam(cache, pos.Add(newDir2), newDir2)
		return
	}

	if r == '|' {
		if dir.X == 0 {
			b.followBeam(cache, pos.Add(dir), dir)
			return
		}
		newDir1 := Point{X: 0, Y: -1}
		newDir2 := Point{X: 0, Y: 1}
		b.followBeam(cache, pos.Add(newDir1), newDir1)
		b.followBeam(cache, pos.Add(newDir2), newDir2)
		return
	}
}

func (b *Board) CountEnergizedTiles() int {
	var count int
	for y := range b.Tiles {
		for x := range b.Tiles[y] {
			if b.Tiles[y][x].Energized {
				count++
			}
		}
	}
	return count
}

func (b *Board) ResetEnergizedTiles() {
	for y := range b.Tiles {
		for x := range b.Tiles[y] {
			b.Tiles[y][x].Energized = false
		}
	}
}

func (b *Board) FindMaxEnergizedTiles() int {
	var max int

	for y := 0; y < b.Height; y++ {
		b.FollowBeam(Point{0, y}, Point{1, 0})
		if count := b.CountEnergizedTiles(); count > max {
			max = count
		}
		b.ResetEnergizedTiles()

		b.FollowBeam(Point{b.Width - 1, y}, Point{-1, 0})
		if count := b.CountEnergizedTiles(); count > max {
			max = count
		}
		b.ResetEnergizedTiles()
	}

	for x := 0; x < b.Width; x++ {
		b.FollowBeam(Point{x, 0}, Point{0, 1})
		if count := b.CountEnergizedTiles(); count > max {
			max = count
		}
		b.ResetEnergizedTiles()

		b.FollowBeam(Point{x, b.Height - 1}, Point{0, -1})
		if count := b.CountEnergizedTiles(); count > max {
			max = count
		}
		b.ResetEnergizedTiles()
	}

	return max
}
