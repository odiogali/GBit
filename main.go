package main

import (
	"fmt"
)

type readable interface {
	isDir() bool
}

type blob struct {
	name     string
	path     string
	contents []byte
}

type tree struct {
	name     string
	path     string
	children []readable
}

type commit struct {
	parents  []commit
	author   string
	message  string
	snapshot tree
}

func (b blob) isDir() bool {
	return false
}

func (t tree) isDir() bool {
	return true
}

func main() {
	fmt.Println("Hello World!")
}
