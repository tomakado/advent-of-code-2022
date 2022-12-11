package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

var (
	program        []instruction
	cyclesToSample = map[int]struct{}{20: {}, 60: {}, 100: {}, 140: {}, 180: {}, 220: {}}
)

func main() {
	parseProgram()

	partOne()
	partTwo()
}

func parseProgram() {
	lines := strings.Split(input, "\n")
	program = make([]instruction, 0, len(lines))

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "noop"):
			program = append(program, instruction{noop, 0})
		case strings.HasPrefix(line, "addx"):
			var v int
			fmt.Sscanf(line, "addx %d", &v)
			program = append(program, instruction{addx, v})
		}
	}
}

func partOne() {
	var (
		signals                 []int
		x                       = 1
		cycles                  int
		addxCounter             int
		currentInstructionIndex = 0
		cyclesTotal             = cyclesToRun()
	)

	for cycles < cyclesTotal {
		inst := program[currentInstructionIndex]

		cycles++

		if _, ok := cyclesToSample[cycles]; ok {
			signals = append(signals, x*cycles)
		}

		switch inst.op {
		case noop:
			currentInstructionIndex++
		case addx:
			if addxCounter < 1 {
				addxCounter++
			} else {
				addxCounter = 0
				x += inst.value
				currentInstructionIndex++
			}
		}
	}

	var sum int
	for _, signal := range signals {
		sum += signal
	}

	fmt.Println(sum)
}

func cyclesToRun() int {
	var res int

	for _, instruction := range program {
		switch instruction.op {
		case noop:
			res++
		case addx:
			res += 2
		}
	}

	return res
}

func partTwo() {
	var (
		signals                 []int
		x                       = 1
		cycles                  int
		addxCounter             int
		currentInstructionIndex = 0
		cyclesTotal             = cyclesToRun()
		crtRow, crtCol          int
		b                       strings.Builder
		pixelsDrawn             int
	)

	const pixelsToDraw = 40 * 6

	for cycles < cyclesTotal {
		if x == crtCol || x-crtCol == 1 || x-crtCol == -1 {
			b.WriteString("#")
		} else {
			b.WriteString(".")
		}

		if pixelsDrawn == pixelsToDraw {
			continue
		}

		crtCol++
		if crtCol == 40 {
			b.WriteString("\n")
			crtCol = 0
			crtRow++
		}

		inst := program[currentInstructionIndex]

		cycles++

		if _, ok := cyclesToSample[cycles]; ok {
			signals = append(signals, x*cycles)
		}

		switch inst.op {
		case noop:
			currentInstructionIndex++
		case addx:
			if addxCounter < 1 {
				addxCounter++
			} else {
				addxCounter = 0
				x += inst.value
				currentInstructionIndex++
			}
		}
	}

	fmt.Println(b.String())
}

type instruction struct {
	op    instructionType
	value int
}

type instructionType uint8

const (
	noop instructionType = iota
	addx
)
