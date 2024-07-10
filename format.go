package main

type readable interface {
	isDir() bool
}

type user struct {
	name  string
	email string
}

type blob struct { // blobs represent files
	name     string
	path     string
	contents []byte
}

type tree struct { // trees represent folders / directories
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
