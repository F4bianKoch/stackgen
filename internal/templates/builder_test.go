package templates

import (
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"
)

func TestBuildProjectFromTemplate(t *testing.T) {
	tests := []struct {
		name       string
		templateFS fstest.MapFS
		metadata   Metadata
		wantErr    bool
	}{
		{
			"Template with no files",
			fstest.MapFS{},
			Metadata{},
			false,
		},
		{
			"Template with Manifest",
			fstest.MapFS{
				"stackgen.json": {Data: []byte("{}")},
			},
			Metadata{},
			false,
		},
		{
			"Template with other files",
			fstest.MapFS{
				".env": {Data: []byte("")},
			},
			Metadata{},
			false,
		},
		{
			"Template with Manifest and other files",
			fstest.MapFS{
				"stackgen.json": {Data: []byte("{}")},
				".env":          {Data: []byte("")},
			},
			Metadata{},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			projectPath := t.TempDir()
			err := BuildProjectFromTemplate(projectPath, test.templateFS, test.metadata)
			requireErr(t, err, test.wantErr)
			validateResult(t, projectPath, test.templateFS)
		})
	}
}

func validateResult(t *testing.T, projectPath string, templateFS fstest.MapFS) {
	t.Helper()
	projectFS := os.DirFS(projectPath)

	fs.WalkDir(projectFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatalf("unexpeted error: %v", err)
		}
		if path == "." {
			return nil
		}

		file, err := fs.ReadFile(templateFS, path)

		if err != nil {
			t.Fatalf("couldn't find %q in project but is in template", path)
		}

		projectFile, err := fs.ReadFile(projectFS, path)
		if err != nil {
			t.Fatalf("unexpected error=%v", err)
		}

		if strings.Compare(string(file), string(projectFile)) != 0 {
			t.Fatalf("file %q contents did not match: wanted=%s\n\ngot=%s", path, projectFile, file)
		}

		return nil
	})
}
