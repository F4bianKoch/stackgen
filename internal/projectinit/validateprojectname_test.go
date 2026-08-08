package projectinit

import (
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	cases := []struct {
		testName    string
		projectName string
		ok          bool
	}{
		{"letters", "abc", true},
		{"hyphen", "abc-123", true},
		{"underscore", "abc_123", true},
		{"space", "bad name", false},
		{"special character", "bad*name", false},
		{"forward slash", "../escape", false},
		{"backslash", `..\escape`, false},
		{"current directory", ".", false},
		{"parent directory", "..", false},
		{"leading hyphen", "-project", false},
		{"empty", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			err := validateProjectName(tc.projectName)
			if (err == nil) != tc.ok {
				t.Fatalf("validateProjectName(%q) error = %v, want valid %v", tc.projectName, err, tc.ok)
			}
		})
	}
}
