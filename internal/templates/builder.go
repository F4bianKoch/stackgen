package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
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
		templateFile := template.Must(template.ParseFS(templateFS, path))

		return renderTemplateToFile(projectFile, templateFile, metadata)
	})
}

func renderTemplateToFile(path string, template *template.Template, metadata Metadata) error {
	fmt.Printf("Creating file: %q\n", path)

	// safety check (normally all directories should be created at this point)
	// overhead is acceptable
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
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
