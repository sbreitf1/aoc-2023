package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ExitOnError(err error, args ...interface{}) {
	// always check args
	if len(args) > 0 {
		if _, ok := args[0].(string); !ok {
			panic("first arg must be string")
		}
	}

	if err != nil {
		msg := "ERR:"
		if len(args) > 0 {
			msg += " " + fmt.Sprintf(args[0].(string), args[1:]...)
		}
		fmt.Println(msg, err.Error())
		os.Exit(1)
	}
}

func ExitWithMessage(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
	os.Exit(1)
}

func ReadLines(file string) []string {
	data, err := os.ReadFile(file)
	ExitOnError(err)
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		lines[i] = strings.Trim(lines[i], "\r")
	}
	return lines
}

func ReadNonEmptyLines(file string) []string {
	lines := ReadLines(file)
	nonEmptyLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	return nonEmptyLines
}

func ReadString(file string) string {
	data, err := os.ReadFile(file)
	ExitOnError(err)
	return string(data)
}

func SplitAndParseInts(str string, separator string) []int {
	parts := strings.Split(str, separator)
	ints := make([]int, 0, len(parts))
	for _, p := range parts {
		if len(p) > 0 {
			num, err := strconv.Atoi(p)
			ExitOnError(err)
			ints = append(ints, num)
		}
	}
	return ints
}

type Point2D struct {
	X, Y int
}

func (p Point2D) Add(p2 Point2D) Point2D {
	return Point2D{X: p.X + p2.X, Y: p.Y + p2.Y}
}

func (p Point2D) Sub(p2 Point2D) Point2D {
	return Point2D{X: p.X - p2.X, Y: p.Y - p2.Y}
}

func (p Point2D) Neg() Point2D {
	return Point2D{X: -p.X, Y: -p.Y}
}

func (p Point2D) Mul(factor int) Point2D {
	return Point2D{X: p.X * factor, Y: p.Y * factor}
}
