package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

func ListTemplates() error {
	templatesPath, err := filepath.Abs("./templates")
	dir, err := os.ReadDir(templatesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template %s does not exist", templatesPath)
		}
		return err
	}

	fmt.Println("Available templates:")
	for _, dirEntry := range dir {
		fmt.Printf("  - %s\n", dirEntry.Name())
	}

	return nil
}
