package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var values map[string]any
	h.store.WithLock(func() {
		values = form.ExtractValues(h.cfg, h.schema)
	})

	focusPath := extractFocusPath(r)

	page := layout.Page(h.schema, values, h.title, h.brand, focusPath)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
