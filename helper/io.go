package helper

import (
	"os"
	"strconv"
	"strings"
)

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

func SplitAndTrim(str string, separator string) []string {
	parts := strings.Split(str, separator)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
