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

// extractHTTPBasePath returns the base path where the handler is mounted.
// For example, if the handler is at /config, this returns "/config".
func extractHTTPBasePath(r *http.Request) string {
	return r.URL.Path
}
