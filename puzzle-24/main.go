package main

// https://adventofcode.com/2023/day/24

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	hails := ParseHails(lines)
	//solution1 := CountIntersectionsInFuture2D(hails, helper.Point2D[int]{X: 7, Y: 7}, helper.Point2D[int]{X: 27, Y: 27})
	solution1 := CountIntersectionsInFuture2D(hails, helper.Point2D[float64]{X: 200000000000000, Y: 200000000000000}, helper.Point2D[float64]{X: 400000000000000, Y: 400000000000000})
	fmt.Println("-> part 1:", solution1)
}

func ParseHails(lines []string) []Hail {
	// 19, 13, 30 @ -2,  1, -2
	pattern := regexp.MustCompile(`^\s*(-?\d+)\s*,\s*(-?\d+)\s*,\s*(-?\d+)\s*@\s*(-?\d+)\s*,\s*(-?\d+)\s*,\s*(-?\d+)\s*$`)
	hails := make([]Hail, 0, len(lines))
	for _, line := range lines {
		m := pattern.FindStringSubmatch(line)
		if len(m) == 7 {
			var pos, dir helper.Point3D[int64]
			pos.X, _ = strconv.ParseInt(m[1], 10, 64)
			pos.Y, _ = strconv.ParseInt(m[2], 10, 64)
			pos.Z, _ = strconv.ParseInt(m[3], 10, 64)
			dir.X, _ = strconv.ParseInt(m[4], 10, 64)
			dir.Y, _ = strconv.ParseInt(m[5], 10, 64)
			dir.Z, _ = strconv.ParseInt(m[6], 10, 64)
			hails = append(hails, Hail{Pos: pos, Dir: dir})
		}
	}
	return hails
}

type Hail struct {
	Pos, Dir helper.Point3D[int64]
}

func CountIntersectionsInFuture2D(hails []Hail, min, max helper.Point2D[float64]) int {
	var count int
	for i := 0; i < len(hails); i++ {
		for j := i + 1; j < len(hails); j++ {
			if p, ok := GetIntersectionInFuture2D(hails[i], hails[j]); ok {
				if p.InBounds(min, max) {
					count++
				}
			}
		}
	}
	return count
}

func GetIntersectionInFuture2D(h1, h2 Hail) (helper.Point2D[float64], bool) {
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	p1 := helper.ConvertPoint2D[int64, float64](h1.Pos.XY())
	d1 := helper.ConvertPoint2D[int64, float64](h1.Dir.XY())
	p2 := helper.ConvertPoint2D[int64, float64](h2.Pos.XY())
	d2 := helper.ConvertPoint2D[int64, float64](h2.Dir.XY())
	t := p2.Sub(p1).Cross(d2.Div(d1.Cross(d2)))
	if t < 0 {
		return helper.Point2D[float64]{}, false
	}
	p := p1.Add(d1.Mul(t))
	var u float64
	if h2.Dir.X != 0 {
		u = (p.X - float64(h2.Pos.X)) / float64(h2.Dir.X)
	} else {
		u = (p.Y - float64(h2.Pos.Y)) / float64(h2.Dir.Y)
	}
	if u < 0 {
		return helper.Point2D[float64]{}, false
	}
	return p, true
}
