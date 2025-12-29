package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
	"github.com/moq77111113/circuit/internal/validation"
)

// renderWithErrors re-renders the form with validation errors and preserved form data.
func (h *Handler) renderWithErrors(w http.ResponseWriter, r *http.Request, result *validation.ValidationResult) {
	configValues := form.ExtractValues(h.cfg, h.schema)
	mergedValues := validation.MergeFormValues(h.schema.Nodes, configValues, r.Form)

	focusPath := extractFocusPath(r)

	page := layout.PageWithErrors(h.schema, mergedValues, focusPath, result, layout.PageOptions{
		Title:    h.title,
		Brand:    h.brand,
		ReadOnly: h.readOnly,
	})

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render form with errors", http.StatusInternalServerError)
	}
}
