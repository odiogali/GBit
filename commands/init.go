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

}
