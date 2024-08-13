package commands

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func Commit(args []string) {
	wd, _ := os.Getwd()
	gbitSubDir := wd + "/.GBit"
	stagePath := gbitSubDir + "/stage"
	commitsPath := gbitSubDir + "/commits"
	logsPath := gbitSubDir + "/logs"

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

	if len(lines) == 0 {
		fmt.Println("Staging area is empty, no commit is created.")
		os.Exit(0)
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
	commitsDir, err := os.Open(commitsPath)
	if err != nil {
		panic(err)
	}
	defer commitsDir.Close()

	if commitInfo, err := commitsDir.Stat(); err == nil {

		if commitInfo.Size() == 0 {

			// Create the commit struct
			noParents := make([]string, 0)
			commit := CommitEntity{commitName, noParents, currentTime, allObjects, nil, userInfo, commitMessage}

			// Create file the particular commit will be stored in; write to it
			commitFile, err := os.Create(commitsPath + "/" + commitName + ".json")
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

			// Write to the 'logs' file that a commit has happened
			if logFile, err := os.OpenFile(logsPath, os.O_APPEND|os.O_WRONLY, 0644); err == nil {

				defer logFile.Close()
				// To the log file, write: branch name, commitID, the word "commit", commit message
				_, err = logFile.WriteString("main " + commitName + " commit \"" + commitMessage + "\"\n")
				if err != nil {
					fmt.Println("Error writing to the log file.")
					os.Exit(1)
				}

			} else {
				fmt.Println("Error opening log file.")
				os.Exit(1)
			}

			// clear stage file - delete it, then create it anew but empty
			if _, err := os.Create(stagePath); err != nil {
				fmt.Println("Error truncating the stage file.")
				os.Exit(1)
			}

			os.Exit(0) // to avoid the next part running
		}

	}

	// Get the parents of this commit; if not on branch, we are on main branch; parent = latest main commit
	var parentName = ""
	logFile, err := os.Open(logsPath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	scanner = bufio.NewScanner(logFile)
	scanner.Split(bufio.ScanLines)
	var logLines []string
	for scanner.Scan() {
		logLines = append(logLines, scanner.Text())
	}
	for i := len(logLines) - 1; i >= 0; i-- {
		words := strings.Split(logLines[i], " ")
		// words[0] - "main", [1] - name, [2] - "commit", [3] - commit message
		if words[0] == "main" {
			parentName = words[1]
			break
		}
	}
	var parents = make([]string, 1)
	parents[0] = parentName
	commit := CommitEntity{commitName, parents, currentTime, allObjects, nil, userInfo, commitMessage}

	// Create file the particular commit will be stored in; write to it
	commitFile, err := os.Create(commitsPath + "/" + commitName + ".json")
	jsonData, err = json.Marshal(commit)
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

	// Write to the 'logs' file that a commit has happened
	if logFile, err := os.OpenFile(logsPath, os.O_APPEND|os.O_WRONLY, 0644); err == nil {

		defer logFile.Close()
		// To the log file, write: branch name, commitID, "commit", and commit message
		_, err = logFile.WriteString("main " + commitName + " commit \"" + commitMessage + "\"\n")
		if err != nil {
			fmt.Println("Error writing to the log file.")
			os.Exit(1)
		}

	} else {
		fmt.Println("Error opening log file.")
		os.Exit(1)
	}

	// clear stage file - delete it, then create it anew but empty
	if _, err := os.Create(stagePath); err != nil {
		fmt.Println("Error truncating the stage file.")
		os.Exit(1)
	}

}
