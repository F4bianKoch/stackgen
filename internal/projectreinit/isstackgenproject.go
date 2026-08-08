package internal

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/f4biankoch/stackgen/internal"
)

func isStackgenProject(projectPath string) (fs.FS, error) {
	projectFS := os.DirFS(projectPath)

	manifestInfo, err := fs.Stat(projectFS, internal.Manifest)
	if err != nil {
		return projectFS, fmt.Errorf("command has to be executed in valid stackgen project: %v", err)
	}
	if !manifestInfo.Mode().IsRegular() {
		return projectFS, fmt.Errorf("command has to be executed in valid stackgen project: %s is not a regular file", internal.Manifest)
	}

	return projectFS, nil
}
