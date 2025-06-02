package cmd

import (
	"fmt"
	"os"
)

func Init() {
	err := os.Mkdir(".mygit", 0755)
	if err != nil {
		fmt.Println("Error initializing repository:", err)
		return
	}
	fmt.Println("Initialized empty MyGit repository in .mygit/")
}
