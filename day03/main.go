package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	rucksacks []rucksack
	groups    []group
)

func main() {
	parseRucksacks()

	partOne()
	partTwo()
}

func partOne() {
	var sum int

	for _, r := range rucksacks {
		sum += priority(r.common())
	}

	fmt.Println(sum)
}

func partTwo() {
	var sum int

	for _, g := range groups {
		sum += priority(g.common())
	}

	fmt.Println(sum)
}

func parseRucksacks() {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		rucksacks = append(
			rucksacks,
			rucksack{
				first:  line[:len(line)/2],
				second: line[len(line)/2:],
			},
		)
	}

	groups = partition(lines)
}

type group struct {
	lines []string
}

func (g group) common() rune {
	sort.Slice(g.lines, func(i, j int) bool {
		return len(g.lines[i]) < len(g.lines[j])
	})

	for _, c := range g.lines[0] {
		if strings.ContainsRune(g.lines[1], c) && strings.ContainsRune(g.lines[2], c) {
			return c
		}
	}

	return 0
}

type rucksack struct {
	first, second string
}

func (r rucksack) common() rune {
	for _, c := range r.first {
		if strings.ContainsRune(r.second, c) {
			return c
		}
	}

	return 0
}

func priority(c rune) int {
	return strings.IndexRune(alphabet, c) + 1
}

func partition(lines []string) []group {
	var groups []group

	for i := 0; i < len(lines); i += 3 {
		groups = append(groups, group{lines: lines[i : min(i+3, len(lines))]})
	}

	return groups
}

func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}
