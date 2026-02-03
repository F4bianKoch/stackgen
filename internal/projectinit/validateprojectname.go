package projectinit

import (
	"fmt"
	"strings"
)

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
