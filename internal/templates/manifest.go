package templates

type Manifest struct {
	Project_name string
	Description  string
	Version      string
	Template     string

	Options map[string]Option
}

type Option struct {
	Value    any
	Default  any
	Required bool
}
