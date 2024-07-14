package commands

type readable interface {
	isDir() bool
}

type User struct {
	Name  string
	Email string
}

const (
	modified int = 1
	staged       = 2
	commited     = 3
)

type Blob struct { // blobs represent files
	Name     string
	Path     string
	Contents []byte
	state    int
}

type Tree struct { // trees represent folders / directories
	Name     string
	Path     string
	Children []readable
	state    int
}

type Commit struct {
	Ref      string
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
