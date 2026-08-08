package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/f4biankoch/stackgen/internal"
)

func BuildProjectFromTemplate(projectPath string, templateFS fs.FS, metadata Metadata) error {
	fmt.Printf("\nBuilding Project at %q:\n", projectPath)

	return fs.WalkDir(templateFS, ".", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			dirPath := filepath.Join(projectPath, path)
			return os.MkdirAll(dirPath, 0755)
		}

		projectFile := filepath.Join(projectPath, path)

		if path == internal.Manifest {
			return renderManifest(projectFile, metadata)
		}

		templateFile, err := template.ParseFS(templateFS, path)
		if err != nil {
			return fmt.Errorf("parse template %q: %w", path, err)
		}

		return renderFileFromTemplate(projectFile, templateFile, metadata)
	})
}

func renderManifest(path string, metadata Metadata) error {
	fmt.Printf("Creating Manifest: %q\n", path)

	if err := validateParentTree(path); err != nil {
		return err
	}

	content, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// truncates existing file content which is wanted
	return os.WriteFile(path, content, 0644)
}

func renderFileFromTemplate(path string, template *template.Template, metadata Metadata) error {
	fmt.Printf("Creating file: %q\n", path)

	if err := validateParentTree(path); err != nil {
		return err
	}

	// truncates existing file content which is wanted
	projectFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer projectFile.Close()

	return template.Execute(projectFile, metadata)
}

func validateParentTree(path string) error {
	// safety check (normally all directories should be created at this point)
	// overhead is acceptable
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return nil
}
