package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

var huffCodes = make(map[string]string)

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
		os.Exit(1)
	} else if len(args) == 1 {

		// Adds all files in working directory to necessary GBit repo
		if args[0] == "." {

			// WARNING: Need to walkthrough all files in working directory

		} else { // Just adds specified file to GBit repo

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
				os.Chdir(wd)
				// WARNING: Still need to write to the file
				freq := countFreq(dat)
				_, _ = encode(freq, dat)
				fmt.Println("This is huffCodes after all is said and done: ", huffCodes)
			}

		}

	} else { // if number of arguments is not one
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
	fmt.Println("Should contain characters and their frequencies: ", freq)
	return freq
}

func encode(frequencies map[string]int, unencodedText []byte) (string, huffTree) {
	huffEntities := []huffEntity{}
	for key, value := range frequencies {
		aLeaf := leaf{key, value, nil, nil}
		huffEntities = append(huffEntities, aLeaf)
	}
	fmt.Println("Should contain just the leaves: ", huffEntities)

	finalTree := huffTree{}     // there's no real reason for me to be using a huffTree struct here tbh
	for len(huffEntities) > 1 { // Go doesn't have while loop lol
		smallest, currentSlice := getSmallestItem(huffEntities)
		secondSmallest, updatedSlice := getSmallestItem(currentSlice)
		huffEntities = updatedSlice

		// create node with relevant children and frequency
		newNode := node{smallest, secondSmallest, (smallest.getFrequency() + secondSmallest.getFrequency())}

		// add node to the huffman tree
		finalTree.root = newNode
		huffEntities = append(huffEntities, newNode)
	}
	fmt.Println("Slice should have one value: ", huffEntities)
	fmt.Println("Slice has one value: ", len(huffEntities) == 1)

	generateHuffCodes(finalTree.root, "")

	return getEncodedText(unencodedText), finalTree
}

// Gets the smallest huffEntity in the slice and returns that and the modified slice
func getSmallestItem(arr []huffEntity) (huffEntity, []huffEntity) {
	smallestItem := arr[0]
	index := 0
	for i, item := range arr {
		if item.getFrequency() < smallestItem.getFrequency() {
			smallestItem = item
			index = i
		}
	}

	if index+1 == len(arr) {
		arr = arr[:index]
	} else {
		arr = append(arr[:index], arr[index+1:]...)
	}

	fmt.Println("This is supposedly the smallest item: ", smallestItem, ". This is the slice without it: ", arr)
	return smallestItem, arr
}

func generateHuffCodes(entity huffEntity, huffcode string) {
	if entity.isLeaf() {
		huffCodes[entity.(leaf).getChar()] = huffcode
	}

	if entity.getLeft() != nil {
		generateHuffCodes(entity.getLeft(), huffcode+"0")
	}

	if entity.getRight() != nil {
		generateHuffCodes(entity.getRight(), huffcode+"1")
	}
}

func getEncodedText(unencodedText []byte) string {
	for _, item := range unencodedText {
		fmt.Println(string(item))
		// char := string(item)
		// NOTE: Use bit writer and append bit by bit
	}
	return "" // return the result
}

type node struct {
	left      huffEntity
	right     huffEntity
	frequency int
}

type leaf struct {
	character   string
	frequency   int
	left, right huffEntity
}

type huffTree struct {
	root huffEntity
}

// for storing both leaves and nodes in the same data structure
type huffEntity interface {
	isLeaf() bool
	getFrequency() int
	getLeft() huffEntity
	getRight() huffEntity
}

func (n node) isLeaf() bool {
	return false
}

func (l leaf) isLeaf() bool {
	return true
}

func (n node) getFrequency() int {
	return n.frequency
}

func (l leaf) getFrequency() int {
	return l.frequency
}

func (l leaf) getChar() string {
	return l.character
}

func (l leaf) getLeft() huffEntity {
	return nil
}

func (l leaf) getRight() huffEntity {
	return nil
}

func (n node) getLeft() huffEntity {
	return n.left
}

func (n node) getRight() huffEntity {
	return n.right
}
