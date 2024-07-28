package commands

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	// "path/filepath"
	"time"
)

func Commit(args []string) {
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"
	stagePath := gbitSubDir + "/stage"
	commitsPath := gbitSubDir + "/commits"
	dagPath := (commitsPath + "/DAG.json")

	if len(args) != 2 {
		fmt.Println("Follow the convention: GBit commit -m <message>")
		os.Exit(1)
	}

	if args[0] != "-m" {
		fmt.Println(args[0] + " not a valid flag.")
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

	var allObjects []string
	allStrings := ""
	for _, line := range lines {
		allStrings = allStrings + " " + line
		allObjects = append(allObjects, line)
	}
	stageFile.Close()

	// create commit name using sha256
	h := sha256.New()
	h.Write([]byte(allStrings))
	bs := h.Sum(nil)
	commitName := hex.EncodeToString(bs)

	// Get the current time to store it
	currentTime := time.Now().String()

	// Get author
	jsonData, err := os.ReadFile(wd + "/config.json")
	if err != nil {
		fmt.Println("Problem getting user info. Make sure you have initialized user information.")
		os.Exit(1)
	}
	var userInfo User
	err = json.Unmarshal(jsonData, &userInfo)

	// If DAG of commits does not already exist, we need to create it
	if _, err = os.Stat(dagPath); errors.Is(err, os.ErrNotExist) {

		// Create the commit struct
		noParents := make([]string, 0)
		commit := CommitEntity{commitName, noParents, currentTime, allObjects, userInfo, commitMessage}

		// Create file the particular commit will be stored in; write to it
		commitFile, e := os.Create(commitsPath + "/" + commitName + ".json")
		jsonData, err := json.Marshal(commit)
		if err != nil {
			fmt.Println("Unable to marshal json data that will store commit info.")
			os.Exit(1)
		}
		_, err = commitFile.Write(jsonData)
		if err != nil {
			fmt.Println("Unable to write to JSON file that will store commit info.")
			os.Exit(1)
		}
		commitFile.Close()

		// Create the DAG
		dagFile, e := os.Create(dagPath)
		if e != nil {
			fmt.Println("Commit file could not be created due to some error.")
			os.Exit(1)
		}
		var dagEdges = make([]map[string][]string, 1)
		dag := CommitDAG{commitName, dagEdges}
		dagJson, err := json.Marshal(dag)
		if err != nil {
			fmt.Println("Unable to marshal JSON data for DAG used to store commit history.")
			os.Exit(1)
		}
		_, err = dagFile.Write(dagJson)
		if err != nil {
			fmt.Println("Unable to write to DAG to JSON file that will store info about commit history.")
			os.Exit(1)
		}

		dagFile.Close()

	}

	// Deal with if file exists
	if dagData, err := os.ReadFile(dagPath); err == nil {

	}
}
