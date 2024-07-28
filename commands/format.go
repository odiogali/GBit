package commands

type readable interface {
	isDir() bool
}

type User struct {
	Name  string
	Email string
}

type Blob struct { // blobs represent files
	Name string
	Path string
}

type Tree struct { // trees represent folders / directories
	Name     string
	Children []readable
}

type CommitEntity struct {
	Ref     string
	Parents []string
	Time    string
	Objects []string
	Author  User
	Message string
	// Snapshot Tree
}

type CommitDAG struct {
	RootCommit string
	Edges      []map[string][]string
}

func (b Blob) isDir() bool {
	return false
}

func (t Tree) isDir() bool {
	return true
}
