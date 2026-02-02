package projectInit

/*
	Cannot test this yet without a good template engine
*/

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateTemplatePath(templateName string) (string, error) {
	templatedir := filepath.Join("templates", templateName)

	info, err := os.Stat(templatedir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found (expected folder %s)", templateName, templatedir)
		}
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("template path %s is not a directory", templatedir)
	}

	return templatedir, nil
}
