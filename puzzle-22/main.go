package main

// https://adventofcode.com/2023/day/22

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"strconv"
)

var (
	dirDown = helper.Point3D{X: 0, Y: 0, Z: -1}
	dirUp   = helper.Point3D{X: 0, Y: 0, Z: 1}
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	world := ParseWorld(lines)
	world.SimulateToEnd()
	desintegratableBricks := world.GetDesintegratableBricks()
	solution1 := len(desintegratableBricks)
	fmt.Println("-> part 1:", solution1)
}

func ParseWorld(lines []string) *World {
	pattern := regexp.MustCompile(`^(\d+),(\d+),(\d+)~(\d+),(\d+),(\d+)$`)
	bricks := make([]Brick, 0, len(lines))
	for _, line := range lines {
		m := pattern.FindStringSubmatch(line)
		if len(m) == 7 {
			x1, _ := strconv.Atoi(m[1])
			y1, _ := strconv.Atoi(m[2])
			z1, _ := strconv.Atoi(m[3])
			x2, _ := strconv.Atoi(m[4])
			y2, _ := strconv.Atoi(m[5])
			z2, _ := strconv.Atoi(m[6])
			brick := Brick{
				Min: helper.Point3D{X: helper.Min(x1, x2), Y: helper.Min(y1, y2), Z: helper.Min(z1, z2)},
				Max: helper.Point3D{X: helper.Max(x1, x2), Y: helper.Max(y1, y2), Z: helper.Max(z1, z2)},
			}
			bricks = append(bricks, brick)
		}
	}
	return &World{
		Bricks: bricks,
	}
}

type World struct {
	Bricks []Brick
}

type Brick struct {
	Min, Max helper.Point3D
	Resting  bool
}

func (b *Brick) Move(dir helper.Point3D) {
	b.Min = b.Min.Add(dir)
	b.Max = b.Max.Add(dir)
}

func (w *World) SimulateToEnd() int {
	var stepCount int
	for {
		if w.MoveBricks() == 0 {
			break
		}
		stepCount++
	}
	return stepCount
}

func (w *World) MoveBricks() int {
	var movedCount int
	for i := range w.Bricks {
		if !w.Bricks[i].Resting && w.MoveBrick(i) {
			movedCount++
		}
	}
	return movedCount
}

func (w *World) MoveBrick(index int) bool {
	if w.Bricks[index].Min.Z <= 1 {
		w.Bricks[index].Resting = true
		return false
	}

	w.Bricks[index].Move(dirDown)
	collidingBricks := w.GetCollidingBricks(index)
	if len(collidingBricks) > 0 {
		w.Bricks[index].Move(dirUp)
		var isOnRestingBrick bool
		for _, otherBrickIndex := range collidingBricks {
			if w.Bricks[otherBrickIndex].Resting {
				isOnRestingBrick = true
			}
		}
		w.Bricks[index].Resting = isOnRestingBrick
		return false
	}
	return true
}

func (w *World) GetCollidingBricks(index int) []int {
	collidingBricks := make([]int, 0)
	for i := 0; i < len(w.Bricks); i++ {
		if i != index {
			if BricksCollide(w.Bricks[index], w.Bricks[i]) {
				collidingBricks = append(collidingBricks, i)
			}
		}
	}
	return collidingBricks
}

func BricksCollide(b1, b2 Brick) bool {
	if b1.Min.X > b2.Max.X || b2.Min.X > b1.Max.X {
		return false
	}
	if b1.Min.Y > b2.Max.Y || b2.Min.Y > b1.Max.Y {
		return false
	}
	if b1.Min.Z > b2.Max.Z || b2.Min.Z > b1.Max.Z {
		return false
	}
	return true
}

func (w *World) GetDesintegratableBricks() []int {
	bricks := make([]int, 0)
	for i := range w.Bricks {
		if w.CanDesintegrateBrick(i) {
			bricks = append(bricks, i)
		}
	}
	return bricks
}

func (w *World) CanDesintegrateBrick(index int) bool {
	w.Bricks[index].Move(dirUp)
	supportedBricks := w.GetCollidingBricks(index)
	w.Bricks[index].Move(dirDown)

	if len(supportedBricks) == 0 {
		// no bricks supported by this brick
		return true
	}

	// check whether this brick is the only support of the other bricks
	for _, sbi := range supportedBricks {
		w.Bricks[sbi].Move(dirDown)
		supportCount := len(w.GetCollidingBricks(sbi))
		w.Bricks[sbi].Move(dirUp)
		if supportCount == 1 {
			// this brick is the only support, cannot remove
			return false
		}
	}
	return true
}
