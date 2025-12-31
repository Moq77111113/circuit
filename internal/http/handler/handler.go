package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/actions"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/sync"
)

// Handler serves the config UI over HTTP.
type Handler struct {
	schema        ast.Schema
	cfg           any
	path          string
	title         string
	brand         bool
	readOnly      bool
	store         *sync.Store
	authenticator auth.Authenticator
	actions       []actions.Def
}

// Config holds configuration for creating a Handler.
type Config struct {
	Schema        ast.Schema
	Cfg           any
	Path          string
	Title         string
	Brand         bool
	ReadOnly      bool
	Store         *sync.Store
	Authenticator auth.Authenticator
	Actions       []actions.Def
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
		readOnly:      c.ReadOnly,
		store:         c.Store,
		authenticator: c.Authenticator,
		actions:       c.Actions,
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
