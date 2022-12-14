package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

var (
	lines      []string
	startPlace int

	stacksPartOne []Stack
	stacksPartTwo []Stack
	instructions  []instruction
)

func main() {
	parseCrates()
	parseInstructions()

	partOne()
	partTwo()
}

func parseCrates() {
	lines = strings.Split(input, "\n")
	matrix := readRawMatrix()	
	stacksPartOne = transformMatrixToStacks(matrix)
	stacksPartTwo = transformMatrixToStacks(matrix)
}

func readRawMatrix() [][]string {
	var rawMatrix [][]string

	for i, line := range lines {
		if strings.HasPrefix(line, " 1 ") {
			startPlace = i + 2
			break
		}

		rawMatrix = append(rawMatrix, strings.Split(strings.ReplaceAll(line, "    ", " "), " "))
	}

	return rawMatrix
}

func transformMatrixToStacks(matrix [][]string) []Stack {
	stacks := make([]Stack, len(matrix[0]))

	for i := 0; i < len(matrix[0]); i++ {
		stacks[i] = Stack{}

		for j := len(matrix) - 1; j >= 0; j-- {
			if matrix[j][i] != "" {
				stacks[i].Push(matrix[j][i])
			}
		}
	}

	return stacks
}

func parseInstructions() {
	for i := startPlace; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}

		var inst instruction
		fmt.Sscanf(lines[i], "move %d from %d to %d", &inst.count, &inst.from, &inst.to)

		instructions = append(instructions, inst)
	}
}

func partOne() {
	for _, inst := range instructions {
		inst.moveOneByOne()
	}

	var b strings.Builder

	for _, stack := range stacksPartOne {
		topElement := stack.Pop()
		b.WriteString(strings.TrimPrefix(strings.TrimSuffix(topElement, "]"), "["))
	}

	fmt.Println(b.String())
}

func partTwo() {
	for _, inst := range instructions {
		inst.moveAll()
	}

	var b strings.Builder

	for _, stack := range stacksPartTwo {
		topElement := stack.Pop()
		b.WriteString(strings.TrimPrefix(strings.TrimSuffix(topElement, "]"), "["))
	}

	fmt.Println(b.String())
}

type instruction struct {
	count, from, to int
}

func (i instruction) moveOneByOne() {
	for j := 0; j < i.count; j++ {
		stacksPartOne[i.to-1].Push(stacksPartOne[i.from-1].Pop())
	}
}

func (i instruction) moveAll() {
	crates := make([]string, 0, i.count)

	for j := 0; j < i.count; j++ {
		crates = append(crates, stacksPartTwo[i.from-1].Pop())
	}

	for j := len(crates) - 1; j >= 0; j-- {
		stacksPartTwo[i.to-1].Push(crates[j])
	}
}
