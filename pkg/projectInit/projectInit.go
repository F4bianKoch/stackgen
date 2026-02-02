package projectInit

import (
	"fmt"
	"regexp"
)

var projectNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func Run(projectName string, force bool, template string) error {
	fmt.Println("Initializing new stackgen project...")

	if err := validateProjectName(projectName); err != nil {
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

	if err := createProjectDir(projectPath, exists, force); err != nil {
		return err
	}

	if err := buildTemplate(templatePath, projectPath, force); err != nil {
		return nil
	}

	fmt.Printf("\nInitialized stack in: %s\n", projectName)
	fmt.Printf("\nNext Steps: \n")
	fmt.Printf(" - cd %s\n", projectName)
	fmt.Printf(" - docker compose up -d\n")
	fmt.Printf(" - docker compose logs -f\n")

	return nil
}
