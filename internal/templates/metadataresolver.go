package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

func ResolveMetadata(template fs.FS, defaults bool) (Metadata, error) {
	var metadata Metadata

	err := manifestToJson(template, &metadata)
	if err != nil {
		return metadata, err
	}

	for option, value := range metadata {
		value = resolveValue(option, value, defaults)

		metadata[option] = value
	}

	return metadata, nil
}

func manifestToJson(template fs.FS, metadata *Metadata) error {
	manifestFile, err := fs.ReadFile(template, Manifest)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%q does not exist in the template: %w", Manifest, err)
		}

		return fmt.Errorf("unexpected Error while validating template: %w", err)
	}

	err = json.Unmarshal(manifestFile, metadata)
	if err != nil {
		return fmt.Errorf("while parsing %q: %v", Manifest, err)
	}

	return nil
}

func resolveValue(option string, value any, defaults bool) any {
	if defaults {
		return value
	}

	var resolvedValue string
	fmt.Printf("Enter value for %q [default: %+v]: ", option, value)
	fmt.Scanln(&resolvedValue)

	if resolvedValue == "" {
		return value
	}

	return resolvedValue
}
