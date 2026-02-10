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
	options.Minimal = true

	err := manifestToJson(template, &options)
	if err != nil {
		return options, err
	}

	err = validateOptions(options, defaults)
	if err != nil {
		return options, err
	}

	resolveOptionValues(options, defaults)

	options.Minimal = false
	return options, nil
}

func resolveOptionValues(options Options, defaults bool) {
	for key, option := range options.Options {
		option.Resolved_Value = option.Default

		if !defaults && option.Value != nil {
			option.Resolved_Value = option.Value
		}

		options.Options[key] = option
	}
}

func validateOptions(options Options, defaults bool) error {
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
		return fmt.Errorf("\n%+v\n", strings.Join(invalidOptions, "\n"))
	}

	return nil
}

func manifestToJson(template fs.FS, options *Options) error {
	manifestFile, err := fs.ReadFile(template, Manifest)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%q does not exist in the template: %w", Manifest, err)
		}

		return fmt.Errorf("unexpected Error while validating template: %w", err)
	}

	err = json.Unmarshal(manifestFile, options)
	if err != nil {
		return fmt.Errorf("while parsing %q: %v", Manifest, err)
	}

	return nil
}
