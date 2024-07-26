package commands

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func Commit(args []string) {
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"
	stagePath := gbitSubDir + "/stage"

	if len(args) != 2 {
		fmt.Println("Follow the convention: GBit commit -m <message>")
		os.Exit(1)
	}

	if args[0] != "-m" {
		fmt.Println("'-m' not a valid flag.")
		os.Exit(1)
	}

	// Otherwise, arg[0] is '-m' and arg[1] is some message
	commitMessage := args[1]

	// Read stage file line by line
	stageFile, err := os.Open(stagePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(stageFile)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	allStrings := ""
	for _, line := range lines {
		allStrings = allStrings + line
	}
	stageFile.Close()

	// create commit name using sha256
	h := sha256.New()
	h.Write([]byte(allStrings))
	bs := h.Sum(nil)
	commitName := hex.EncodeToString(bs)

	fmt.Println(commitName)
}
