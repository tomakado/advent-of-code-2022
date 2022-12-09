package main

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

//go:embed input.txt
var input string

var (
	lines = strings.Split(input, "\n")
	tree  *node
)

func main() {
	parseTree()

	partOne()
	partTwo()
}

func partOne() {
	fmt.Println(tree)
	var (
		dirs = findAllDirsWithMaxSize(tree, 100_000)
		sum int
	)

	for _, dir := range dirs {
		sum += dir.Size()
	}

	fmt.Println(sum)
}

func findAllDirsWithMaxSize(n *node, maxSize int) []*node {
	var dirs []*node

	if n.type_ != nodeDir {
		return dirs
	}

	if n.Size() <= maxSize {
		dirs = append(dirs, n)
	}

	for _, child := range n.children {
		dirs = append(dirs, findAllDirsWithMaxSize(child, maxSize)...)
	}

	return dirs
}

func partTwo() {
}

func parseTree() {
	tree, _ = parseDir("/", 0)
}

func parseDir(name string, pos int) (*node, int) {
	var (
		dir = &node{
			name:     name,
			type_:    nodeDir,
			children: map[string]*node{},
		}
		i = pos
	)

	for i := pos; i < len(lines); i++ {
		line := lines[i]

		switch line {
		case "$ cd /", "$ ls", "":
			continue
		case "$ cd ..":
			if name == "/" {
				continue
			}
			return dir, i + 1
		}

		if strings.HasPrefix(line, "dir ") {
			continue
		}

		if strings.HasPrefix(line, "$ cd ") {
			targetName := strings.TrimPrefix(line, "$ cd ")
			parsedDir, shift := parseDir(targetName, i+1)
			dir.children[targetName] = parsedDir
			i += shift
			continue
		}

		sizeStr, fileName, _ := strings.Cut(line, " ")
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			log.Fatal(err, line, i)
		}

		dir.children[fileName] = &node{
			name:  fileName,
			type_: nodeFile,
			size:  size,
		}
	}

	return dir, i + 1
}

type node struct {
	name     string
	size     int
	type_    nodeType
	children map[string]*node
}

func (n *node) String() string {
	return stringNode(n, 0)
}

func (n *node) Size() int {
	if n.type_ == nodeFile {
		return n.size
	}

	size := 0
	for _, child := range n.children {
		size += child.Size()
	}

	return size
}

func stringNode(n *node, depth int) string {
	if n.type_ == nodeFile {
		return fmt.Sprintf("%s- %s (file, size=%d)\n", strings.Repeat("\t", depth), n.name, n.size)
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s- %s (total=%d)\n", strings.Repeat("\t", depth), n.name, n.Size()))

	childrenNames := make([]string, 0, len(n.children))
	for name := range n.children {
		childrenNames = append(childrenNames, name)
	}
	sort.Strings(childrenNames)

	for _, name := range childrenNames {
		b.WriteString(stringNode(n.children[name], depth+1))
	}

	if n.Size() <= 100_000 {
		return color.New(color.FgGreen).Sprint(b.String())
	}

	return b.String()
}

type nodeType int

const (
	nodeFile nodeType = iota
	nodeDir
)
