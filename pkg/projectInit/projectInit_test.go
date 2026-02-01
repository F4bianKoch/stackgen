package projectInit

import (
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"abc", true},
		{"abc-123", true},
		{"abc_123", true},
		{"bad name", false},
		{"bad*name", false},
		{"../escape", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateProjectName(tc.name)
			if (err == nil) != tc.ok {
				t.Fatalf("name=%q ok=%v err=%v", tc.name, tc.ok, err)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	projectName := "testProject"
	path, err := resolvePath(projectName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasSuffix(path, projectName) {
		t.Fatalf("expected path to end with %q, got %q", projectName, path)
	}
}
