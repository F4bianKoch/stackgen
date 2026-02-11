package projectinit

import (
	"fmt"
	"os"
)

func validateTargetPath(projectPath string, force bool) error {
	fileInfo, err := os.Stat(projectPath)

	if os.IsNotExist(err) {
		fmt.Printf("Creating project directory at: %s\n\n", projectPath)
		return nil
	}

	if err == nil {
		if !fileInfo.IsDir() {
			return fmt.Errorf("a file with the name %s already exists", projectPath)
		}

		entries, err := os.ReadDir(projectPath)

		if err != nil {
			return fmt.Errorf("error reading target directory: %v", err)
		}

		if len(entries) == 0 {
			fmt.Printf("Using existing project directory at: %s\n", projectPath)
			return nil
		}

		if force {
			fmt.Printf("Overriding existing project directory at: %s\n\n", projectPath)
			return nil
		}

		return fmt.Errorf("%s already exists. Use --force to overwrite", projectPath)
	}

	return err
}
