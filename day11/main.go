package main

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

var notes []*monkeyNote

func main() {
	partOne()
	partTwo()
}

func parseNotes() {
	blocks := strings.Split(input, "\n\n")
	notes = make([]*monkeyNote, 0, len(blocks))

	for _, block := range blocks {
		notes = append(notes, parseNote(block))
	}
}

func parseNote(block string) *monkeyNote {
	var (
		n     monkeyNote
		lines = strings.Split(block, "\n")
	)

	for _, line := range lines {
		var (
			cleanLine     = strings.TrimSpace(line)
			key, value, _ = strings.Cut(cleanLine, ": ")
		)

		switch key {
		case "Starting items":
			n.worryLevels = parseInts(value)
		case "Operation":
			fmt.Sscanf(value, "new = old %s %d", &n.op, &n.opModifier)
		case "Test":
			fmt.Sscanf(value, "divisible by %d", &n.testValue)
		case "If true":
			fmt.Sscanf(value, "throw to monkey %d", &n.trueTarget)
		case "If false":
			fmt.Sscanf(value, "throw to monkey %d", &n.falseTarget)
		}
	}

	return &n
}

func parseInts(s string) []int {
	var (
		strs = strings.Split(s, ", ")
		ints = make([]int, 0, len(strs))
	)

	for _, str := range strs {
		val, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(err)
		}
		ints = append(ints, val)
	}

	return ints
}

func partOne() {
	parseNotes()
	playRounds(20, 3)
}

func partTwo() {
	parseNotes()
	playRounds(10_000, 1)
}

func playRounds(rounds, worryLevelDivider int) {
	for i := 0; i < rounds; i++ {
		round(i, worryLevelDivider)
	}

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].inspectionsCounter > notes[j].inspectionsCounter
	})

	twoMostActiveMonkeys := notes[:len(notes)-2]
	fmt.Println(twoMostActiveMonkeys[0].inspectionsCounter * twoMostActiveMonkeys[1].inspectionsCounter)
}

func round(roundNum, worryLevelDivider int) {
	limit := limit()
	
	for j, note := range notes {
		var numItemsToInspect = len(note.worryLevels)
		notes[j].inspectionsCounter += numItemsToInspect

		for len(note.worryLevels) > 0 {
			var (
				worryLevel = note.worryLevels[0]
				newLevel   = calc(note.op, worryLevel, note.opModifier) % limit / worryLevelDivider
				targetIdx  = note.trueTarget
			)

			if newLevel%note.testValue != 0 {
				targetIdx = note.falseTarget
			}

			notes[targetIdx].worryLevels = append(notes[targetIdx].worryLevels, newLevel)

			note.worryLevels = note.worryLevels[1:]
		}
	}
}

func limit() int {
	var res = 1

	for _, note := range notes {
		res *= note.testValue
	}

	return res
}

func calc(op string, a, b int) int {
	// HACK: there are no zeroes as modifiers in input. So, if we failed to parse a number,
	// we can assume that it's `old` identifier.
	if b == 0 {
		b = a
	}

	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		return a / b
	}

	return 0
}

type monkeyNote struct {
	worryLevels        []int
	op                 string
	opModifier         int
	testValue          int
	trueTarget         int
	falseTarget        int
	inspectionsCounter int
}
