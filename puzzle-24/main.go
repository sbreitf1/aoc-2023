package main

// https://adventofcode.com/2023/day/23

import (
	"aoc/helper"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	lines := helper.ReadNonEmptyLines("example-1.txt")

	hails := ParseHails(lines)
	solution1 := CountIntersectionsInFuture2D(hails, helper.Point2D{X: 7, Y: 7}, helper.Point2D{X: 27, Y: 27})
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
				fmt.Println(x, y)
				os.Exit(1)
				if x >= float64(min.X) && x <= float64(max.X) && y >= float64(min.Y) && y <= float64(max.Y) {
					count++
				}
			}
		}
	}
	return count
}

func GetIntersectionInFuture2D(h1, h2 Hail) (float64, float64, bool) {

	return 0, 0, true
}
