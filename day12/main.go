package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	heightMap             [][]int
	position, destination point
)

func main() {
	readHeightMap()

	partOne()
	partTwo()
}

func readHeightMap() {
	heightMap = make([][]int, 0)
	lines := strings.Split(input, "\n")

	for y, line := range lines {
		if line == "" {
			continue
		}

		row := make([]int, len(line))

		for x, c := range line {
			switch c {
			case 'S':
				position = point{x, y}
				c = 'a'
			case 'E':
				destination = point{x, y}
				c = 'z'
			}

			row[x] = strings.IndexRune(alphabet, c)
		}

		heightMap = append(heightMap, row)
	}
}

func partOne() {
	distances := dijkstra(heightMap, position, destination)
	fmt.Println(distances[destination])
}

func partTwo() {
	var lowestPoints []point

	for y := 0; y < len(heightMap); y++ {
		for x := 0; x < len(heightMap[y]); x++ {
			if heightMap[y][x] == 0 {
				lowestPoints = append(lowestPoints, point{x, y})
			}
		}
	}

	distances := make([]int, 0, len(lowestPoints))
	for _, p := range lowestPoints {
		dst := dijkstra(heightMap, p, destination)

		// HACK: IDK why dijkstra returns zeroes for about half of the lowest points
		// so this is a workaround
		if dst[destination] == 0 {
			continue
		}

		distances = append(distances, dst[destination])
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i] < distances[j]
	})

	fmt.Println(distances[0])
}

func dijkstra(heightMap [][]int, start, end point) map[point]int {
	var (
		visited   = newSet[point]()
		toVisit   = []point{start}
		distances = map[point]int{start: 0}
	)

	for {
		if len(toVisit) == 0 {
			break
		}

		current := toVisit[0]
		visited.add(current)
		toVisit = toVisit[1:]

		if current == end {
			break
		}

		neighbours := []point{
			current.up(),
			current.down(),
			current.left(),
			current.right(),
		}

		for _, neighbour := range neighbours {
			if canVisit(heightMap, current, neighbour) {
				if distances[neighbour] == 0 {
					toVisit = append(toVisit, neighbour)
					distances[neighbour] = distances[current] + 1
				}

				if distances[neighbour] >= distances[current]+1 {
					distances[neighbour] = distances[current] + 1
				}
			}
		}

		sort.Slice(toVisit, func(i, j int) bool {
			return distances[toVisit[i]] < distances[toVisit[j]]
		})
	}

	return distances
}

func getPointsToVisit(heightMap [][]int, pos point) []entry {
	// potential directions:
	var (
		up      = pos.up()
		down    = pos.down()
		left    = pos.left()
		right   = pos.right()
		toVisit = make([]entry, 0, 4)
	)

	if canVisit(heightMap, pos, right) {
		toVisit = append(toVisit, entry{point: right, dist: 1})
	}

	if canVisit(heightMap, pos, down) {
		toVisit = append(toVisit, entry{point: down, dist: 1})
	}

	if canVisit(heightMap, pos, up) {
		toVisit = append(toVisit, entry{point: up, dist: 1})
	}

	if canVisit(heightMap, pos, left) {
		toVisit = append(toVisit, entry{point: left, dist: 1})
	}

	return toVisit
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("(x=%d, y=%d)", p.x, p.y)
}

func (p point) up() point {
	return point{p.x, p.y - 1}
}

func (p point) down() point {
	return point{p.x, p.y + 1}
}

func (p point) left() point {
	return point{p.x - 1, p.y}
}

func (p point) right() point {
	return point{p.x + 1, p.y}
}

func (p point) vectorLen() int {
	return p.x*p.x + p.y*p.y
}

func canVisit(heightMap [][]int, current, target point) bool {
	if target.x < 0 || target.y < 0 || target.y >= len(heightMap) || target.x >= len(heightMap[0]) {
		return false
	}

	var (
		currentHeight = heightMap[current.y][current.x]
		targetHeight  = heightMap[target.y][target.x]
		delta         = targetHeight - currentHeight
	)

	return delta <= 1
}

func keys(m map[point]struct{}) []point {
	var keys []point

	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].vectorLen() < keys[j].vectorLen()
	})

	return keys
}

type set[T comparable] map[T]struct{}

func newSet[T comparable](elements ...T) set[T] {
	s := set[T]{}

	for _, e := range elements {
		s.add(e)
	}

	return s
}

func (s set[T]) add(v T) {
	s[v] = struct{}{}
}

func (s set[T]) remove(v T) {
	delete(s, v)
}

func (s set[T]) contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s set[T]) first() T {
	return s.nth(0)
}

func (s set[T]) nth(n int) T {
	if len(s) < n {
		panic("nth out of range")
	}

	setKeys := make([]T, 0, len(s))

	for k := range s {
		setKeys = append(setKeys, k)
	}

	for i, p := range setKeys {
		if i == n {
			return p
		}
	}

	panic("nth out of range")
}

type entry struct {
	point point
	dist  int
}
