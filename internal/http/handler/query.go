package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/ast/path"
)

func extractFocusPath(r *http.Request) path.Path {
	focusParam := r.URL.Query().Get("focus")
	if focusParam == "" {
		return path.Root()
	}
	return path.ParsePath(focusParam)
}
