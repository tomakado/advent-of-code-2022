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

var inventories []inventory

func main() {
	parseInventories()

	partOne()
	partTwo()
}

func partOne() {
	max := inventories[0].calories.sum()

	for i := 1; i < len(inventories); i++ {
		calories := inventories[i].calories.sum()
		if calories > max {
			max = calories
		}
	}

	fmt.Println(max)
}

func partTwo() {
	sort.Slice(inventories, func(i, j int) bool {
		return inventories[i].calories.sum() > inventories[j].calories.sum()
	})

	sum := inventories[0].calories.sum() + inventories[1].calories.sum() + inventories[2].calories.sum()
	fmt.Println(sum)
}

func parseInventories() {
	var current inventory

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			inventories = append(inventories, current)
			current = inventory{}
			continue
		}

		calories, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		current.calories = append(current.calories, calories)
	}
}

type inventory struct {
	calories calories
}

type calories []int

func (c calories) sum() int {
	var sum int
	for _, c := range c {
		sum += c
	}
	return sum
}
