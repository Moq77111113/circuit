package handler

import (
	"net/http"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/ui/layout"
	"github.com/moq77111113/circuit/internal/ui/render"
)

func (h *Handler) renderPreview(w http.ResponseWriter, r *http.Request) {
	values := form.ExtractValues(h.cfg, h.schema)

	for key, vals := range r.Form {
		if len(vals) > 0 && key != "action" {
			values[key] = vals[0]
		}
	}

	focusPath := extractFocusPath(r)

	// Create RenderContext
	rc := render.NewRenderContext(&h.schema, values)
	rc.Focus = focusPath
	rc.ReadOnly = h.readOnly

	// Create PageContext with preview banner
	pc := layout.NewPageContext(rc)
	pc.Title = h.title
	pc.Brand = h.brand
	pc.TopContent = []g.Node{previewBanner(r.Form)}

	page := layout.Page(pc)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render preview", http.StatusInternalServerError)
	}
}

func previewBanner(formData map[string][]string) g.Node {
	hiddenFields := []g.Node{}
	for key, vals := range formData {
		if key != "action" && len(vals) > 0 {
			hiddenFields = append(hiddenFields, h.Input(
				h.Type("hidden"),
				h.Name(key),
				h.Value(vals[0]),
			))
		}
	}

	return h.Div(
		h.Class("preview-banner"),
		h.Div(
			h.Class("preview-banner__content"),
			h.Div(
				h.Class("preview-banner__message"),
				h.Strong(g.Text("Preview Mode")),
				g.Text(" — Review your changes before applying"),
			),
			h.Form(
				h.Method("POST"),
				h.Class("preview-banner__actions"),
				g.Group(hiddenFields),
				h.Button(
					h.Type("submit"),
					h.Name("action"),
					h.Value("confirm"),
					h.Class("btn btn-primary"),
					g.Text("Confirm Changes"),
				),
				h.A(
					h.Href("?"),
					h.Class("btn btn-secondary"),
					g.Text("Cancel"),
				),
			),
		),
	)
}
