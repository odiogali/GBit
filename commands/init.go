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

		if err = os.RemoveAll(folder); err != nil {
			panic(err)
		}

		if err = os.Mkdir(folder, 0755); err != nil {
			panic(err)
		}

		fmt.Println("Reinitialized existing GBit repository in ", wd)
	} else if err == nil {
		fmt.Println("Initialized empty GBit repository in ", folder)
	}

	// create folder 'objects'
	subDirObjects := folder + "/objects"
	err := os.Mkdir(subDirObjects, 0755)
	if err != nil {
		panic(err)
	}

	subDirRefs := folder + "/refs"
	err = os.Mkdir(subDirRefs, 0755)
	if err != nil {
		panic(err)
	}

	headFileName := folder + "/HEAD"
	headFilePtr, err := os.Create(headFileName)
	if err != nil {
		panic(err)
	}

	headFilePtr.Write([]byte("ref: refs/heads/main\n"))
	headFilePtr.Close()

}
