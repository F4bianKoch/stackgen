package templates

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/f4biankoch/stackgen/internal"
)

func ResolveMetadata(template fs.FS, defaults bool) (Metadata, error) {
	var metadata Metadata

	manifestFile, err := resolveManifest(template)
	if err != nil {
		return metadata, err
	}

	metadata, err = manifestToJson(manifestFile, &metadata)
	if err != nil {
		return metadata, err
	}

	for option, value := range metadata {
		if option == "template_source" {
			continue
		}

		value = resolveValue(os.Stdin, os.Stdout, option, value, defaults)
		metadata[option] = value
	}

	return metadata, nil
}

func manifestToJson(manifestFile []byte, metadata *Metadata) (Metadata, error) {
	err := json.Unmarshal(manifestFile, metadata)
	if err != nil {
		return nil, fmt.Errorf("while parsing %q: %v", internal.Manifest, err)
	}

	return *metadata, nil
}

func resolveManifest(template fs.FS) ([]byte, error) {
	manifestFile, err := fs.ReadFile(template, internal.Manifest)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%q does not exist in the template: %w", internal.Manifest, err)
		}

		return nil, fmt.Errorf("unexpected Error while validating template: %w", err)
	}
	return manifestFile, nil
}

func resolveValue(in io.Reader, out io.Writer, option string, value any, defaults bool) any {
	if defaults {
		return value
	}

	var resolvedValue string
	fmt.Fprintf(out, "%s [%v]: ", option, value)

	_, err := fmt.Fscanln(in, &resolvedValue)
	if err != nil {
		resolvedValue = ""
	}

	if resolvedValue == "" {
		return value
	}

	return resolvedValue
}
