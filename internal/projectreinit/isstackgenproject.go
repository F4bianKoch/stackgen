package internal

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/f4biankoch/stackgen/internal"
)

func isStackgenProject(projectPath string) (fs.FS, error) {
	projectFS := os.DirFS(projectPath)

	_, err := projectFS.Open(internal.Manifest)
	if err != nil {
		return projectFS, fmt.Errorf("command has to be executed in valid stackgen project: %v", err)
	}

	return projectFS, nil
}
