package templates

import "fmt"

type Options struct {
	Minimal         bool
	Project_name    string
	Template_source string
	Options         map[string]Option
}

type Option struct {
	Value    any
	Default  any
	Required bool

	Resolved_Value any
}

func (o Option) String() string {
	return fmt.Sprint(o.Resolved_Value)
}
