package commands

import (
	"crypto/sha256"
	"fmt"
	"hex"
	"os"
)

func Remove(args []string) {
	/*** The plan:
		- There will be a json file containing a slice of the objects to remove FOR EACH COMMIT
		- When we want to build a commit, we will go backwards and add object by object so long as an object
			is not in "remove slice"
	***/

	wd, _ := os.Getwd()
	gbitSubDir := wd + "./GBit"
	commitsDir := gbitSubDir + "/commits"

	if _, err := os.Stat(gbitSubDir); os.IsNotExist(err) {
		fmt.Println("fatal: not a GBit repository")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("Nothing specified, nothing removed.")
		os.Exit(1)
	}

	// Get the latest object for the filenames specified in the argument
	var toRemove []string

	for _, filename := range args {
		readBytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file: ", filename)
			continue
		}

		h := sha256.New()
		h.Write(readBytes)
		bs := h.Sum(nil)
		commitName := hex.EncodeToString(bs)
		toRemove = append(toRemove, commitName)

		fmt.Println(toRemove)

		// NOTE: We now need a staging area for remove
	}
}
