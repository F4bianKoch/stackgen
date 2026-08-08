package templates

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/f4biankoch/stackgen/internal"
)

func TestManifestResolver(t *testing.T) {
	tests := []struct {
		name     string
		template fstest.MapFS
		metadata Metadata
		wantErr  bool
	}{
		{
			"Valid Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{}")},
			},
			Metadata{},
			false,
		},
		{
			"Valid null Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("null")},
			},
			nil,
			false,
		},
		{
			"Valid minimal filled Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{\"template_source\": \"\"}")},
			},
			Metadata{"template_source": ""},
			false,
		},
		{
			"Valid filled Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("{\"template_source\": \"\", \"my_option\": \"value\"}")},
			},
			Metadata{"template_source": "", "my_option": "value"},
			false,
		},
		{
			"Invalid Manifest",
			fstest.MapFS{
				internal.Manifest: {Data: []byte("")},
			},
			Metadata{},
			true,
		},
		{
			"No Manifest",
			fstest.MapFS{},
			Metadata{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata, err := ResolveMetadata(test.template, true)

			requireErr(t, err, test.wantErr)
			if test.wantErr {
				return
			}

			if !reflect.DeepEqual(metadata, test.metadata) {
				t.Fatalf("got %#v, want %#v", metadata, test.metadata)
				return
			}
		})
	}
}

func TestResolveValue(t *testing.T) {
	tests := []struct {
		name     string
		option   string
		defaults bool
		value    any
		input    string
		expected any
	}{
		{"empty user input", "port", false, 8000, "\n", 8000},
		{"EOF", "port", false, 8000, "", 8000},
		{"override default", "port", false, 8000, "8001", "8001"},
		{"defaults true returns default", "port", true, 8000, "8001", 8000},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := strings.NewReader(test.input)
			out := new(bytes.Buffer)

			value := resolveValue(in, out, test.option, test.value, test.defaults)
			if value != test.expected {
				t.Fatalf("got=%v, want=%v", value, test.expected)
			}
		})
	}
}
