package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func ValidateManifest(template fs.FS, defaults bool) (Options, error) {
	var options Options
	options.minimal = true

	manifestFile, err := fs.ReadFile(template, "stackgen.json")
	if err != nil {
		if os.IsNotExist(err) {
			return options, fmt.Errorf("stackgen.json does not exist in the template: %w", err)
		}

		return options, fmt.Errorf("unexpected Error while validating template: %w", err)
	}

	err = json.Unmarshal(manifestFile, &options)
	if err != nil {
		return options, fmt.Errorf("while parsing stackgen.json: %v", err)
	}

	var invalidOptions []string
	for name, option := range options.Options {
		if !option.Required {
			continue
		}

		if defaults && option.Default == nil {
			invalidOptions = append(invalidOptions, fmt.Sprintf("option %q does not have a default value", name))
		}

		if !defaults && option.Value == nil {
			invalidOptions = append(invalidOptions, fmt.Sprintf("value of option %q needs to be specified", name))
		}
	}

	if len(invalidOptions) > 0 {
		fmt.Printf("\nError:\n%v\n", strings.Join(invalidOptions, "\n"))
		return options, nil
	}

	options.minimal = false
	return options, nil
}
