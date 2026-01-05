package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
	"github.com/moq77111113/circuit/internal/ui/render"
	"github.com/moq77111113/circuit/internal/validation"
)

// renderWithErrors re-renders the form with validation errors and preserved form data.
func (h *Handler) renderWithErrors(w http.ResponseWriter, r *http.Request, result *validation.ValidationResult) {
	configValues := form.ExtractValues(h.cfg, h.schema)
	mergedValues := validation.MergeFormValues(h.schema.Nodes, configValues, r.Form)

	focusPath := extractFocusPath(r)
	httpBasePath := extractHTTPBasePath(r)

	rc := render.NewRenderContext(&h.schema, mergedValues)
	rc.Focus = focusPath
	rc.HTTPBasePath = httpBasePath
	rc.ReadOnly = h.readOnly
	rc.Errors = result

	pc := layout.NewPageContext(rc)
	pc.Title = h.title
	pc.Brand = h.brand

	page := layout.Page(pc)

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render form with errors", http.StatusInternalServerError)
	}
}
