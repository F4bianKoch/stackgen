package projectinit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_BuildsEmbeddedTemplateWithDefaults(t *testing.T) {
	workingDir := t.TempDir()
	chdir(t, workingDir)

	if err := Run("example", false, "basic", true); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	assertProjectFile(t, filepath.Join(workingDir, "example", ".env"), strings.Join([]string{
		"POSTGRES_USER=stackgen",
		"POSTGRES_PASSWORD=stackgen",
		"POSTGRES_DB=stackgen",
		"PORT=8000",
	}, "\n")+"\n")
	assertProjectFileContains(t, filepath.Join(workingDir, "example", "stackgen.json"), `"template_source":"basic"`)
}

func TestRun_DoesNotModifyNonEmptyTargetWithoutForce(t *testing.T) {
	workingDir := t.TempDir()
	chdir(t, workingDir)
	projectPath := filepath.Join(workingDir, "example")
	if err := os.Mkdir(projectPath, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(projectPath, "keep.txt")
	if err := os.WriteFile(sentinel, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Run("example", false, "basic", true)
	if err == nil {
		t.Fatal("Run() error = nil, want existing target error")
	}
	assertProjectFile(t, sentinel, "keep")
	if _, err := os.Stat(filepath.Join(projectPath, "stackgen.json")); !os.IsNotExist(err) {
		t.Fatalf("manifest should not have been created, stat error = %v", err)
	}
}

func TestRun_RejectsInvalidNameBeforeCreatingFiles(t *testing.T) {
	workingDir := t.TempDir()
	chdir(t, workingDir)

	if err := Run("../escape", false, "basic", true); err == nil {
		t.Fatal("Run() error = nil, want invalid project name error")
	}
	entries, err := os.ReadDir(workingDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("working directory contains %d entries after validation failure", len(entries))
	}
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})
}

func assertProjectFile(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s contents = %q, want %q", path, got, want)
	}
}

func assertProjectFileContains(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if !strings.Contains(string(got), want) {
		t.Fatalf("%s contents = %q, want substring %q", path, got, want)
	}
}
