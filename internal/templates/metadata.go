package templates

import "fmt"

type Metadata map[string]any

func (m Metadata) String(key string) string {
	return fmt.Sprintf("%v", m[key])
}
