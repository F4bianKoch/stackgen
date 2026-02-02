package projectInit

import (
	"os"
	"path/filepath"
)

func resolvePath(projectName string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projectPath := filepath.Join(workingDir, projectName)
	return projectPath, nil
}
