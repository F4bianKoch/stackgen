package projectinit

import (
	"fmt"
	"regexp"

	"github.com/f4biankoch/stackgen/internal/templates"
)

var projectNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func Run(projectName string, force bool, templateName string, defaults bool) error {
	fmt.Println("Initializing new stackgen project...")

	if err := validateProjectName(projectName); err != nil {
		return err
	}

	projectPath, err := resolvePath(projectName)
	if err != nil {
		return err
	}

	if err := validateTargetPath(projectPath, force); err != nil {
		return err
	}

	templateFS, err := templates.ResolveTemplateFS(templateName)
	if err != nil {
		return err
	}

	// Folder and File creation begins here!!!

	if err := templates.BuildProjectFromTemplate(projectPath, templateFS, defaults); err != nil {
		return err
	}

	fmt.Printf("\nInitialized stack in: %s\n", projectName)
	fmt.Printf("\nNext Steps: \n")
	fmt.Printf(" - cd %s\n", projectName)
	fmt.Printf(" - docker compose up -d\n")
	fmt.Printf(" - docker compose logs -f\n")

	return nil
}
