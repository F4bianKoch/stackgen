package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func BuildProjectFromTemplate(projectPath string, template fs.FS, defaults bool) error {
	if err := validateTemplate(template, defaults); err != nil {
		return buildMinimalProject(projectPath, template)
	}

	return buildProject(template, projectPath)
}

func validateTemplate(template fs.FS, defaults bool) error {
	manifestFile, err := fs.ReadFile(template, "stackgen.json")
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("stackgen.json does not exist in the template: %w", err)
		}

		return fmt.Errorf("unexpected Error while validating template: %w", err)
	}

	var manifest Manifest
	err = json.Unmarshal(manifestFile, &manifest)
	if err != nil {
		return fmt.Errorf("while parsing stackgen.json: %v", err)
	}

	var invalidOptions []string
	for name, option := range manifest.Options {
		if !option.Required {
			continue
		}

		if defaults && option.Default == nil {
			invalidOptions = append(invalidOptions, fmt.Sprintf("option %q does not have a default value", name))
		}

		if !defaults && option.Value == nil {
			if option.Default == nil {
				invalidOptions = append(invalidOptions, fmt.Sprintf("value of option %q needs to be specified", name))
			}
		}
	}

	for _, invalidOption := range invalidOptions {
		return fmt.Errorf(invalidOption)
	}

	return nil
}

func buildMinimalProject(projectPath string, template fs.FS) error {
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
