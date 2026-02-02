package templates

/*
	Cannot test this yet without a good template engine
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ResolveTemplatePath(templateName string) (string, error) {
	source := "embed"
	templatePath := templateName

	if sourcePath := strings.SplitN(templateName, ":", 2); len(sourcePath) > 1 {
		source = sourcePath[0]
		templatePath = sourcePath[1]
	}

	var template Resolver

	switch source {
	case "embed":
		template = EmbeddedTemplate{"templates"}

	case "local":
		template = LocalTemplate{}

	default:
		return "", fmt.Errorf("unknown template source %q", source)
	}

	return template.ResolveAbsPath(templatePath)
}

type Resolver interface {
	ResolveAbsPath(templateName string) (string, error)
}

type LocalTemplate struct {
}

func (lt LocalTemplate) ResolveAbsPath(templatePath string) (string, error) {
	templatedir, err := filepath.Abs(templatePath)
	if err != nil {
		return "", err
	}

	return resolveFullPath(templatedir)
}

type EmbeddedTemplate struct {
	templateRoot string
}

func (et EmbeddedTemplate) ResolveAbsPath(templateName string) (string, error) {
	templatedir := filepath.Join(et.templateRoot, templateName)

	return resolveFullPath(templatedir)
}

func resolveFullPath(templatedir string) (string, error) {
	info, err := os.Stat(templatedir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found (expected folder %s)", filepath.Base(templatedir), templatedir)
		}
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("template path %s is not a directory", templatedir)
	}

	return templatedir, nil
}
