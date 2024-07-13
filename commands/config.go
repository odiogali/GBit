package commands

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	Name  string
	Email string
}

func Config(args []string) {
	file, err := os.OpenFile("config.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0766)
	if err != nil {
		panic(err)
	}

	fileInfo, err := os.Stat("config.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if fileInfo.Size() == 0 {
		if len(args) == 1 {
			fmt.Println("Cannot obtain value: ", args[0])
		} else if len(args) == 2 && args[0] == "user.email" {
			u := User{"", args[1]}

			info, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}

			_, err = file.Write(info)
			if err != nil {
				panic(err)
			}
		} else if len(args) == 2 && args[0] == "user.name" {
			u := User{args[1], ""}

			info, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}

			_, err = file.Write(info)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Too many arguments were provided.")
			return
		}
	}

	if len(args) == 2 {
		if args[0] == "user.name" {
		} else if args[0] == "user.email" {
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else if len(args) == 1 {
		if args[0] == "user.name" {
		} else if args[0] == "user.email" {
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else {
		fmt.Println("Invalid number of arguments.")
	}

}
