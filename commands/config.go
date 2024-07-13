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
	file, err := os.OpenFile("config.json", os.O_CREATE|os.O_RDWR, 0766)
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
			fmt.Println("You must specify this value: ", args[0])

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
		}
		return
	}

	// read config.json to byte slice
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var u User

	if len(args) == 2 {
		if args[0] == "user.name" {
			// clear the file before writing
			if err := os.Truncate("config.json", 0); err != nil {
				fmt.Printf("Failed to truncate: %v", err)
			}

			err := json.Unmarshal(data, &u)
			if err != nil {
				panic(err)
			}

			u.Name = args[1]

			info, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}

			_, err = file.Write(info)
			if err != nil {
				panic(err)
			}
		} else if args[0] == "user.email" {
			// clear the file before writing
			if err := os.Truncate("config.json", 0); err != nil {
				fmt.Printf("Failed to truncate: %v", err)
			}

			err := json.Unmarshal(data, &u)
			if err != nil {
				panic(err)
			}

			u.Email = args[1]

			info, err := json.Marshal(u)
			if err != nil {
				panic(err)
			}

			_, err = file.Write(info)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else if len(args) == 1 {
		if args[0] == "user.name" {
			err := json.Unmarshal(data, &u)
			if err != nil {
				panic(err)
			}

			if u.Name != "" {
				fmt.Println(u.Name)
			} else {
				fmt.Println("You must specify this value.")
			}
		} else if args[0] == "user.email" {
			err := json.Unmarshal(data, &u)
			if err != nil {
				panic(err)
			}

			if u.Email != "" {
				fmt.Println(u.Email)
			} else {
				fmt.Println("You must specify this value.")
			}
		} else {
			fmt.Println("Subcommand: ", args[0], " not supported.")
		}
	} else {
		fmt.Println("Invalid number of arguments.")
	}

}
