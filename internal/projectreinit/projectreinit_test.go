package internal

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsStackgenProject(t *testing.T) {
	tests := []struct {
		name         string
		createMarker func(t *testing.T, dir string)
		wantErr      bool
	}{
		{
			name: "regular manifest",
			createMarker: func(t *testing.T, dir string) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(dir, "stackgen.json"), []byte("{}"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{name: "missing manifest", wantErr: true},
		{
			name: "manifest is a directory",
			createMarker: func(t *testing.T, dir string) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(dir, "stackgen.json"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectPath := t.TempDir()
			if test.createMarker != nil {
				test.createMarker(t, projectPath)
			}

			projectFS, err := isStackgenProject(projectPath)
			if (err != nil) != test.wantErr {
				t.Fatalf("isStackgenProject() error = %v, wantErr %v", err, test.wantErr)
			}
			if projectFS == nil {
				t.Fatal("isStackgenProject() returned a nil filesystem")
			}
			if !test.wantErr {
				if _, err := fs.ReadFile(projectFS, "stackgen.json"); err != nil {
					t.Fatalf("returned filesystem cannot read manifest: %v", err)
				}
			}
		})
	}
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "short yes", input: "y\n", want: true},
		{name: "case insensitive yes", input: "YES\n", want: true},
		{name: "no", input: "n\n", want: false},
		{name: "empty defaults to no", input: "\n", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			withStdin(t, test.input)
			if got := confirm(); got != test.want {
				t.Fatalf("confirm() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRun_RebuildsProjectFromManifest(t *testing.T) {
	templateDir := t.TempDir()
	templateManifest := `{"template_source":"local:` + filepath.ToSlash(templateDir) + `","name":"default"}`
	writeTestFile(t, filepath.Join(templateDir, "stackgen.json"), templateManifest)
	writeTestFile(t, filepath.Join(templateDir, "config.txt"), `name={{ index . "name" }}`)

	projectDir := t.TempDir()
	projectManifest := `{"template_source":"local:` + filepath.ToSlash(templateDir) + `","name":"custom"}`
	writeTestFile(t, filepath.Join(projectDir, "stackgen.json"), projectManifest)
	writeTestFile(t, filepath.Join(projectDir, "config.txt"), "old")
	chdirForTest(t, projectDir)
	withStdin(t, "yes\n")

	if err := Run(true); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(projectDir, "config.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "name=custom" {
		t.Fatalf("config contents = %q, want metadata from project manifest", got)
	}
}

func TestRun_CancelledDoesNotRequireValidProject(t *testing.T) {
	chdirForTest(t, t.TempDir())
	withStdin(t, "no\n")

	if err := Run(true); err != nil {
		t.Fatalf("Run() after cancellation error = %v", err)
	}
}

func withStdin(t *testing.T, input string) {
	t.Helper()
	original := os.Stdin
	inputFile, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := inputFile.WriteString(input); err != nil {
		t.Fatal(err)
	}
	if _, err := inputFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	os.Stdin = inputFile
	t.Cleanup(func() {
		os.Stdin = original
		if err := inputFile.Close(); err != nil {
			t.Errorf("close test stdin: %v", err)
		}
	})
}

func chdirForTest(t *testing.T, dir string) {
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

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestIsStackgenProject_ErrorIncludesContext(t *testing.T) {
	_, err := isStackgenProject(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "valid stackgen project") {
		t.Fatalf("isStackgenProject() error = %v, want project context", err)
	}
}
