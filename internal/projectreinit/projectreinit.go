package internal

import (
	"fmt"
	"os"

	"github.com/f4biankoch/stackgen/internal"
	"github.com/f4biankoch/stackgen/internal/templates"
)

func Run(defaults bool) error {
	if !confirm() {
		return nil
	}

	fmt.Printf("Reinitializing stackgen project...\n\n")

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	projectFS, err := isStackgenProject(projectPath)
	if err != nil {
		return err
	}

	metadata, err := templates.ResolveMetadata(projectFS, defaults)
	if err != nil {
		return err
	}

	templateFS, err := templates.ResolveTemplateFS(metadata.String("template_source"))
	if err != nil {
		return err
	}

	// Folder and File creation begins here!!!

	if err := templates.BuildProjectFromTemplate(projectPath, templateFS, metadata); err != nil {
		return err
	}

	fmt.Printf("\nReinitialized stack in: %q\n", projectPath)
	internal.PrintNextSteps(projectPath)

	return nil
}
