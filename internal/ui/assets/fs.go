package assets

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed css/*.css
var cssFS embed.FS

//go:embed js/*.js
var jsFS embed.FS

var (
	DefaultCSS string
	DefaultJS  string
)

func init() {
	// Load CSS
	files, err := fs.Glob(cssFS, "css/*.css")
	if err == nil {
		var b strings.Builder
		for _, f := range files {
			data, err := cssFS.ReadFile(f)
			if err != nil {
				continue
			}
			b.Write(data)
			b.WriteString("\n")
		}
		DefaultCSS = b.String()
	}

	// Load JS
	files, err = fs.Glob(jsFS, "js/*.js")
	if err == nil {
		var b strings.Builder
		for _, f := range files {
			data, err := jsFS.ReadFile(f)
			if err != nil {
				continue
			}
			b.Write(data)
			b.WriteString("\n")
		}
		DefaultJS = b.String()
	}
}
