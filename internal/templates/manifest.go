package templates

type Options struct {
	minimal bool
	Options map[string]Option
}

type Option struct {
	Value    any
	Default  any
	Required bool
}
