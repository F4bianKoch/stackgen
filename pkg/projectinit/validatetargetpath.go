package projectInit

import (
	"fmt"
	"os"
)

func validateTargetPath(projectPath string, force bool) (bool, error) {
	fileInfo, err := os.Stat(projectPath)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err == nil {
		if !fileInfo.IsDir() {
			return false, fmt.Errorf("a file with the name %s already exists", projectPath)
		}

		entries, err := os.ReadDir(projectPath)

		if err != nil {
			return true, fmt.Errorf("error reading target directory: %v", err)
		}

		if len(entries) == 0 {
			return true, nil
		}

		if force {
			return true, nil
		}

		return true, fmt.Errorf("%s already exists. Use --force to overwrite", projectPath)
	}

	return false, err
}
