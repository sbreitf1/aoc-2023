package main

// https://adventofcode.com/2023/day/7

import (
	"aoc/helper"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

const (
	TypeFiveOfAKind  Type = 6
	TypeFourOfAKind  Type = 5
	TypeFullHouse    Type = 4
	TypeThreeOfAKind Type = 3
	TypeTwoPair      Type = 2
	TypeOnePair      Type = 1
	TypeHighCard     Type = 0
)

func main() {
	lines := helper.ReadNonEmptyLines("input.txt")

	bids := ParseBids(lines)
	solution1 := ComputeSolution1(bids)

	fmt.Println("-> part 1:", solution1)
}

var patternBid = regexp.MustCompile(`^([23456789TJQKA]{5})\s+(\d+)$`)

func ParseBids(lines []string) []Bid {
	bids := make([]Bid, 0, len(lines))
	for _, line := range lines {
		m := patternBid.FindStringSubmatch(line)
		if len(m) != 3 {
			helper.ExitWithMessage("invalid bid line %q", line)
		}
		bid, _ := strconv.Atoi(m[2])
		bids = append(bids, Bid{
			Hand: [5]Card{Card(m[1][0]), Card(m[1][1]), Card(m[1][2]), Card(m[1][3]), Card(m[1][4])},
			Bid:  bid,
		})
	}
	return bids
}

type Bid struct {
	Hand Hand
	Bid  int
}

type Hand [5]Card

func (h Hand) String() string {
	return string(h[0]) + string(h[1]) + string(h[2]) + string(h[3]) + string(h[4])
}

type Card rune

func (c Card) String() string {
	return string(c)
}

func (c Card) Value() int {
	if c >= '2' && c <= '9' {
		return int(c-'2') + 2
	}
	switch c {
	case 'T':
		return 10
	case 'J':
		return 11
	case 'Q':
		return 12
	case 'K':
		return 13
	case 'A':
		return 14
	default:
		panic(fmt.Sprintf("unknown card %q", c))
	}
}

type Type int

func (t Type) String() string {
	switch t {
	case TypeFiveOfAKind:
		return "five-of-a-kind"
	case TypeFourOfAKind:
		return "four-of-a-kind"
	case TypeFullHouse:
		return "full-house"
	case TypeThreeOfAKind:
		return "three-of-a-kind"
	case TypeTwoPair:
		return "two-pair"
	case TypeOnePair:
		return "one-pair"
	default:
		return "high-card"
	}
}

func (h Hand) GetType() Type {
	cardCounts := make(map[Card]int)
	for _, c := range h {
		currentCount := cardCounts[c]
		cardCounts[c] = currentCount + 1
	}

	takeNum := func(searchNum int) bool {
		for c, num := range cardCounts {
			if num == searchNum {
				delete(cardCounts, c)
				return true
			}
		}
		return false
	}

	if takeNum(5) {
		return TypeFiveOfAKind
	}
	if takeNum(4) {
		return TypeFourOfAKind
	}
	if takeNum(3) {
		if takeNum(2) {
			return TypeFullHouse
		}
		return TypeThreeOfAKind
	}
	if takeNum(2) {
		if takeNum(2) {
			return TypeTwoPair
		}
		return TypeOnePair
	}
	return TypeHighCard
}

func LessHand(h1, h2 Hand) bool {
	t1 := h1.GetType()
	t2 := h2.GetType()
	if t1 < t2 {
		return true
	}
	if t1 > t2 {
		return false
	}
	for i := 0; i < 5; i++ {
		if h1[i].Value() < h2[i].Value() {
			return true
		}
		if h1[i].Value() > h2[i].Value() {
			return false
		}
	}
	return false
}

func ComputeSolution1(bids []Bid) int {
	sort.Slice(bids, func(i, j int) bool {
		return !LessHand(bids[i].Hand, bids[j].Hand)
	})
	totalWinnings := 0
	for i := range bids {
		rank := len(bids) - i
		totalWinnings += bids[i].Bid * rank
	}
	return totalWinnings
}
