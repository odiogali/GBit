package commands

import (
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

func CommitTree(args []string) {

	var commitMessage string
	var treeToCommit string
	parentCommit := ""

	// This code only handles the instance where only one parent commit is passed in
	if len(args) == 3 {

		if len(args[0]) == 40 && args[1] == "-m" {
			treeToCommit = args[0]
			commitMessage = args[2]
		} else {
			fmt.Printf("Error with 3 arguments provided.\n")
			os.Exit(1)
		}

	} else if len(args) == 5 {

		if len(args[0]) == 40 && args[1] == "-p" && len(args[2]) == 40 && args[3] == "-m" {
			treeToCommit = args[0]
			parentCommit = args[2]
			commitMessage = args[4]
		} else {
			fmt.Printf("Error with 5 arguments provided.\n")
			os.Exit(1)
		}

	} else {
		fmt.Println("Invalid number of arguments.")
		os.Exit(1)
	}

	timeAt1970 := time.Date(1970, time.January, 1, 0, 0, 0, 1, time.Local)
	currentTime := time.Now().String()
	splitTime := strings.Split(currentTime, " ")
	timeZone := splitTime[2]
	secSince1970 := time.Since(timeAt1970).Seconds()

	var commitContentString strings.Builder
	commitContentString.WriteString("tree " + treeToCommit)
	commitContentString.WriteString(parentCommit + "\n")
	commitContentString.WriteString("author Okabe Rintaro example@gmail.com " + fmt.Sprintf("%v", secSince1970) + " " + timeZone + "\n")
	commitContentString.WriteString("committer Okabe Rintaro example@gmail.com " + fmt.Sprintf("%v", secSince1970) + " " + timeZone + "\n")
	commitContentString.WriteString(commitMessage)

	h := sha1.New()
	h.Write([]byte(commitContentString.String()))
	hashedByte := h.Sum(nil)
	hashedString := hex.EncodeToString(hashedByte)

	wd, _ := os.Getwd()

	var fileLocation string
	firstTwoLetters := hashedString[:2]
	restOfName := hashedString[2:]
	fileLocation = wd + "/.GBit/objects/" + firstTwoLetters
	var commitFile *os.File
	if err := os.Mkdir(fileLocation, 0755); os.IsExist(err) || err == nil {
		commitFile, err = os.OpenFile(fileLocation+"/"+restOfName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	defer commitFile.Close()

	// Compress content using gzip and write it to the file
	w := gzip.NewWriter(commitFile)
	_, err := w.Write([]byte(commitContentString.String()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.Close()

	fmt.Println(hashedString)

}
