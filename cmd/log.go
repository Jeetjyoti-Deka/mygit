package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Log() {
	headPath := ".mygit/Head"
	objectsDir := ".mygit/objects"

	headData, err := os.ReadFile(headPath)
	if err != nil {
		fmt.Println("No commits yet.")
		return
	}

	hash := string(headData)

	for hash != "" {
		path := filepath.Join(objectsDir, hash)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to read commit %s: %v\n", hash, err)
			break
		}

		lines := strings.Split(string(data), "\n")
		var tree, parent, date, message string
		state := "headers"
		println(tree)

		for _, line := range lines {
			if line == "" {
				state = "message"
				continue
			}

			if state == "headers" {
				if strings.HasPrefix(line, "tree ") {
					tree = strings.TrimPrefix(line, "tree ")
				} else if strings.HasPrefix(line, "parent ") {
					parent = strings.TrimPrefix(line, "parent ")
				} else if strings.HasPrefix(line, "date ") {
					date = strings.TrimPrefix(line, "date ")
				}
			} else {
				message = line + "\n"
			}
		}

		fmt.Printf("Commit: %s\nDate: %s\n\n    %s\n", hash, date, strings.TrimSpace(message))
		fmt.Println("----")

		hash = parent
	}
}
