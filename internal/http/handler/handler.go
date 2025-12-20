package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/http/action"
	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/reload"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ui/layout"
)

// Handler serves the config UI over HTTP.
type Handler struct {
	schema ast.Schema
	cfg    any
	path   string
	title  string
	brand  bool
	loader *reload.Loader
}

// New creates a new HTTP handler for the config UI.
func New(schema ast.Schema, cfg any, path, title string, brand bool, loader *reload.Loader) *Handler {
	return &Handler{
		schema: schema,
		cfg:    cfg,
		path:   path,
		title:  title,
		brand:  brand,
		loader: loader,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPost:
		h.post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) get(w http.ResponseWriter, _ *http.Request) {
	var values map[string]any
	h.loader.WithLock(func() {
		values = form.ExtractValues(h.cfg, h.schema)
	})

	page := layout.Page(h.schema, values, h.title, h.brand)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Render(w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	act := action.Parse(r.Form)

	switch act.Type {
	case action.ActionAdd:
		err = h.handleAdd(act.Field)
	case action.ActionRemove:
		err = h.handleRemove(act.Field, act.Index)
	case action.ActionSave:
		err = h.handleSave(r.Form)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
