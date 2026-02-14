package templates

import (
	"fmt"
	"testing"
)

func TestTemplateResolver(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"No source but valid embedded template name", "basic", false},
		{"embed", fmt.Sprint("embed:", "basic"), false},
		{"local", fmt.Sprint("local:", t.TempDir()), false},

		{"Empty template path", "", true},
		{"No source and invalid template name", "invalidTemplate", true},
		{"invalid embedded path", fmt.Sprint("embed:", "invalidTemplate"), true},
		{"invalid local path", fmt.Sprint("local:", t.TempDir(), "/invalid"), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ResolveTemplateFS(test.path)
			requireErr(t, err, test.wantErr)
		})
	}
}
