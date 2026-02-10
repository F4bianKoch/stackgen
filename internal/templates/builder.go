package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

func BuildProjectFromTemplate(projectPath string, template fs.FS, options Options) error {
	if options.Minimal {
		return buildMinimalProject(projectPath, template, options)
	}

	return buildProject(projectPath, template, options)
}

func buildMinimalProject(projectPath string, templateFS fs.FS, options Options) error {
	fmt.Printf("\nBuilding minimal Project at %q:\n", projectPath)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	manifestPath := filepath.Join(projectPath, Manifest)
	manifestTemplate := template.Must(template.ParseFS(templateFS, Manifest))

	return renderTemplateToFile(manifestPath, manifestTemplate, options)
}

func buildProject(projectPath string, templateFS fs.FS, options Options) error {
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
		if err != nil {
			return err
		}

		return renderTemplateToFile(projectFile, templateFile, options)
	})
}

func renderTemplateToFile(path string, template *template.Template, options Options) error {
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

	return template.Execute(projectFile, options)
}
