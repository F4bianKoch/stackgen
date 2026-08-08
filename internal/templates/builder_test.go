package templates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestBuildProjectFromTemplate(t *testing.T) {
	templateFS := fstest.MapFS{
		"stackgen.json": {Data: []byte(`{"ignored":"template manifest contents"}`)},
		".env":          {Data: []byte("PORT={{ index . \"port\" }}\n")},
		"nested/config": {Data: []byte("name={{ index . \"name\" }}\n")},
	}
	metadata := Metadata{"name": "example&co", "port": 8080}
	projectPath := t.TempDir()

	if err := BuildProjectFromTemplate(projectPath, templateFS, metadata); err != nil {
		t.Fatalf("BuildProjectFromTemplate() error = %v", err)
	}

	assertFileContents(t, filepath.Join(projectPath, ".env"), "PORT=8080\n")
	assertFileContents(t, filepath.Join(projectPath, "nested/config"), "name=example&co\n")

	manifestData, err := os.ReadFile(filepath.Join(projectPath, "stackgen.json"))
	if err != nil {
		t.Fatalf("read generated manifest: %v", err)
	}
	var gotMetadata Metadata
	if err := json.Unmarshal(manifestData, &gotMetadata); err != nil {
		t.Fatalf("generated manifest is invalid JSON: %v", err)
	}
	if gotMetadata["name"] != "example&co" || gotMetadata["port"] != float64(8080) {
		t.Fatalf("generated manifest = %#v, want metadata values preserved", gotMetadata)
	}
}

func TestBuildProjectFromTemplate_OverwritesExistingFile(t *testing.T) {
	projectPath := t.TempDir()
	path := filepath.Join(projectPath, "config.txt")
	if err := os.WriteFile(path, []byte("old contents that are longer"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := BuildProjectFromTemplate(projectPath, fstest.MapFS{
		"config.txt": {Data: []byte("new")},
	}, Metadata{})
	if err != nil {
		t.Fatalf("BuildProjectFromTemplate() error = %v", err)
	}

	assertFileContents(t, path, "new")
}

func TestBuildProjectFromTemplate_ReturnsTemplateParseError(t *testing.T) {
	err := BuildProjectFromTemplate(t.TempDir(), fstest.MapFS{
		"broken.txt": {Data: []byte("{{")},
	}, Metadata{})
	if err == nil {
		t.Fatal("BuildProjectFromTemplate() error = nil, want malformed template error")
	}
	if !strings.Contains(err.Error(), `parse template "broken.txt"`) {
		t.Fatalf("BuildProjectFromTemplate() error = %q, want file context", err)
	}
}

func TestBuildProjectFromTemplate_ReturnsFilesystemError(t *testing.T) {
	projectPath := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectPath, "nested"), []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := BuildProjectFromTemplate(projectPath, fstest.MapFS{
		"nested/config": {Data: []byte("value")},
	}, Metadata{})
	if err == nil {
		t.Fatal("BuildProjectFromTemplate() error = nil, want parent path error")
	}
}

func assertFileContents(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s contents = %q, want %q", path, got, want)
	}
}
