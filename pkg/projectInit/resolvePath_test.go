package projectInit

import (
	"strings"
	"testing"
)

func TestResolvePath(t *testing.T) {
	projectName := "testProject"
	t.Run("resolve_path", func(t *testing.T) {
		path, err := resolvePath(projectName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.HasSuffix(path, projectName) {
			t.Fatalf("expected path to end with %q, got %q", projectName, path)
		}
	})
}
