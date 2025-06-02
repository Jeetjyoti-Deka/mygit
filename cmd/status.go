package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Status() {
	objectsDir := ".mygit/objects"
	index := make(map[string]string)
	committed := map[string]string{}
	indexPath := ".mygit/index"

	if data, err := os.ReadFile(indexPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				index[parts[1]] = parts[0]
			}

		}
	}

	headPath := ".mygit/HEAD"
	if headHashBytes, err := os.ReadFile(headPath); err == nil {
		headHash := strings.TrimSpace(string(headHashBytes))
		commitPath := filepath.Join(objectsDir, headHash)

		if commitBytes, err := os.ReadFile(commitPath); err == nil {
			lines := strings.Split(string(commitBytes), "\n")

			for _, line := range lines {
				if strings.HasPrefix(line, "tree ") {
					treeHash := strings.TrimPrefix(line, "tree ")
					treePath := filepath.Join(objectsDir, treeHash)

					if treeBytes, err := os.ReadFile(treePath); err == nil {
						treeLines := strings.Split(string(treeBytes), "\n")
						for _, entry := range treeLines {
							if entry == "" {
								continue
							}

							parts := strings.SplitN(entry, " ", 2)
							if len(parts) == 2 {
								committed[parts[1]] = parts[0]
							}
						}
					}

				}
			}
		}
	}

	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	staged := []string{}
	modified := []string{}
	untracked := []string{}

	for _, file := range files {
		name := file.Name()

		if file.IsDir() || name == ".mygit" {
			continue
		}

		content, err := os.ReadFile(name)
		if err != nil {
			continue
		}

		hash := fmt.Sprintf("%x", sha1.Sum(content))

		if stagedHash, ok := index[name]; ok {
			if stagedHash == hash {
				staged = append(staged, name)
			} else {
				modified = append(modified, name)
			}
		} else if committedHash, ok := committed[name]; ok {
			if committedHash == hash {
				// file is committed and clean — do nothing
			} else {
				modified = append(modified, name)
			}
		} else {
			untracked = append(untracked, name)
		}

	}

	// Output result
	fmt.Println("=== Staged Files ===")
	for _, f := range staged {
		fmt.Println("  +", f)
	}

	fmt.Println("\n=== Modified (Not Staged) ===")
	for _, f := range modified {
		fmt.Println("  ~", f)
	}

	fmt.Println("\n=== Untracked Files ===")
	for _, f := range untracked {
		fmt.Println("  ?", f)
	}
}
