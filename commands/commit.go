package commands

import (
	// "crypto/sha256"
	// "flag"
	"fmt"
	"os"
)

func Commit(args []string) {
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"
	fmt.Println(gbitSubDir)

}
