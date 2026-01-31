package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var force bool
var template string
var projectNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
var initCmd = &cobra.Command{
	Use:   "init <project-name>",
	Short: "Initializes a new stackgen project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		return runInit(projectName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&force, "force", false, "overwrite existing files if the target directory is not empty")
	initCmd.Flags().StringVar(&template, "template", "basic", "specify project template to use (default: basic)")
}

func runInit(projectName string) error {
	fmt.Println("Initializing new stackgen project...")

	err := validateProjectName(projectName)
	if err != nil {
		return err
	}

	projectPath, err := resolvePath(projectName)
	if err != nil {
		return err
	}

	exists, err := validateTargetPath(projectPath, force)
	if err != nil {
		return err
	}

	templatePath, err := validateTemplatePath(template)
	if err != nil {
		return err
	}

	// Folder and File creation begins here!!!

	err = createProjectDir(projectPath, exists, force)
	if err != nil {
		return err
	}

	var content []byte

	content, err = os.ReadFile(filepath.Join(templatePath, "compose.yml"))
	if err != nil {
		return err
	}
	err = createFile(filepath.Join(projectPath, "compose.yml"), content, force)
	if err != nil {
		return err
	}

	content, err = os.ReadFile(filepath.Join(templatePath, ".env"))
	if err != nil {
		return err
	}
	err = createFile(filepath.Join(projectPath, ".env"), content, force)
	if err != nil {
		return err
	}

	content, err = os.ReadFile(filepath.Join(templatePath, "README.md"))
	if err != nil {
		return err
	}
	err = createFile(filepath.Join(projectPath, "README.md"), content, force)
	if err != nil {
		return err
	}

	fmt.Printf("\nInitialized stack in: %s\n", projectName)
	fmt.Printf("\nNext Steps: \n")
	fmt.Printf(" - cd %s\n", projectName)
	fmt.Printf(" - docker compose up -d\n")
	fmt.Printf(" - docker compose logs -f\n")

	return nil
}

func validateProjectName(projectName string) error {
	if projectName == "." || projectName == ".." {
		return fmt.Errorf("invalid project name: %s", projectName)
	}

	if strings.Contains(projectName, "/") || strings.Contains(projectName, "\\") {
		return fmt.Errorf("project name cannot contain path separators")
	}

	if strings.HasPrefix(projectName, "-") {
		return fmt.Errorf("project name cannot start with a hyphen")
	}

	if !projectNameRe.MatchString(projectName) {
		return fmt.Errorf("project name can only contain alphanumeric characters, hyphens, and underscores")
	}

	return nil
}

func resolvePath(projectName string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projectPath := filepath.Join(workingDir, projectName)
	return projectPath, nil
}

func validateTargetPath(projectPath string, force bool) (bool, error) {
	fileInfo, err := os.Stat(projectPath)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err == nil {
		if !fileInfo.IsDir() {
			return false, fmt.Errorf("a file with the name %s already exists", projectPath)
		}

		entries, err := os.ReadDir(projectPath)

		if err != nil {
			return true, fmt.Errorf("error reading target directory: %v", err)
		}

		if len(entries) == 0 {
			return true, nil
		}

		if force {
			return true, nil
		}

		return true, fmt.Errorf("%s already exists. Use --force to overwrite", projectPath)
	}

	return false, err
}

func validateTemplatePath(templateName string) (string, error) {
	templateDir := filepath.Join("templates", templateName)
	info, err := os.Stat(templateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found (expected folder %s)", templateName, templateDir)
		}
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("template path %s is not a directory", templateDir)
	}

	return templateDir, nil
}

func createProjectDir(projectPath string, exists bool, force bool) error {
	if exists && force {
		fmt.Printf("Using existing project directory at: %s\n", projectPath)
		return nil
	}

	fmt.Printf("Creating project directory at: %s\n", projectPath)

	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		return err
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
