package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/ast/path"
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/http/action"
	"github.com/moq77111113/circuit/internal/http/form"
	"github.com/moq77111113/circuit/internal/sync"
	"github.com/moq77111113/circuit/internal/ui/layout"
)

// Handler serves the config UI over HTTP.
type Handler struct {
	schema        ast.Schema
	cfg           any
	path          string
	title         string
	brand         bool
	store         *sync.Store
	authenticator auth.Authenticator
}

// Config holds configuration for creating a Handler.
type Config struct {
	Schema        ast.Schema
	Cfg           any
	Path          string
	Title         string
	Brand         bool
	Store         *sync.Store
	Authenticator auth.Authenticator
}

// New creates a new HTTP handler for the config UI.
func New(c Config) *Handler {
	if c.Authenticator == nil {
		c.Authenticator = auth.None{}
	}
	return &Handler{
		schema:        c.Schema,
		cfg:           c.Cfg,
		path:          c.Path,
		title:         c.Title,
		brand:         c.Brand,
		store:         c.Store,
		authenticator: c.Authenticator,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := h.authenticator.Authenticate(r)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Circuit"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPost:
		h.post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	var values map[string]any
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
