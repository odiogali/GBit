package main

import (
	"GBit/commands"
	"flag"
	"fmt"
	"os"
)

var commandsMap = map[string]func([]string){
	"cat-file":    commands.CatFile,
	"hash-object": commands.HashObject,
	"ls-tree":     commands.LsTree,
	"config":      commands.Config,
	"branch":      commands.Branch,
	"init":        commands.Init,
	"add":         commands.Add,
	"remove":      commands.Remove,
	"pull":        commands.Pull,
	"push":        commands.Push,
	"clone":       commands.Clone,
	"commit":      commands.Commit,
	"merge":       commands.Merge,
	"write-tree":  commands.WriteTree,
	"commit-tree": commands.CommitTree,
}

func main() {
	version := flag.Bool("version", false, "print the version number of GBit")
	v := flag.Bool("v", false, "also print the version number of GBit")

	flag.Parse()

	if *version || *v {
		fmt.Println("GBit 1.0.0")
		os.Exit(1)
	}

	subcommand := flag.Args()

	if len(subcommand) == 0 {
		fmt.Println(usage())
		os.Exit(1)
	}

	cmd, ok := commandsMap[subcommand[0]]
	if !ok {
		fmt.Println(usage())
		os.Exit(1)
	}

	cmd(subcommand[1:])
}

func usage() string {
	s := "Usage: gbit [command] [options]\nAvailable commands:\n"
	for k := range commandsMap {
		s += "-" + k + "\n"
	}
	return s
}
