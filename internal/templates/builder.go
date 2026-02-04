package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func BuildTemplate(projectPath string, template fs.FS) error {
	return fs.WalkDir(template, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			dirPath := filepath.Join(projectPath, path)
			return os.MkdirAll(dirPath, 0755)
		}

		projectFile := filepath.Join(projectPath, path)
		templateFile, err := fs.ReadFile(template, path)
		if err != nil {
			return err
		}

		return generateFile(projectFile, templateFile)
	})
}

func generateFile(path string, content []byte) error {
	fmt.Printf("Creating file: %s\n", path)

	// safety check (normally all directories should be created at this point)
	// overhead is acceptable
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// truncates existing file content which is wanted
	return os.WriteFile(path, content, 0644)
}
