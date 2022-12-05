package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

var pairs []pair

func main() {
	parsePairs()

	partOne()
	partTwo()
}

func parsePairs() {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		pairs = append(pairs, parsePair(line))
	}
}

func parsePair(line string) pair {
	var p pair

	fmt.Sscanf(line, "%d-%d,%d-%d", &p.min1, &p.max1, &p.min2, &p.max2)
	
	return p
}

func partOne() {
	var sum int

	for _, p := range pairs {
		if p.oneFullyContainsOther() {
			sum++
		}
	}

	fmt.Println(sum)
}

func partTwo() {
	var sum int

	for _, p := range pairs {
		if p.hasOverlap() {
			sum++
		}
	}

	fmt.Println(sum)
}

type pair struct {
	min1, max1, min2, max2 int
}

func (p pair) oneFullyContainsOther() bool {
	return p.min1 <= p.min2 && p.max1 >= p.max2 || p.min2 <= p.min1 && p.max2 >= p.max1
}

func (p pair) hasOverlap() bool {
	return p.min1 <= p.max2 && p.max1 >= p.min2
}
