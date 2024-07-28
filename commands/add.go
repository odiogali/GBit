package commands

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/icza/bitio"
	"os"
	"path/filepath"
	"strings"
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
	}

	if len(args) == 1 {

		// Adds all files in working directory to necessary GBit directory
		if args[0] == "." {

			var added []string

			err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Println("Error walking through file at: ", path)
					return err
				}

				if !strings.Contains(path, ".GBit") && !strings.Contains(path, ".git") && !info.IsDir() {
					huffCodes = make(map[string]string)
					dat, _ := os.ReadFile(path)               // find, read, and store file contents
					hashed := Hash(dat)                       // hash the file's contents
					asciiString := hex.EncodeToString(hashed) // name of new file needs to be readable
					hashedStringPath := objectsDir + "/" + asciiString

					// if file does not already exist, create it and write to it
					if _, err := os.Stat(hashedStringPath); os.IsNotExist(err) {
						file, error := os.Create(hashedStringPath)
						if error != nil {
							fmt.Println("Create file error: ", error)
						}

						defer file.Close()

						// Encode file contents
						freq := countFreq(dat)
						encodedText := encode(freq, dat)

						// Write encoded text
						if _, err := file.Write(encodedText); err != nil {
							panic(err)
						}

						relPath, err := filepath.Rel(wd, path)
						codesStruct := JsonCodes{relPath, huffCodes} // create struct for writing to json file
						jsonData, err := json.Marshal(codesStruct)
						if err != nil {
							fmt.Printf("Failed to marshal: %s", args[0])
						}

						// create new json file and write to it
						jsonFile, err := os.Create(objectsDir + "/" + asciiString + ".json")
						if err != nil {
							fmt.Println("Create error (json for decoding): ", error)
							os.Exit(1)
						}
						if _, err = jsonFile.Write(jsonData); err != nil {
							panic(err)
						}

						added = append(added, asciiString)

					}
				}

				return nil

			})

			if err != nil {
				panic(err)
			}

			stageFile, err := os.OpenFile(gbitSubDir+"/stage", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening staging file.")
			}
			for _, item := range added {
				fmt.Fprintln(stageFile, item)
			}
			stageFile.Close()

			os.Exit(0)
		}

		// Single argument is not ". " so we add specified file to GBit repo
		dat, err := os.ReadFile(args[0]) // find, read, and store file contents
		if err != nil {
			fmt.Printf("fatal: pathspec '%s' did not match any files\n", args[0])
			os.Exit(1)
		}

		hashed := Hash(dat)                       // hash the file's contents
		asciiString := hex.EncodeToString(hashed) // name of new file needs to be readable
		hashedStringPath := objectsDir + "/" + asciiString

		// if file does not already exist, create it and write to it
		_, error := os.Open(hashedStringPath)
		if error == nil { // no need to add a file that has already been added
			os.Exit(0)
		}

		file, error := os.Create(hashedStringPath)
		if error != nil {
			fmt.Println("Error creating encoded file.")
			os.Exit(1)
		}

		defer file.Close()

		// Actually write to file
		freq := countFreq(dat)
		encodedText := encode(freq, dat)
		// fmt.Println("This is huffCodes after all is said and done: ", huffCodes)

		// Write encoded text
		if _, err := file.Write(encodedText); err != nil {
			panic(err)
		}

		relPath, _ := filepath.Rel(wd, wd+"/"+args[0])
		var codesStruct = JsonCodes{relPath, huffCodes} // create struct for writing to json file
		jsonData, err := json.Marshal(codesStruct)
		if err != nil {
			fmt.Printf("Failed to marshal: %s", args[0])
		}

		// create new json file and write to it
		jsonFile, err := os.Create(hashedStringPath + ".json")
		if err != nil {
			fmt.Println("Create error (json for decoding): ", error)
			return
		}
		if _, err = jsonFile.Write(jsonData); err != nil {
			panic(err)
		}

		stage, err := os.OpenFile(gbitSubDir+"/stage", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening staging file.")
		}
		fmt.Fprintln(stage, asciiString)
		stage.Close()

		os.Exit(0)
	}

	var added []string
	// if number of arguments is not one
	for _, item := range args {
		huffCodes = make(map[string]string)
		dat, err := os.ReadFile(item) // find, read, and store file contents
		if err != nil {
			fmt.Printf("fatal: pathspec '%s' did not match any files\n", item)
			continue // I don't think execution halts after a bad filename
		}
		hashed := Hash(dat)                       // hash the file's contents
		asciiString := hex.EncodeToString(hashed) // name of new file needs to be readable
		hashedStringPath := objectsDir + "/" + asciiString

		// if file does not already exist, create it and write to it
		if _, err := os.Stat(hashedStringPath); os.IsNotExist(err) {
			file, error := os.Create(hashedStringPath)
			if error != nil {
				fmt.Println("Create file error: ", error)
				return
			}

			defer file.Close()

			// Actually write to file
			freq := countFreq(dat)
			encodedText := encode(freq, dat)

			// Write encoded text
			if _, err := file.Write(encodedText); err != nil {
				panic(err)
			}

			relPath, _ := filepath.Rel(wd, wd+"/"+item)
			codesStruct := JsonCodes{relPath, huffCodes} // create struct for writing to json file
			jsonData, err := json.Marshal(codesStruct)
			if err != nil {
				fmt.Printf("Failed to marshal: %s", args[0])
			}

			// create new json file and write to it
			jsonFile, err := os.Create(hashedStringPath + ".json")
			if err != nil {
				fmt.Println("Create error (json for decoding): ", error)
				os.Exit(1)
			}
			if _, err = jsonFile.Write(jsonData); err != nil {
				panic(err)
			}

			added = append(added, asciiString)

		}
	}

	stage, err := os.OpenFile(gbitSubDir+"/stage", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening staging file.")
	}
	for _, item := range added {
		fmt.Fprintln(stage, item)
	}
	stage.Close()

}

func Hash(fileContents []byte) []byte {
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
	// Change the following code to an assert perhaps
	// fmt.Println("Should contain characters and their frequencies: ", freq)
	return freq
}

func encode(frequencies map[string]int, unencodedText []byte) []byte {
	huffEntities := []huffEntity{}
	for key, value := range frequencies {
		aLeaf := leaf{key, value, nil, nil}
		huffEntities = append(huffEntities, aLeaf)
	}
	// Change the following code to an assert perhaps
	// fmt.Println("Contains just the leaves: ", len(huffEntities) == len(frequencies))

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
	// Change the following code to an assert perhaps
	// fmt.Println("Slice should have one value: ", huffEntities)
	// fmt.Println("Slice has one value: ", len(huffEntities) == 1)

	generateHuffCodes(finalTree.root, "")

	return getEncodedText(unencodedText)
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

	// Change the following code to an assert perhaps
	// fmt.Println("This is supposedly the smallest item: ", smallestItem, ". This is the slice without it: ", arr)
	return smallestItem, arr
}

func generateHuffCodes(entity huffEntity, huffcode string) {
	if entity == nil {
		return
	}

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

func getEncodedText(unencodedText []byte) []byte {
	b := &bytes.Buffer{}
	writer := bitio.NewWriter(b)

	for _, item := range unencodedText {
		char := string(item)
		code := huffCodes[char]

		for _, letter := range code {
			if string(letter) == "1" {
				_ = writer.WriteBool(true)
			} else if string(letter) == "0" {
				_ = writer.WriteBool(false)
			} else {
				panic("Invalid huffCode")
			}
		}
	}

	return b.Bytes()
}

func Decode(fileName string, jsonFileName string) []byte {
	var res []byte

	// read the json file
	jsonData, err := os.ReadFile(jsonFileName)
	if err != nil {
		panic(err)
	}
	var codeStruct JsonCodes
	if err := json.Unmarshal(jsonData, &codeStruct); err != nil {
		fmt.Println("Failed to unmarshal json file.")
		os.Exit(1)
	}
	// for verification
	// fmt.Println("Contents of the json file: ", codeStruct.Codes)

	// read encoded file
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading encoded file: ", fileName)
		os.Exit(1)
	}
	reader := bitio.NewReader(bytes.NewBuffer(file))

	bit, err := reader.ReadBool()
	strBuffer := ""
	unEncodedText := ""

	for err == nil {
		// read bit
		if bit {
			strBuffer = strBuffer + "1"
		} else {
			strBuffer = strBuffer + "0"
		}

		// check if bit is in map
		key, absent := checkForValue(codeStruct.Codes, strBuffer)
		if !absent {
			unEncodedText = unEncodedText + key
			strBuffer = ""
		}

		bit, err = reader.ReadBool()
	}

	// fmt.Println("This is the unencoded text: ", unEncodedText)
	res = []byte(unEncodedText)

	return res
}

func checkForValue(mapping map[string]string, searchValue string) (string, bool) {

	for key, value := range mapping {
		if searchValue == value {
			return key, false
		}
	}

	return "", true
}

type JsonCodes struct {
	Name  string
	Codes map[string]string `json:"jsonCodes"`
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
