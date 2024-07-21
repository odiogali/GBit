package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func Add(args []string) {
	// what happens when we call 'add' and not in GBit repository
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"
	objectsDir := gbitSubDir + "/objects"

	if _, err := os.Stat(gbitSubDir); os.IsNotExist(err) {
		fmt.Println("fatal: not a GBit repository")
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println("Nothing specified, nothing added.")
	} else if len(args) == 1 {

		if args[0] == "." {
			// WARNING: Need to walkthrough all files in working directory
		} else {

			dat, err := os.ReadFile(args[0]) // find, read, and store file contents
			if err != nil {
				fmt.Printf("fatal: pathspec '%s' did not match any files\n", args[0])
				os.Exit(1)
			}
			hashed := hash(dat)                       // hash the file's contents
			asciiString := hex.EncodeToString(hashed) // name of new file needs to be readable
			hashedStringPath := objectsDir + "/" + asciiString

			// if file does not already exist, create it and write to it
			if _, err := os.Stat(hashedStringPath); os.IsNotExist(err) {
				error := os.Chdir(objectsDir)
				if error != nil {
					fmt.Println(error)
					os.Exit(1)
				}

				file, error := os.Create(asciiString)
				if error != nil {
					fmt.Println("Create error: ", error)
					return
				}

				defer file.Close()
				os.Chdir(wd) // no error handling beccause we know 'wd' exists
				// WARNING: Still need to write to the file
				freq := countFreq(dat)
				fmt.Println(freq)
			}

		}

	} else {
		for _, item := range args {
			dat, err := os.ReadFile(item) // find, read, and store file contents
			if err != nil {
				fmt.Printf("fatal: pathspec '%s' did not match any files\n", item)
				os.Exit(1)
			}
			hashed := hash(dat)                       // hash the file's contents
			asciiString := hex.EncodeToString(hashed) // name of new file needs to be readable
			hashedStringPath := objectsDir + "/" + asciiString

			// if file does not already exist, create it and write to it
			if _, err := os.Stat(hashedStringPath); os.IsNotExist(err) {
				error := os.Chdir(objectsDir)
				if error != nil {
					fmt.Println(error)
					os.Exit(1)
				}

				file, error := os.Create(asciiString)
				if error != nil {
					fmt.Println("Create error: ", error)
					return
				}

				defer file.Close()
				os.Chdir(wd) // no error handling beccause we know 'wd' exists
				// WARNING: Still need to write to the file
				freq := countFreq(dat)
				fmt.Println(freq)
			}
		}
	}
}

func hash(fileContents []byte) []byte {
	h := sha256.New()
	h.Write(fileContents)
	bs := h.Sum(nil)
	return bs
}

func countFreq(data []byte) map[string]int {
	freq := make(map[string]int)
	for _, letter := range data {
		// Count number of occurences
		letterAsString := string(letter)
		num, exists := freq[letterAsString]
		if !exists {
			freq[letterAsString] = 1
		} else {
			freq[letterAsString] = num + 1
		}
	}
	return freq
}

type node struct {
	left      huffEntity
	right     huffEntity
	frequency int
}

type leaf struct {
	character string
	frequency int
}

type huffEntity interface {
	isLeaf() bool
}

func (n node) isLeaf() bool {
	return false
}

func (l leaf) isLeaf() bool {
	return true
}
