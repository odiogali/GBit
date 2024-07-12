package commands

import (
	"bufio"
	"fmt"
	"os"
)

func Config(args []string) {
	readFile, err := os.Create("/home/odi/Desktop/GBit/config.txt")

	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	if len(args) == 2 {
		if args[0] == "user.name" {
			fmt.Fprintln(readFile, args[1])
			if len(fileLines) == 0 {
				fmt.Fprintln(readFile, "")
			} else {
				fmt.Fprintln(readFile, fileLines[0])
			}
		} else if args[0] == "user.email" {
			if len(fileLines) == 0 {
				fmt.Fprintln(readFile, "")
			} else {
				fmt.Fprintln(readFile, fileLines[1])
			}
			fmt.Fprintln(readFile, args[1])
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else if len(args) == 1 {
		if args[0] == "user.name" {
			if len(fileLines) == 0 {
				fmt.Println("User's name has not been specified.")
			} else {
				fmt.Println(fileLines[0])
				os.Exit(1)
			}
		} else if args[0] == "user.email" {
			if len(fileLines) == 0 {
				fmt.Println("User's email has not been specified.")
			} else {
				fmt.Println(fileLines[1])
				os.Exit(1)
			}
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else {
		fmt.Println("Invalid number of arguments.")
	}

	readFile.Close()
}
