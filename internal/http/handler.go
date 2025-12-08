package http

import (
	"net/http"
	"os"

	"github.com/moq77111113/circuit/internal/reload"
	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/ui"
	"github.com/moq77111113/circuit/internal/yaml"
)

// Handler serves the config UI over HTTP.
type Handler struct {
	schema schema.Schema
	cfg    any
	path   string
	title  string
	brand  bool
	loader *reload.Loader
}

// New creates a new HTTP handler for the config UI.
func New(schema schema.Schema, cfg any, path, title string, brand bool, loader *reload.Loader) *Handler {
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
		values = extractValues(h.cfg, h.schema)
	})

	page := ui.Page(h.schema, values, h.title, h.brand)

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

	h.loader.WithLock(func() {
		err = applyForm(h.cfg, h.schema, r.Form)
	})

	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	var data []byte
	h.loader.WithLock(func() {
		data, err = yaml.Encode(h.cfg)
	})

	if err != nil {
		http.Error(w, "Failed to encode config", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(h.path, data, 0644)
	if err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	h.get(w, r)
}
