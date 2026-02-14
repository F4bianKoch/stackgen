package templates

import (
	"testing"
	"testing/fstest"

	"github.com/f4biankoch/stackgen/internal"
)

func TestManifestResolver(t *testing.T) {
	tests := []struct {
		name     string
		template fstest.MapFS
		wantErr  bool
	}{
		{
			"Valid Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{}")},
			},
			false,
		},
		{
			"Valid minimal filled Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{\"template_source\": \"\"}")},
			},
			false,
		},
		{
			"Valid filled Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{\"template_source\": \"\", \"my_option\": \"value\"}")},
			},
			false,
		},
		{
			"Invalid Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("")},
			},
			true,
		},
		{
			"No Manifest",
			fstest.MapFS{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ResolveMetadata(test.template, true)
			requireErr(t, err, test.wantErr)
		})
	}
}
