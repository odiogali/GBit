package commands

import (
	"fmt"
)

func Log(args []string) {
	fmt.Printf("Log got the arguments: %v", args)
}
