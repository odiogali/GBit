package main

type readable interface {
	isDir() bool
}

type User struct {
	name  string
	email string
}

type Blob struct { // blobs represent files
	name     string
	path     string
	contents []byte
}

type Tree struct { // trees represent folders / directories
	name     string
	path     string
	children []readable
}

type Commit struct {
	parents  []Commit
	author   string
	message  string
	snapshot Tree
}

func (b Blob) isDir() bool {
	return false
}

func (t Tree) isDir() bool {
	return true
}
