package embedded_templates

import "embed"

//go:embed ** */.*
var FS embed.FS
