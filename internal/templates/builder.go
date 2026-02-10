package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func BuildProjectFromTemplate(projectPath string, template fs.FS, options Options) error {
	if options.minimal {
		return buildMinimalProject(projectPath, template)
	}

	return buildProject(template, projectPath)
}

func buildMinimalProject(projectPath string, template fs.FS) error {
	fmt.Printf("\nBuilding minimal Project at %q:\n", projectPath)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	manifestFile, err := fs.ReadFile(template, "stackgen.json")
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(projectPath, "stackgen.json")
	return generateFile(manifestPath, manifestFile)
}

func buildProject(template fs.FS, projectPath string) error {
	fmt.Printf("\nBuilding Project at %q:\n", projectPath)

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
	fmt.Printf("Creating file: %q\n", path)

	// safety check (normally all directories should be created at this point)
	// overhead is acceptable
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	// truncates existing file content which is wanted
	return os.WriteFile(path, content, 0644)
}
