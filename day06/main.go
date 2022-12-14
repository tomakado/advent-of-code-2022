package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input string

func main() {
	partOne()
	partTwo()
}

func partOne() {
	solution(4)
}

func partTwo() {
	solution(14)
}

func solution(num int) {
	for i := num; i < len(input); i++ {
		set := make(map[rune]struct{})

		for j := i - 1; j >= i - num; j-- {
			set[rune(input[j])] = struct{}{}
		}

		if len(set) == num {
			fmt.Println(i)
			break
		}
	}

}
