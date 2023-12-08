package main

// https://adventofcode.com/2023/day/8

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	DirLeft  Dir = 'L'
	DirRight Dir = 'R'
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	sequence, nodes := ParseInput(lines)
	mover := NetworkMover{
		Sequence:    sequence,
		Network:     nodes,
		DirectLinks: make(map[DirectLinkHeader]string),
	}
	solution1 := GetPathLength(&mover, "AAA", "ZZZ")
	solution2 := GetGhostPathLength(&mover)

	fmt.Println("-> part 1:", solution1)
	fmt.Println("-> part 2:", solution2)
}

type Network map[string]Node

type Node struct {
	Left  string
	Right string
}

func (n Node) GetNext(d Dir) string {
	if d == DirLeft {
		return n.Left
	}
	return n.Right
}

type Dir rune

func (d Dir) String() string {
	return string(d)
}

func ParseInput(lines []string) ([]Dir, Network) {
	sequence := []Dir(lines[0])
	patternNode := regexp.MustCompile(`^([A-Z0-9]+)\s*=\s*\(\s*([A-Z0-9]+)\s*,\s*([A-Z0-9]+)\s*\)$`)
	nodes := make(Network, 0)
	for i := 1; i < len(lines); i++ {
		m := patternNode.FindStringSubmatch(lines[i])
		if len(m) == 4 {
			nodes[m[1]] = Node{
				Left:  m[2],
				Right: m[3],
			}
		}
	}
	return sequence, nodes
}

func GetPathLength(mover *NetworkMover, from, to string) int64 {
	if _, ok := mover.Network[from]; !ok {
		fmt.Println("skip solution 1")
		return -1
	}

	var count int64
	for ; ; count++ {
		if from == to {
			break
		}
		from = mover.Move(from, count, 1)
	}
	return count
}

type NetworkMover struct {
	Sequence    []Dir
	Network     Network
	DirectLinks map[DirectLinkHeader]string
}

type DirectLinkHeader struct {
	Node          string
	SequenceIndex int
	Steps         int64
}

func GetGhostPathLength(mover *NetworkMover) int64 {
	return -1
	currentPositions := GetStartPositions(mover.Network)
	var count int64
	for ; ; count++ {
		if IsEndPosition(currentPositions) {
			break
		}
		currentPositions = mover.MoveMany(currentPositions, count, 1)
	}
	return count
}

func GetStartPositions(nodes Network) []string {
	startPositions := make([]string, 0)
	for k := range nodes {
		if strings.HasSuffix(k, "A") {
			startPositions = append(startPositions, k)
		}
	}
	sort.Strings(startPositions)
	return startPositions
}

func (nm *NetworkMover) MoveMany(positions []string, sequenceIndex, steps int64) []string {
	newPositions := make([]string, len(positions))
	for i := range positions {
		newPositions[i] = nm.Move(positions[i], sequenceIndex, steps)
	}
	return newPositions
}

func (nm *NetworkMover) Move(pos string, sequenceIndex, steps int64) string {
	sequenceIndex = sequenceIndex % int64(len(nm.Sequence))
	hdr := DirectLinkHeader{
		Node:          pos,
		SequenceIndex: int(sequenceIndex),
		Steps:         steps,
	}

	if newPos, ok := nm.DirectLinks[hdr]; ok {
		return newPos
	}
	return nm.newMove(pos, sequenceIndex, steps)
}

func (nm *NetworkMover) newMove(pos string, sequenceIndex, steps int64) string {
	if steps == 1 {
		return nm.Network[pos].GetNext(nm.Sequence[sequenceIndex])
	}
	//TODO optimize
	return nm.Move(nm.Move(pos, sequenceIndex, 1), sequenceIndex+1, steps-1)
}

func DetectLoop(sequence []Dir, nodes Network, start string) (int, int) {
	fmt.Println("detect loop for", start)
	// 8811050362409 = least common multiple 19637,12643,11567,15871,14257,19099
	knownSequenceStarts := map[string]int{}
	currentPos := start
	count := 0
	for ; ; count++ {

		if count%len(sequence) == 0 {
			if loopStart, ok := knownSequenceStarts[currentPos]; ok {
				return loopStart, count - loopStart
			}
			knownSequenceStarts[currentPos] = count
		}
		currentPos = nodes[currentPos].GetNext(sequence[count%len(sequence)])
	}
}

func IsEndPosition(positions []string) bool {
	for _, p := range positions {
		if !strings.HasSuffix(p, "Z") {
			return false
		}
	}
	return true
}
