package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var values ast.ValuesByPath
	h.store.WithLock(func() {
		values = form.ExtractValues(h.cfg, h.schema)
	})

	focusPath := extractFocusPath(r)

	page := layout.Page(h.schema, values, focusPath, layout.PageOptions{
		Title:    h.title,
		Brand:    h.brand,
		ReadOnly: h.readOnly,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
