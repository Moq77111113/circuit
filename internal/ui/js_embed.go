package ui

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed js/*.js
var jsFS embed.FS

var DefaultJS string

func init() {
	files, err := fs.Glob(jsFS, "js/*.js")
	if err != nil {
		return
	}
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
