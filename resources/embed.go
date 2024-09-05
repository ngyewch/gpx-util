package resources

import "embed"

var (
	//go:embed *.gohtml
	TemplateFS embed.FS
)
