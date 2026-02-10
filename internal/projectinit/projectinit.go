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

	options, err := templates.ValidateManifest(templateFS, defaults)
	if err != nil {
		return err
	}

	options.Project_name = projectName
	options.Template_source = templateName

	// Folder and File creation begins here!!!

	if err := templates.BuildProjectFromTemplate(projectPath, templateFS, options); err != nil {
		return err
	}

	fmt.Printf("\nInitialized stack in: %q\n", projectPath)
	fmt.Printf("\nNext Steps: \n")
	fmt.Printf(" - cd %s\n", projectPath)
	fmt.Printf(" - docker compose up -d\n")
	fmt.Printf(" - docker compose logs -f\n")

	return nil
}
