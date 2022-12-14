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

const (
	maxSize        = 100_000
	spaceAvailable = 70_000_000
	updateSize     = 30_000_000
)

//go:embed input.txt
var input string

var (
	lines = strings.Split(input, "\n")
	tree  = &node{
		name:     "/",
		children: map[string]*node{},
		type_:    nodeDir,
	}
)

func main() {
	parseTree()

	partOne()
	partTwo()
}

func partOne() {
	var (
		dirs = findAllDirsWithMaxSize(tree, maxSize)
		sum  int
	)

	for _, dir := range dirs {
		sum += dir.Size()
	}

	fmt.Println(sum)
}

func partTwo() {
	var (
		flatDirs           = flattenTree(tree)
		totalSpaceUsed     = tree.Size()
		currentUnusedSpace = spaceAvailable - totalSpaceUsed
		spaceToFree        = updateSize - currentUnusedSpace
	)

	sort.Slice(flatDirs, func(i, j int) bool {
		return flatDirs[i].Size() < flatDirs[j].Size()
	})


	for _, dir := range flatDirs {
		dirSize := dir.Size()

		if dirSize < spaceToFree {
			continue
		}

		fmt.Println(dirSize)
		break
	}
}

func flattenTree(n *node) []*node {
	var flatDirs []*node

	if n.type_ == nodeDir {
		flatDirs = append(flatDirs, n)
	}

	for _, child := range n.children {
		if child.type_ != nodeDir {
			continue
		}

		flatDirs = append(flatDirs, flattenTree(child)...)
	}

	return flatDirs
}

func parseTree() {
	var (
		currentDirectory = tree
	)

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}

		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "$":
			switch tokens[1] {
			case "cd":
				// change dir
				switch tokens[2] {
				case "..":
					if currentDirectory == tree {
						continue
					}
					currentDirectory = currentDirectory.parent
				case "/":
					// go to root
					currentDirectory = tree
				default:
					// go to dir
					child, ok := currentDirectory.children[tokens[2]]
					if !ok {
						child = &node{
							name:     tokens[2],
							children: map[string]*node{},
							type_:    nodeDir,
							parent:   currentDirectory,
						}
						currentDirectory.children[tokens[2]] = child
					}
					currentDirectory = child
				}
			case "ls":
				// list
				continue
			}
		case "dir":
			// directory
			currentDirectory.children[tokens[1]] = &node{
				name:     tokens[1],
				children: map[string]*node{},
				type_:    nodeDir,
				parent:   currentDirectory,
			}
			continue
		default:
			// file
			size, _ := strconv.Atoi(tokens[0])
			currentDirectory.children[tokens[1]] = &node{
				name:   tokens[1],
				size:   size,
				parent: currentDirectory,
			}
		}
	}
}

func findAllDirsWithMaxSize(n *node, maxSize int) []*node {
	if n.type_ == nodeFile {
		return nil
	}

	var dirs []*node

	if n.Size() <= maxSize {
		dirs = append(dirs, n)
	}

	for _, child := range n.children {
		dirs = append(dirs, findAllDirsWithMaxSize(child, maxSize)...)
	}

	return dirs
}

type node struct {
	name     string
	size     int
	type_    nodeType
	children map[string]*node
	parent   *node
}

func (n *node) insert(targetPath string, v *node) {
	log.Printf("target path = %s, v = %+v", targetPath, v.name)
	if targetPath == n.name {
		n.children[v.name] = v
		return
	}

	segments := strings.Split(targetPath, "/")

	if n.children == nil {
		n.children = map[string]*node{}
	}

	if _, ok := n.children[segments[0]]; !ok {
		n.children[segments[0]] = &node{
			name:     segments[0],
			children: map[string]*node{},
			type_:    nodeDir,
		}
	}

	n.children[segments[0]].insert(strings.Join(segments[1:], "/"), v)
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

	if n.Size() <= maxSize {
		return color.New(color.FgGreen).Sprint(b.String())
	}

	return b.String()
}

type nodeType int

const (
	nodeFile nodeType = iota
	nodeDir
)
