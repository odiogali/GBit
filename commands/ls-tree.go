package commands

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

func LsTree(args []string) {

	var hashName string
	var nameOnly bool = false

	if len(args) == 1 {
		hashName = args[0]
	} else if len(args) == 2 {

		if args[0] != "--name-only" {
			fmt.Printf("Commands: %s and %s are unsupported.\n", args[0], args[1])
			os.Exit(1)
		}

		hashName = args[1]
		nameOnly = true

	} else {
		fmt.Println("Insufficient number of arguments provided.")
		os.Exit(1)
	}

	wd, _ := os.Getwd()

	//	// Navigate to the object - it will be in .GBit/objects/{first two letters of SHA hash}/{rest of SHA}
	firstTwoLetters := hashName[:2]
	restOfName := hashName[2:]
	fileLocation := wd + "/.GBit/objects/" + firstTwoLetters + "/" + restOfName

	// Open and read file
	byteContents, err := os.ReadFile(fileLocation)
	if err != nil {
		fmt.Println("The file you have specified doesn't exist.")
		os.Exit(1)
	}

	// Decompress contents of the file
	r, err := gzip.NewReader(bytes.NewReader(byteContents))
	defer r.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get file contents in a string format
	buf := new(strings.Builder)
	_, err = io.Copy(buf, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stringContents := buf.String()
	splitContents := strings.Split(stringContents, " ")

	if splitContents[0] != "tree" {
		fmt.Println("Specified object is not a tree.")
		os.Exit(1)
	}

	// nameOnly flag is false
	var entryBuf strings.Builder
	var found int = 0
	if !nameOnly {
		for i := 0; i < len(splitContents); i++ {
			//fmt.Println(splitContents[i])
			if i < 1 { // Ignore the [0]'th entry
				continue
			}

			if i == 1 { // We need to get the last 6 letters of the [1]'st entry
				res := string(splitContents[i][len(splitContents[i])-6:])
				_, err = entryBuf.WriteString(res + " ")

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				continue
			}

			if i == (4+(found*3))-1 { // formula for the entries containing the names
				border := len(splitContents[i]) - 40
				entryBuf.WriteString(splitContents[i][border:] + "\t")
				entryBuf.WriteString(splitContents[i][:border] + " ")
				fmt.Println(entryBuf.String())
				entryBuf.Reset()
				found++
				continue
			}

			_, err = entryBuf.WriteString(splitContents[i] + " ")

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}
		os.Exit(0)
	}

	// nameOnly is true
	found = 0
	for i := 0; i < len(splitContents); i++ {
		if i == (4+(found*3))-1 { // formula for the entries containing the names
			border := len(splitContents[i]) - 40
			fmt.Println(splitContents[i][:border])
			found++
		}
	}

	//For writing compressed content to file if ever necessary
	//fmt.Println("File location: ", fileLocation)
	//file, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//w := gzip.NewWriter(file)
	//written, err := w.Write([]byte("tree 15\000040000 tree 8F5721F4996E30623267B961B3FA7E2A18609A32\000dir1 040000 tree F265E88DB6DD33450DCDA2662FD3A0FC48934720\000dir2 100644 blob F265E88DB6DD33450DCDA2662FD3A0FC48934720\000file1"))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Printf("Number of written files: %d", written)
	//w.Close()
	//file.Close()

}
