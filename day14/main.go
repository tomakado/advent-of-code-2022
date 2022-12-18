package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

const sample = `498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9`

//go:embed input.txt
var input string

var sandStartPoint = point{x: 500, y: 0}

func main() {
	c := parseCave(input)

	fmt.Println(partOne(&c))
	fmt.Println(partTwo(&c))
}

func parseCave(input string) cave {
	var (
		c = cave{
			minX:  math.MaxInt,
			floor: map[point]int{},
		}
		lines = strings.Split(input, "\n")
	)

	for _, line := range lines {
		if line == "" {
			continue
		}

		s := parseShape(line)
		c.shapes = append(c.shapes, s)

		for _, p := range s.points {
			switch {
			case p.x > c.maxX:
				c.maxX = p.x
			case p.x < c.minX:
				c.minX = p.x
			case p.y < c.minY:
				c.minY = p.y
			case p.y > c.maxY:
				c.maxY = p.y
			}
		}
	}

	return c
}

func parseShape(line string) shape {
	var (
		rawPoints = strings.Split(line, " -> ")
		s         = shape{points: make([]point, 0, len(rawPoints))}
	)

	for _, rawPoint := range rawPoints {
		var p point
		fmt.Sscanf(rawPoint, "%d,%d", &p.x, &p.y)

		s.points = append(s.points, p)
	}

	return s
}

func partOne(c *cave) int {
	c.fillMatrix()

	var (
		sandIsFalling        = true
		crossedBounds        = false
		blocksDropped        int
		currentBlockPosition = sandStartPoint
	)

	for !crossedBounds {
		if sandIsFalling {
			// continue falling...

			pointsToTry := []point{
				currentBlockPosition.down(),
				currentBlockPosition.leftDown(),
				currentBlockPosition.rightDown(),
			}

			pointsTried := 0

			for _, p := range pointsToTry {
				nextBlock := c.getBlockPartOne(p)
				if nextBlock == air {
					currentBlockPosition = p
					break
				}

				if nextBlock == void {
					return blocksDropped
				}

				pointsTried++
			}

			if pointsTried == len(pointsToTry) {
				sandIsFalling = false
				c.setBlockPartOne(currentBlockPosition, sand)
				blocksDropped++
			}

			continue
		}

		currentBlockPosition = sandStartPoint
		sandIsFalling = true
	}

	return blocksDropped
}

func partTwo(c *cave) int {
	c.fillMatrix()

	var (
		sandIsFalling        = true
		sandStuckInSource    = false
		blocksDropped        int
		currentBlockPosition = sandStartPoint
	)

	for !sandStuckInSource {
		if sandIsFalling {
			// continue falling...

			pointsToTry := []point{
				currentBlockPosition.down(),
				currentBlockPosition.leftDown(),
				currentBlockPosition.rightDown(),
			}

			pointsTried := 0

			for _, p := range pointsToTry {
				nextBlock := c.getBlockPartTwo(p)
				if nextBlock == air {
					currentBlockPosition = p
					break
				}

				pointsTried++
			}

			if pointsTried == len(pointsToTry) {
				if currentBlockPosition == sandStartPoint {
					sandStuckInSource = true
				}

				sandIsFalling = false
				c.setBlockPartTwo(currentBlockPosition, sand)
				blocksDropped++
			}

			continue
		}

		currentBlockPosition = sandStartPoint
		sandIsFalling = true
	}

	return blocksDropped
}

type cave struct {
	minX, maxX, minY, maxY int
	shapes                 []shape
	matrix                 [][]int
	floor                  map[point]int
}

func (c *cave) fillMatrix() {
	c.matrix = make([][]int, c.maxY-c.minY+1)

	for y := 0; y < c.maxY+1; y++ {
		c.matrix[y] = make([]int, c.maxX-c.minX+1)
	}

	for _, s := range c.shapes {
		for i := 0; i < len(s.points)-1; i++ {
			var (
				p     = s.points[i]
				pNext = s.points[i+1]
			)

			switch {
			case p.x < pNext.x:
				for x := p.x; x <= pNext.x; x++ {
					c.setBlockPartOne(point{x: x, y: p.y}, rock)
				}
			case p.x > pNext.x:
				for x := p.x; x >= pNext.x; x-- {
					c.setBlockPartOne(point{x: x, y: p.y}, rock)
				}
			case p.y < pNext.y:
				for y := p.y; y <= pNext.y; y++ {
					c.setBlockPartOne(point{x: p.x, y: y}, rock)
				}
			case p.y > pNext.y:
				for y := p.y; y >= pNext.y; y-- {
					c.setBlockPartOne(point{x: p.x, y: y}, rock)
				}
			}
		}
	}
}

func (c cave) printMatrix() {
	for y := 0; y < len(c.matrix); y++ {
		for x := 0; x < len(c.matrix[y]); x++ {
			block, ok := c.floor[point{x, y}]
			if !ok {
				block = c.matrix[y][x]
			}

			switch block {
			case air:
				fmt.Print(".")
			case rock:
				fmt.Print("#")
			case sand:
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
}

func (c *cave) setBlockPartOne(p point, block int) {
	normalized := c.normalizePoint(p)
	c.matrix[normalized.y][normalized.x] = block
}

func (c *cave) setBlockPartTwo(p point, block int) {
	normalized := c.normalizePoint(p)

	if c.isOutOfBounds(normalized) {
		c.floor[normalized] = block
		return
	}

	c.matrix[normalized.y][normalized.x] = block

}

func (c cave) getBlockPartOne(p point) int {
	normalized := c.normalizePoint(p)

	if c.isOutOfBounds(normalized) {
		return void
	}

	return c.matrix[p.y-c.minY][p.x-c.minX]
}

func (c cave) getBlockPartTwo(p point) int {
	normalized := c.normalizePoint(p)

	if c.isOutOfBounds(normalized) {
		// void:
		if normalized.y >= c.maxY+2 {
			return void
		}

		block, ok := c.floor[normalized]
		if !ok {
			return air
		}
		return block
	}

	return c.matrix[p.y-c.minY][p.x-c.minX]
}

func (c cave) normalizePoint(p point) point {
	return point{
		x: p.x - c.minX,
		y: p.y - c.minY,
	}
}

func (c cave) isOutOfBounds(p point) bool {
	return p.y < 0 || p.y >= len(c.matrix) || p.x < 0 || p.x >= len(c.matrix[p.y])
}

type shape struct {
	points []point
}

type point struct {
	x, y int
}

func (p point) down() point {
	return point{x: p.x, y: p.y + 1}
}

func (p point) leftDown() point {
	return point{x: p.x - 1, y: p.y + 1}
}

func (p point) rightDown() point {
	return point{x: p.x + 1, y: p.y + 1}
}

const (
	air = iota
	rock
	sand
	void
)
