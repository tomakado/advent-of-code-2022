package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
)

const sample = `[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]`

//go:embed input.txt
var input string

func main() {
	packets := readPackets(input)

	partOne(packets)
	partTwo(packets)
}

func partOne(packets []packet) {
	var sum int

	for i, p := range packets {
		if compare(p.left, p.right) == -1 {
			sum += i + 1
		}
	}

	fmt.Println(sum)
}

func partTwo(packets []packet) {
	lists := make(lists, 0, len(packets))

	for _, p := range packets {
		lists = append(lists, p.left, p.right)
	}

	lists = append(lists, readList("[[2]]"))
	lists = append(lists, readList("[[6]]"))

	sort.Sort(lists)

	decoderKey := 1

	for i, l := range lists {
		strList := fmt.Sprintf("%v", l)
		if strList == "{[{[2]}]}" || strList == "{[{[6]}]}" {
			decoderKey *= i + 1
		}
	}

	fmt.Println(decoderKey)
}

func compare(left, right expr) int {
	if bothAreNumbers(left, right) {
		if left.(num) < right.(num) {
			return -1
		} else if left.(num) == right.(num) {
			return 0
		}

		return 1
	}

	if bothAreLists(left, right) {
		var i int

		for i < len(left.(list).elements) && i < len(right.(list).elements) {
			c := compare(left.(list).elements[i], right.(list).elements[i])

			if c != 0 {
				return c
			}

			i++
		}

		if i == len(left.(list).elements) && i < len(right.(list).elements) {
			return -1
		} else if i == len(right.(list).elements) && i < len(left.(list).elements) {
			return 1
		}

		return 0
	}

	if onlyLeftIsNumber(left, right) {
		return compare(list{elements: []expr{left}}, right)
	}

	return compare(left, list{elements: []expr{right}})
}

func bothAreNumbers(left, right expr) bool {
	return left.type_() == exprTypeNum && right.type_() == exprTypeNum
}

func bothAreLists(left, right expr) bool {
	return left.type_() == exprTypeList && right.type_() == exprTypeList
}

func onlyLeftIsNumber(left, right expr) bool {
	return left.type_() == exprTypeNum && right.type_() != exprTypeNum
}

func onlyRightIsNumber(left, right expr) bool {
	return left.type_() != exprTypeNum && right.type_() == exprTypeNum
}

func readPackets(input string) []packet {
	var (
		blocks  = strings.Split(input, "\n\n")
		packets = make([]packet, 0, len(blocks))
	)

	for _, block := range blocks {
		packets = append(packets, readBlock(block))
	}

	return packets
}

func readBlock(block string) packet {
	var (
		list1Str, list2Str, _ = strings.Cut(block, "\n")
		list1, list2          = readList(list1Str), readList(list2Str)
	)

	return packet{list1, list2}
}

func readList(listStr string) list {
	var l list

	gjson.
		Parse(listStr).
		ForEach(func(key, value gjson.Result) bool {
			if value.Type == gjson.Number {
				l.elements = append(l.elements, num(value.Int()))
				return true
			}

			l.elements = append(l.elements, readList(value.String()))
			return true
		})

	return l
}

type packet struct {
	left, right list
}

type list struct {
	elements []expr
}

func (list) expr() {}

func (list) type_() exprType { return exprTypeList }

type num int

func (num) expr() {}

func (num) type_() exprType { return exprTypeNum }

type expr interface {
	expr()
	type_() exprType
}

type exprType uint8

const (
	exprTypeNum exprType = iota
	exprTypeList
)

type lists []list

func (l lists) Len() int { return len(l) }

func (l lists) Less(i, j int) bool {
	return compare(l[i], l[j]) == -1
}

func (l lists) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

