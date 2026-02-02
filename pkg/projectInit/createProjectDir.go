package projectInit

import (
	"fmt"
	"os"
)

func createProjectDir(projectPath string, exists bool, force bool) error {
	if exists && force {
		fmt.Printf("Using existing project directory at: %s\n", projectPath)
		return nil
	}

	fmt.Printf("Creating project directory at: %s\n", projectPath)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	return nil
}
