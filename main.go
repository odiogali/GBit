package main

import (
	"flag"
	"fmt"
	"os"
)

var commandMap = map[string]func([]string){
	"config": commands.Config,
	"log":    commands.Log,
	"help":   commands.Help,
	"init":   commands.Init,
	"add":    commands.Add,
	"diff":   commands.Diff,
	"rm":     commands.Remove,
	"mv":     commands.Move,
	"branch": commands.Branch,
	"merge":  commands.Merge,
	"pull":   commands.Pull,
	"push":   commands.Push,
	"blame":  commands.Blame,
	"clone":  commands.Clone,
}

func main() {
	var message string
	flag.StringVar(&message, "message", "", "commit message")
	flag.StringVar(&message, "m", "", "commit message")

	var version string
	flag.StringVar(&version, "version", "", "version number")
	flag.StringVar(&version, "v", "", "version number")

	flag.Parse()

	subcommand := flag.Args()

	if len(subcommand) == 0 {
		fmt.Println(usage())
		os.Exit(1)
	}

	cmd, ok := commands[subcommand[0]]
	if !ok {
		fmt.Println(usage())
		os.Exit(1)
	}

	cmd(subcommand[1:])

	fmt.Println("Message: ", message)
}

func usage() string {
	s := "Usage: gbit [command] [options]\nAvailable commands:\n"
	for k := range commands {
		s += "-" + k + "\n"
	}
	return s
}
