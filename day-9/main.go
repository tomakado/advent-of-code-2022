package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

var (
	cmds      []cmd
	nameToCmd = map[string]cmd{
		"L": L,
		"R": R,
		"U": U,
		"D": D,
	}
)

func main() {
	parseHistory()

	partOne()
	partTwo()
}

func parseHistory() {
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		direction, distanceStr, _ := strings.Cut(line, " ")
		distance, err := strconv.Atoi(distanceStr)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < distance; i++ {
			cmds = append(cmds, nameToCmd[direction])
		}
	}
}

func partOne() {
	fmt.Println(simulation(2))
}

func partTwo() {
	fmt.Println(simulation(10))
}

func simulation(knotsN int) int {
	var (
		knots  = make([]point, knotsN)
		visits = map[point]struct{}{{}: {}}
	)

	for _, cmd := range cmds {
		dx, dy := move(cmd)
		knots[0].x += dx
		knots[0].y += dy

		for i := 1; i < knotsN; i++ {
			knots[i] = follow(knots[i-1], knots[i])
			// fmt.Println(i, knots[i])

			visits[knots[knotsN-1]] = struct{}{}
		}
	}

	return len(visits)
}

func move(cmd cmd) (int, int) {
	switch cmd {
	case L:
		return -1, 0
	case R:
		return 1, 0
	case U:
		return 0, 1
	case D:
		return 0, -1
	default:
		return 0, 0
	}
}

func follow(head, tail point) point {
	newTail := tail

	switch (point{head.x - tail.x, head.y - tail.y}) {
	case point{-2, 1}, point{-1, 2}, point{0, 2}, point{1, 2}, point{2, 1}, point{2, 2}, point{-2, 2}:
		newTail.y++
	}

	switch (point{head.x - tail.x, head.y - tail.y}) {
	case point{1, 2}, point{2, 1}, point{2, 0}, point{2, -1}, point{1, -2}, point{2, 2}, point{2, -2}:
		newTail.x++
	}

	switch (point{head.x - tail.x, head.y - tail.y}) {
	case point{2, -1}, point{1, -2}, point{0, -2}, point{-1, -2}, point{-2, -1}, point{2, -2}, point{-2, -2}:
		newTail.y--
	}

	switch (point{head.x - tail.x, head.y - tail.y}) {
	case point{-1, -2}, point{-2, -1}, point{-2, -0}, point{-2, 1}, point{-1, 2}, point{-2, 2}, point{-2, -2}:
		newTail.x--
	}

	return newTail
}

type point struct {
	x, y int
}

func (p point) distance(other point) int {
	return abs(p.x-other.x) + abs(p.y-other.y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type cmd uint8

func (c cmd) String() string {
	switch c {
	case L:
		return "L"
	case R:
		return "R"
	case U:
		return "U"
	case D:
		return "D"
	default:
		return "unknown"
	}
}

const (
	L cmd = iota
	R
	U
	D
)
