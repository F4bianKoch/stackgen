package templates

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	embedded_templates "github.com/f4biankoch/stackgen/templates"
)

func ResolveTemplateFS(templatePath string) (fs.FS, error) {
	source := "embed"
	templateName := templatePath

	if sourcePath := strings.SplitN(templatePath, ":", 2); len(sourcePath) > 1 {
		source = sourcePath[0]
		templateName = sourcePath[1]
	}

	if templateName == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}

	switch source {
	case "embed":
		templateRoot, err := fs.Sub(embedded_templates.FS, templateName)
		if err != nil {
			return nil, fmt.Errorf("embedded template %q not found: %w", templateName, err)
		}

		if _, err := fs.Stat(templateRoot, "."); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("embedded template %q not found", templateName)
			}
			return nil, fmt.Errorf("cannot access embedded template %q: %w", templateName, err)
		}

		return templateRoot, nil

	case "local":
		info, err := os.Stat(templateName)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("local template %q not found", templateName)
			}

			return nil, fmt.Errorf("cannot access local template folder %q: %w", templateName, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("local template %q must be a directory", templateName)
		}

		return os.DirFS(templateName), nil

	default:
		return nil, fmt.Errorf("unknown template source %q", source)
	}
}
