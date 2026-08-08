package projectinit

import (
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	workingDir := t.TempDir()
	chdir(t, workingDir)
	projectName := "testProject"
	path, err := resolvePath(projectName)
	if err != nil {
		t.Fatalf("resolvePath() error = %v", err)
	}

	want := filepath.Join(workingDir, projectName)
	if path != want {
		t.Fatalf("resolvePath() = %q, want %q", path, want)
	}
}
