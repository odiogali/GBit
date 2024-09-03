package commands

import (
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func WriteTree(args []string) {
	if len(args) > 0 {
		fmt.Println("write-tree does not take any arguments.")
		os.Exit(1)
	}

	wd, _ := os.Getwd()

	var entries []string

	var totalSize int = 0
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Unable to access path: ", path)
			return err
		}

		if strings.Contains(path, ".GBit") || strings.Contains(path, ".git") { // skip these particular hidden file
			return nil
		}

		var entry strings.Builder

		if info.IsDir() {

			// Hash tree and get the hash of it
			treeHash := hashTree(path)
			entry.WriteString("040000 tree " + info.Name() + "\000" + treeHash + " ")
			totalSize += len([]byte(entry.String())) // Add written bytes to totalSize
			entries = append(entries, entry.String())

		} else {

			// Permissions checking
			mode := info.Mode().String()
			if string(mode[6]) == "x" {
				entry.WriteString("100755 ")
			} else {
				entry.WriteString("100644 ")
			}

			entry.WriteString("blob ")

			// Write name of the file
			entry.WriteString(info.Name() + "\000")
			// Create blob object and write its SHA
			entry.WriteString(hashBlob(path) + " ")
			entries = append(entries, entry.String())

			totalSize += len([]byte(entry.String())) // Add written bytes to totalSize

		}

		// fmt.Println("Walk entry: ", entry.String())

		entry.Reset()

		return nil
	})

	if err != nil {
		panic(err)
	}

	//fmt.Println(entries)
}

// Recursively create tree objects returning the name of the tree; path is always a directory
func hashTree(path string) string {
	if strings.Contains(path, ".GBit") || strings.Contains(path, ".git") { // skip these particular hidden file
		return "" // return nothing, doing nothing
	}

	// fmt.Println("hash tree dir call for: ", path)
	var hashedTreeIdentifier string

	// Get children of directory
	children, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Base Case 1: if directory is empty
	if len(children) == 0 {
		hashedTreeIdentifier := createTree("tree 0\000")
		return hashedTreeIdentifier
	}

	var treeContents strings.Builder

	for _, item := range children {

		info, err := item.Info()
		if err != nil {
			fmt.Println(err)
			continue
		}

		//fmt.Printf("Child of %s: %s\n", path, info.Name())
		if info.IsDir() {
			if strings.Contains(path+"/"+item.Name(), ".GBit") || strings.Contains(path+"/"+item.Name(), ".git") {
				// If we come across any hidden files, skip over them
				continue
			}
			id := hashTree(path + "/" + item.Name())
			entry := "040000 tree " + item.Name() + "\000" + id + " "
			treeContents.WriteString(entry)
		} else {

			var entry string
			// Permissions checking
			mode := info.Mode().String()
			if string(mode[6]) == "x" {
				entry += "100755 "
			} else {
				entry += "100644 "
			}

			entry += "blob "

			// Write name of the file
			entry += info.Name() + "\000"
			// Create blob object and write its SHA
			entry += hashBlob(path+"/"+info.Name()) + " "

			treeContents.WriteString(entry)
		}

	}

	// Base Case 2: if we have finished iterating through all the children
	hashedTreeIdentifier = createTree(treeContents.String())
	return hashedTreeIdentifier
}

// Create, compress and writes the tree data after the hash has been created
func createTree(treeContent string) string {
	//fmt.Println("Create tree for this content: ", treeContent)
	wd, _ := os.Getwd()
	objectsDir := wd + "/.GBit/objects"

	// Calculate the SHA1 hash of the file based on its contents
	contentBytes := []byte(treeContent)
	size := strconv.Itoa(len(contentBytes))
	hashHeader := "tree " + size + "\000"
	contentBytes = append([]byte(hashHeader), contentBytes...)
	h := sha1.New()
	h.Write(contentBytes)
	hashedByte := h.Sum(nil)
	hashedString := hex.EncodeToString(hashedByte)

	// Get the first two letters of the hash string and check if a folder of that name already exists
	var fileLocation string
	firstTwoLetters := hashedString[:2]
	restOfName := hashedString[2:]
	fileLocation = objectsDir + "/" + firstTwoLetters
	// If the folder exists, then put our new file there, else, create it and then put our new file there
	var file *os.File
	if err := os.Mkdir(fileLocation, 0755); os.IsExist(err) || err == nil {
		// folder: /objects/__ exists already or its been created
		file, err = os.OpenFile(fileLocation+"/"+restOfName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			fmt.Println(err)
			os.Exit(1)
		} else if err != nil && strings.Contains(err.Error(), "file exists") {
			return hashedString
		}

	} else {
		fmt.Println(err)
		os.Exit(1)
	}

	// Compress content using gzip and write it to the file
	w := gzip.NewWriter(file)
	_, err := w.Write(contentBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	w.Close()
	file.Close()

	return hashedString
}

func hashBlob(path string) string {
	wd, _ := os.Getwd()
	objectsDir := wd + "/.GBit/objects"

	// Calculate the SHA1 hash of the file based on its contents
	contentBytes, err := os.ReadFile(path)
	fileInfo, err := os.Stat(path)
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

	// Get the first two letters of the hash string and check if a folder of that name already exists
	var fileLocation string
	firstTwoLetters := hashedString[:2]
	restOfName := hashedString[2:]
	fileLocation = objectsDir + "/" + firstTwoLetters
	// If the folder exists, then put our new file there, else, create it and then put our new file there
	var file *os.File
	if err := os.Mkdir(fileLocation, 0755); os.IsExist(err) || err == nil {
		// folder: /objects/__ exists already or its been created
		file, err = os.OpenFile(fileLocation+"/"+restOfName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil && !strings.Contains(err.Error(), "file exists") {
			fmt.Println(err)
			os.Exit(2)
		} else if err != nil && strings.Contains(err.Error(), "file exists") {
			return hashedString
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

	return hashedString
}
