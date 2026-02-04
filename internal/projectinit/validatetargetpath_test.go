package projectinit

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var cases = []struct {
	force bool
}{
	{true},
	{false},
}

func TestValidateTargetPath_NonExistentPaths(t *testing.T) {
	projectPath := t.TempDir()

	for _, tc := range cases {
		t.Run(fmt.Sprintf("force=%t", tc.force), func(t *testing.T) {
			invalidPath := filepath.Join(projectPath, "invalidPath")
			err := validateTargetPath(invalidPath, tc.force)
			if err != nil {
				t.Fatalf("unexpected err=%v", err)
			}
		})
	}
}

// This test does not really work!
func TestValidateTargetPath_IsFile(t *testing.T) {
	projectWorkingDir := t.TempDir()
	projectDir := filepath.Join(projectWorkingDir, "testFile")

	_, err := os.Create(projectDir)
	if err != nil {
		t.Fatalf("unexpected err=%v", err)
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("force=%t", tc.force), func(t *testing.T) {
			err := validateTargetPath(projectDir, tc.force)

			if err == nil {
				t.Fatalf("path %v cannot be a File", projectDir)
			}
		})
	}
}

func TestValidateTargetPath_AlreadyExistsButEmpty(t *testing.T) {
	projectDir := t.TempDir()

	for _, tc := range cases {
		t.Run(fmt.Sprintf("force=%t", tc.force), func(t *testing.T) {
			err := validateTargetPath(projectDir, tc.force)

			if err != nil {
				t.Fatalf("err should be nil but err=%v", err)
			}
		})
	}
}

func TestValidateTargetPath_AlreadyExistsNotEmpty(t *testing.T) {
	projectDir := t.TempDir()
	_, err := os.Create(filepath.Join(projectDir, "testFile"))
	if err != nil {
		t.Fatalf("unexpected err=%v", err)
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("force=%t", tc.force), func(t *testing.T) {
			err := validateTargetPath(projectDir, tc.force)

			if tc.force && err != nil {
				t.Fatalf("err must be nil! force=%t, err=%v", tc.force, err)
			}

			if !tc.force && err == nil {
				t.Fatalf("err cannot be nil! force=%t, err=%v", tc.force, err)
			}
		})
	}
}
