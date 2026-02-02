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
	var template Resolver

	if strings.HasPrefix(templateName, "local:") {
		template = LocalTemplate{}
		templateName = strings.Split(templateName, ":")[1]
	}

	if template == nil {
		if name := strings.Split(templateName, ":"); len(name) > 1 {
			templateName = name[1]
		}

		template = EmbeddedTemplate{"templates"}
	}

	return template.getTemplatePath(templateName)
}

type Resolver interface {
	getTemplatePath(templateName string) (string, error)
}

type LocalTemplate struct {
}

func (lt LocalTemplate) getTemplatePath(templatePath string) (string, error) {
	templatedir, err := filepath.Abs(templatePath)
	if err != nil {
		return "", err
	}

	return resolveFullPath(templatedir)
}

type EmbeddedTemplate struct {
	templateRoot string
}

func (et EmbeddedTemplate) getTemplatePath(templateName string) (string, error) {
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
