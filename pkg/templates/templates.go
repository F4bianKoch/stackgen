package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run(list bool) error {
	if list {
		listTemplates()
	}

	return nil
}

func listTemplates() {
	templatesPath, err := filepath.Abs("./templates")
	dir, err := os.ReadDir(templatesPath)
	if err != nil {
		fmt.Println("Error reading templates directory:", err)
		return
	}

	fmt.Println("Available templates:")
	for _, dirEntry := range dir {
		fmt.Printf("  - %s\n", dirEntry.Name())
	}
}
