package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var values path.ValuesByPath
	h.store.WithLock(func() {
		values = form.ExtractValues(h.cfg, h.schema)
	})

	focusParam := r.URL.Query().Get("focus")
	var focusPath path.Path
	if focusParam == "" {
		focusPath = path.Root()
	} else {
		focusPath = path.ParsePath(focusParam)
	}

	page := layout.Page(h.schema, values, h.title, h.brand, focusPath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
