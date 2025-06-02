package main

import (
	"fmt"
	"mygit/cmd"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mygit <command> [<args>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmd.Init()
	case "add":
		cmd.Add(os.Args[2:])
	case "commit":
		message := "No Message"
		if len(os.Args) > 3 && os.Args[2] == "-m" {
			message = os.Args[3]
		}
		cmd.Commit(message)
	case "log":
		cmd.Log()
	case "status":
		cmd.Status()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

}
