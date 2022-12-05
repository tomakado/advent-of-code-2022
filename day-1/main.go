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

func main() {
	var (
		inventories = parseInventories()
		max         = inventories[0].calories.sum()
	)

	for i := 1; i < len(inventories); i++ {
		calories := inventories[i].calories.sum()
		if calories > max {
			max = calories
		}
	}

	fmt.Println(max)
}

func parseInventories() []inventory {
	var (
		inventories []inventory
		current     inventory
	)

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

	return inventories
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
