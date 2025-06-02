package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Commit(message string) {
	indexPath := ".mygit/index"
	headPath := ".mygit/HEAD"
	objectsDir := ".mygit/objects"

	indexContent, err := os.ReadFile(indexPath)
	if err != nil || len(indexContent) == 0 {
		fmt.Println("Nothing to commit. Use `add` to stage files.")
		return
	}

	lines := strings.Split(string(indexContent), "\n")
	var entries []string
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		hash, file := parts[0], parts[1]
		entries = append(entries, fmt.Sprintf("%s %s", hash, file))
	}

	treeContent := strings.Join(entries, "\n")
	treeHash := fmt.Sprintf("%x", sha1.Sum([]byte(treeContent)))
	treePath := filepath.Join(objectsDir, treeHash)
	os.WriteFile(treePath, []byte(treeContent), 0644)

	commitMessage := message
	timestamp := time.Now().Format(time.RFC3339)

	parent := ""
	if data, err := os.ReadFile(headPath); err == nil {
		parent = string(data)
	}

	commitContent := fmt.Sprintf(
		"tree %s\nparent %s\ndate %s\n\n%s\n",
		treeHash, parent, timestamp, commitMessage,
	)

	commitHash := fmt.Sprintf("%x", sha1.Sum([]byte(commitContent)))
	commitPath := filepath.Join(objectsDir, commitHash)
	os.WriteFile(commitPath, []byte(commitContent), 0644)

	os.WriteFile(headPath, []byte(commitHash), 0644)

	os.Remove(indexPath)
	fmt.Println("Committed with hash:", commitHash)
}
