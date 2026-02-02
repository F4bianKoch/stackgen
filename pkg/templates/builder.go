package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateTemplate(projectPath string, templatePath string, exists bool, force bool) error {
	if err := createProjectDir(projectPath, exists, force); err != nil {
		return err
	}

	if err := buildTemplate(templatePath, projectPath, force); err != nil {
		return err
	}

	return nil
}

func createProjectDir(projectPath string, exists bool, force bool) error {
	if exists && force {
		fmt.Printf("Using existing project directory at: %s\n", projectPath)
		return nil
	}

	fmt.Printf("Creating project directory at: %s\n", projectPath)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	return nil
}

func buildTemplate(templatePath string, projectPath string, force bool) error {
	dirEntries, err := os.ReadDir(templatePath)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			subDirPath := filepath.Join(projectPath, dirEntry.Name())
			subTemplatePath := filepath.Join(templatePath, dirEntry.Name())

			if err := os.Mkdir(subDirPath, 0755); err != nil {
				return err
			}

			if err := buildTemplate(subTemplatePath, subDirPath, force); err != nil {
				return err
			}

			continue
		}

		filePath := filepath.Join(templatePath, dirEntry.Name())
		if err := generateFile(filePath, projectPath, force); err != nil {
			return err
		}
	}

	return nil
}

func generateFile(filePath string, projectPath string, force bool) error {
	var content []byte

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fileName := filepath.Join(projectPath, filepath.Base(filePath))
	if err := createFile(fileName, content, force); err != nil {
		return nil
	}

	return nil
}

func createFile(path string, content []byte, force bool) error {
	fmt.Printf("Creating file: %s\n", path)

	fileAttributes := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if !force {
		fileAttributes |= os.O_EXCL
	}

	file, err := os.OpenFile(path, fileAttributes, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file %s already exists (use --force to overwrite)", path)
		}
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}
