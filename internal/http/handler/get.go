package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/actions"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
	"github.com/moq77111113/circuit/internal/ui/render"
)

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var values ast.ValuesByPath
	h.store.WithLock(func() {
		values = form.ExtractValues(h.cfg, h.schema)
	})

	focusPath := extractFocusPath(r)

	// Create RenderContext
	rc := render.NewRenderContext(&h.schema, values)
	rc.Focus = focusPath
	rc.ReadOnly = h.readOnly

	// Create PageContext
	pc := layout.NewPageContext(rc)
	pc.Title = h.title
	pc.Brand = h.brand
	pc.Actions = convertActions(h.actions)
	pc.ErrorMessage = r.URL.Query().Get("error")

	page := layout.Page(pc)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

func convertActions(actions []actions.Def) []layout.ActionButton {
	buttons := make([]layout.ActionButton, len(actions))
	for i, a := range actions {
		buttons[i] = layout.ActionButton{
			Name:                a.Name,
			Label:               a.Label,
			Description:         a.Description,
			RequireConfirmation: a.RequireConfirmation,
		}
	}
	return buttons
}
