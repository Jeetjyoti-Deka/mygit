package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Add(files []string) {
	objectsDir := ".mygit/objects"
	indexPath := ".mygit/index"

	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		fmt.Println("Failed to create objects directory:", err)
		return
	}

	index := make(map[string]string)
	if data, err := os.ReadFile(indexPath); err == nil {
		lines := strings.Split(string(data), "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 2 {
				index[parts[1]] = parts[0] // filename => hash
			}
		}
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read file %s: %v\n", file, err)
			continue
		}

		hash := fmt.Sprintf("%x", sha1.Sum(content))
		objectPath := filepath.Join(objectsDir, hash)

		if _, err := os.Stat(objectPath); os.IsNotExist(err) {
			if err := os.WriteFile(objectPath, content, 0644); err != nil {
				fmt.Printf("Failed to write object for %s: %v\n", file, err)
				continue
			}
		}

		index[file] = hash

		fmt.Printf("Added: %s -> %s\n", file, hash)
	}

	var lines []string
	for filename, hash := range index {
		lines = append(lines, fmt.Sprintf("%s %s", hash, filename))
	}

	os.WriteFile(indexPath, []byte(strings.Join(lines, "\n")), 0644)
}
