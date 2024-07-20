package commands

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func Add(args []string) {
	// what happens when we call 'add' and not in GBit repository
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"

	if _, err := os.Stat(gbitSubDir); os.IsNotExist(err) {
		fmt.Println("fatal: not a GBit repository")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("Nothing specified, nothing added.")
	} else if len(args) == 1 {
		if args[0] == "." {
			// walkthrough all files in working directory

		} else {
			// add the single file
			dat, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Printf("fatal: pathspec '%s' did not match any files\n", args[0])
				os.Exit(1)
			}
			hashed := hash(dat)
			fmt.Printf("%x\n", hashed)
		}
	} else {
		for _, item := range args {
			// iterate through files specified in arguments and add them
			dat, err := os.ReadFile(item)
			if err != nil {
				fmt.Printf("fatal: pathspec '%s' did not match any files\n", item)
			}
			hashed := hash(dat)
			fmt.Printf("%x\n", hashed)
		}
	}
}

func hash(fileContents []byte) []byte {
	h := sha256.New()
	h.Write(fileContents)
	bs := h.Sum(nil)
	return bs
}
