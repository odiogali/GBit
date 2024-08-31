package commands

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func CatFile(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	gbitSubDir := wd + "/.GBit"

	if _, err := os.Stat(gbitSubDir); os.IsNotExist(err) {
		fmt.Println("fatal: not a GBit repository")
		os.Exit(1)
	}

	if len(args) != 2 {
		fmt.Println("Invalid number of arguments.")
		os.Exit(1)
	}

	if args[0] != "-p" {
		fmt.Println("Invalid arguments: ", args[0], args[1])
		os.Exit(1)
	}

	// Navigate to the object - it will be in .GBit/objects/{__}/{SHA}
	firstTwoLetters := args[1][:2]
	restOfName := args[1][2:]
	fileLocation := gbitSubDir + "/objects/" + firstTwoLetters + "/" + restOfName

	// Open and read file
	byteContents, err := os.ReadFile(fileLocation)
	if err != nil {
		fmt.Println("The file you have specified doesn't exist.")
		os.Exit(1)
	}

	// Decompress contents of the file
	r, err := gzip.NewReader(bytes.NewReader(byteContents))
	if err != nil {
		fmt.Println(err)
		r.Close()
		os.Exit(1)
	}
	io.Copy(os.Stdout, r)
	fmt.Println()
	r.Close()

	// For writing compressed content to file if ever necessary
	//fmt.Println("File location: ", fileLocation)
	//file, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//w := gzip.NewWriter(file)
	//written, err := w.Write([]byte("blob 11\000hello world"))
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Printf("Number of written files: %d", written)
	//w.Close()
	//file.Close()

}
