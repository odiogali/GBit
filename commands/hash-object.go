package commands

import (
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func HashObject(args []string) {
	var nameToHash string
	var wFlag bool = false

	if len(args) == 1 {
		nameToHash = args[0]
	} else if len(args) == 2 {

		if args[0] != "-w" {
			fmt.Println("Usage: GBit hash-object -w <filename>")
			os.Exit(1)
		}

		nameToHash = args[1]
		wFlag = true
	} else {
		fmt.Println("Invalid number of arguments.")
		os.Exit(1)
	}

	wd, _ := os.Getwd()
	gbitDir := wd + "/.GBit"
	objectsDir := gbitDir + "/objects"

	// Calculate the SHA1 hash of the file based on its contents
	contentBytes, err := os.ReadFile(nameToHash)
	fileInfo, err := os.Stat(nameToHash)
	if err != nil {
		fmt.Println("Error with os.stat")
		os.Exit(1)
	}
	size := fmt.Sprintf("%d", fileInfo.Size())
	hashHeader := "blob " + size + "\000"
	contentBytes = append([]byte(hashHeader), contentBytes...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	h := sha1.New()
	h.Write(contentBytes)
	hashedByte := h.Sum(nil)
	hashedString := hex.EncodeToString(hashedByte)
	fmt.Println(hashedString)

	// Get the first two letters of the hash string and check if a folder of that name already exists
	var fileLocation string
	firstTwoLetters := hashedString[:2]
	restOfName := hashedString[2:]
	fileLocation = objectsDir + "/" + firstTwoLetters + "/"
	// If the folder exists, and wFlag is true, then put our new file there, else, create it and then put our new file there
	if !wFlag {
		os.Exit(0)
	}
	var file *os.File
	if err := os.Mkdir(fileLocation, 0755); os.IsExist(err) || err == nil {
		// folder: /objects/__ exists already or its been created
		file, err = os.Create(fileLocation + "/" + restOfName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		fmt.Println(err)
		os.Exit(1)
	}

	// Compress content using gzip and write it to the file
	w := gzip.NewWriter(file)
	_, err = w.Write(contentBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.Close()
	file.Close()
}
