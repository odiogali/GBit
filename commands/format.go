package commands

type readable interface {
	isDir() bool
}

type User struct {
	Name  string
	Email string
}

type Blob struct { // blobs represent files
	Name     string
	Path     string
	Contents []byte
}

type Tree struct { // trees represent folders / directories
	Name     string
	Path     string
	Children []readable
}

type Commit struct {
	Parents  []Commit
	Author   string
	Message  string
	Snapshot Tree
}

func (b Blob) isDir() bool {
	return false
}

func (t Tree) isDir() bool {
	return true
}
