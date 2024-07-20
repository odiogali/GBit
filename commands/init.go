package commands

import (
	"fmt"
	"os"
)

func Init(args []string) {
	// Create folder GBit info will be stored in
	wd, _ := os.Getwd()
	var folder string

	if len(args) == 0 {
		folder = wd + "/.GBit"
	} else if len(args) == 1 {
		folder = args[0] + "/.GBit"
	}

	if err := os.Mkdir(folder, 0755); os.IsExist(err) {
		fmt.Println("Reinitialized existing GBit repository in ", wd)
	} else if err == nil {
		fmt.Println("Initialized empty GBit repository in ", folder)
	}

	// will store logs of commands that have been run
	subDirLogs := folder + "/logs"
	err := os.Mkdir(subDirLogs, 0755)
	if err != nil {
		panic(err)
	}

	// create folder where staged files will be stored
	subDirStage := folder + "/stage"
	err = os.Mkdir(subDirStage, 0755)
	if err != nil {
		panic(err)
	}

	// create record of commits
	subDirCommits := folder + "/commits"
	err = os.Mkdir(subDirCommits, 0755)
	if err != nil {
		panic(err)
	}

	// create folder 'objects'
	subDirObjects := folder + "/objects"
	err = os.Mkdir(subDirObjects, 0755)
	if err != nil {
		panic(err)
	}
}
