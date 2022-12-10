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

var forest [][]int

func main() {
	parseForest()

	partOne()
	partTwo()
}

func parseForest() {
	rows := strings.Split(input, "\n")
	forest = make([][]int, 0, len(rows))

	for i, row := range rows {
		if row == "" {
			continue
		}

		forest = append(forest, make([]int, len(row)))

		for j, char := range row {
			digit, err := strconv.Atoi(string(char))
			if err != nil {
				log.Fatal(err)
			}
			forest[i][j] = digit
		}
	}
}

func partOne() {
	var sum int

	for i, row := range forest {
		for j := range row {
			if isVisible(i, j) {
				sum++
				continue
			}
		}
	}

	fmt.Println(sum)
}

func partTwo() {
	var maxScore int

	for i, row := range forest {
		for j := range row {
			if isVisible(i, j) {
				score := rayScore(i, j)
				if score > maxScore {
					maxScore = score
				}
			}
		}
	}

	fmt.Println(maxScore)
}

func isVisible(i, j int) bool {
	if i == 0 || j == 0 || i == len(forest)-1 || j == len(forest[i])-1 {
		return true
	}

	return isVisibleFromLeft(i, j) ||
		isVisibleFromRight(i, j) ||
		isVisibleFromTop(i, j) ||
		isVisibleFromBottom(i, j)
}

func isVisibleFromLeft(i, j int) bool {
	targetTree := forest[i][j]

	for k := j - 1; k >= 0; k-- {
		if forest[i][k] >= targetTree {
			return false
		}
	}

	return true
}

func isVisibleFromRight(i, j int) bool {
	targetTree := forest[i][j]

	for k := j + 1; k < len(forest[i]); k++ {
		if forest[i][k] >= targetTree {
			return false
		}
	}

	return true
}

func isVisibleFromTop(i, j int) bool {
	targetTree := forest[i][j]

	for k := i - 1; k >= 0; k-- {
		if forest[k][j] >= targetTree {
			return false
		}
	}

	return true
}

func isVisibleFromBottom(i, j int) bool {
	targetTree := forest[i][j]

	for k := i + 1; k < len(forest); k++ {
		if forest[k][j] >= targetTree {
			return false
		}
	}

	return true
}

func rayScore(i, j int) int {
	return rayLeft(i, j) * rayRight(i, j) * rayTop(i, j) * rayBottom(i, j)
}

func rayLeft(i, j int) int {
	rayLength := 0

	for k := j - 1; k >= 0; k-- {
		rayLength++

		if forest[i][k] >= forest[i][j] {
			break
		}

	}

	return rayLength
}

func rayRight(i, j int) int {
	rayLength := 0

	for k := j + 1; k < len(forest[i]); k++ {
		rayLength++

		if forest[i][k] >= forest[i][j] {
			break
		}
	}

	return rayLength
}

func rayTop(i, j int) int {
	rayLength := 0

	for k := i - 1; k >= 0; k-- {
		rayLength++

		if forest[k][j] >= forest[i][j] {
			break
		}
	}

	return rayLength
}

func rayBottom(i, j int) int {
	rayLength := 0

	for k := i + 1; k < len(forest); k++ {
		rayLength++

		if forest[k][j] >= forest[i][j] {
			break
		}
	}

	return rayLength
}

