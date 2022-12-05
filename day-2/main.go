package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

// A,X — Rock
// C,Z — Scissors
// B,Y — Paper
// Rock defeats Scissors, Scissors defeat Paper, Paper defeats Rock.

var (
	loseMap = map[string]string{
		"A": "C",
		"B": "A",
		"C": "B",
	}
	playerToOpponentDefeats = map[string]string{
		"X": "C",
		"Y": "A",
		"Z": "B",
	}
	opponentsToPlayerDefeats = map[string]string{
		"C": "X",
		"A": "Y",
		"B": "Z",
	}
	opponentToPlayer = map[string]string{
		"A": "X",
		"B": "Y",
		"C": "Z",
	}
	rounds []round
)

func main() {
	parseRounds()

	partOne()
	partTwo()
}

func partOne() {
	var totalScore int

	for _, round := range rounds {
		totalScore += round.playerScorePartOne()
	}

	fmt.Println(totalScore)
}

func partTwo() {
	var totalScore int

	for _, round := range rounds {
		totalScore += round.playerScorePartTwo()
	}

	fmt.Println(totalScore)
}

func parseRounds() {
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		rounds = append(rounds, parseRound(line))
	}
}

func parseRound(line string) round {
	opponent, player, _ := strings.Cut(line, " ")
	return round{
		opponent: opponent,
		player:   player,
	}
}

type round struct {
	opponent string
	player   string
}

func (r round) playerScorePartOne() int {
	return roundScore(r.opponent, r.player) + shapeScore(r.player)
}

func (r round) playerScorePartTwo() int {
	switch r.player {
	case "X": // lose
		return shapeScore(loseMap[r.opponent])
	case "Y": // draw
		return shapeScore(r.opponent) + 3
	case "Z": // win
		return shapeScore(opponentsToPlayerDefeats[r.opponent]) + 6
	}

	return 0
}

func roundScore(opponent, player string) int {
	if isDraw(opponent, player) {
		return 3
	}

	if playerToOpponentDefeats[player] == opponent {
		return 6
	}

	return 0
}

func isDraw(opponent, player string) bool {
	return opponent == "A" && player == "X" ||
		opponent == "B" && player == "Y" ||
		opponent == "C" && player == "Z"
}

func shapeScore(shape string) int {
	switch shape {
	case "A", "X":
		return 1
	case "B", "Y":
		return 2
	case "C", "Z":
		return 3
	}

	return 0
}
