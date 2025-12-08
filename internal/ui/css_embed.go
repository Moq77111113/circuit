package ui

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed css/*.css
var cssFS embed.FS

var DefaultCSS string

func init() {

	files, err := fs.Glob(cssFS, "css/*.css")
	if err != nil {
		return
	}
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

func SetCSS(s string) {
	DefaultCSS = s
}
