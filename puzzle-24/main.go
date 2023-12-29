package main

// https://adventofcode.com/2023/day/23

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	hails := ParseHails(lines)
	//solution1 := CountIntersectionsInFuture2D(hails, helper.Point2D{X: 7, Y: 7}, helper.Point2D{X: 27, Y: 27})
	solution1 := CountIntersectionsInFuture2D(hails, helper.Point2D{X: 200000000000000, Y: 200000000000000}, helper.Point2D{X: 400000000000000, Y: 400000000000000})
	fmt.Println("-> part 1:", solution1)
}

func ParseHails(lines []string) []Hail {
	// 19, 13, 30 @ -2,  1, -2
	pattern := regexp.MustCompile(`^\s*(-?\d+)\s*,\s*(-?\d+)\s*,\s*(-?\d+)\s*@\s*(-?\d+)\s*,\s*(-?\d+)\s*,\s*(-?\d+)\s*$`)
	hails := make([]Hail, 0, len(lines))
	for _, line := range lines {
		m := pattern.FindStringSubmatch(line)
		if len(m) == 7 {
			var pos, dir helper.Point3D
			pos.X, _ = strconv.Atoi(m[1])
			pos.Y, _ = strconv.Atoi(m[2])
			pos.Z, _ = strconv.Atoi(m[3])
			dir.X, _ = strconv.Atoi(m[4])
			dir.Y, _ = strconv.Atoi(m[5])
			dir.Z, _ = strconv.Atoi(m[6])
			hails = append(hails, Hail{Pos: pos, Dir: dir})
		}
	}
	return hails
}

type Hail struct {
	Pos, Dir helper.Point3D
}

func CountIntersectionsInFuture2D(hails []Hail, min, max helper.Point2D) int {
	var count int
	for i := 0; i < len(hails); i++ {
		for j := i + 1; j < len(hails); j++ {
			if x, y, ok := GetIntersectionInFuture2D(hails[i], hails[j]); ok {
				if x >= float64(min.X) && x <= float64(max.X) && y >= float64(min.Y) && y <= float64(max.Y) {
					count++
				}
			}
		}
	}
	return count
}

func GetIntersectionInFuture2D(h1, h2 Hail) (float64, float64, bool) {
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	d1d2Cross := h1.Dir.XY().Cross(h2.Dir.XY())
	if d1d2Cross < 0.001 && d1d2Cross > -0.001 {
		return 0, 0, false
	}
	d2InvX := float64(h2.Dir.X) / d1d2Cross
	d2InvY := float64(h2.Dir.Y) / d1d2Cross
	dp := h2.Pos.XY().Sub(h1.Pos.XY())
	t := float64(dp.X)*d2InvY - float64(dp.Y)*d2InvX
	if t < 0 {
		return 0, 0, false
	}
	cx := float64(h1.Pos.X) + float64(h1.Dir.X)*t
	cy := float64(h1.Pos.Y) + float64(h1.Dir.Y)*t
	var u float64
	if h2.Dir.X != 0 {
		u = (cx - float64(h2.Pos.X)) / float64(h2.Dir.X)
	} else {
		u = (cy - float64(h2.Pos.Y)) / float64(h2.Dir.Y)
	}
	if u < 0 {
		return 0, 0, false
	}
	return cx, cy, true
}
