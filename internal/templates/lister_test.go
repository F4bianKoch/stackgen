package templates

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	workingDir := t.TempDir()
	templatesDir := filepath.Join(workingDir, "templates")
	if err := os.MkdirAll(filepath.Join(templatesDir, "basic"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(templatesDir, "api"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(templatesDir, "embed.go"), []byte("package templates"), 0o644); err != nil {
		t.Fatal(err)
	}
	chdirForLister(t, workingDir)

	output := captureStdout(t, func() {
		if err := ListTemplates(); err != nil {
			t.Fatalf("ListTemplates() error = %v", err)
		}
	})

	if !strings.Contains(output, "Available templates:") ||
		!strings.Contains(output, "  - api\n") ||
		!strings.Contains(output, "  - basic\n") {
		t.Fatalf("ListTemplates() output = %q, want directory names", output)
	}
	if strings.Contains(output, "embed.go") {
		t.Fatalf("ListTemplates() output = %q, file should not be listed", output)
	}
}

func TestListTemplates_MissingDirectory(t *testing.T) {
	chdirForLister(t, t.TempDir())
	err := ListTemplates()
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("ListTemplates() error = %v, want missing directory error", err)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	fn()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = original

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func chdirForLister(t *testing.T, dir string) {
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
