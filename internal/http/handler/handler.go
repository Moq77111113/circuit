package handler

import (
	"net/http"

	"github.com/moq77111113/circuit/internal/actions"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/sync"
)

// Authenticator is the interface required by the handler.
// Matches circuit.Authenticator via structural typing.
type Authenticator interface {
	Authenticate(r *http.Request) (*auth.Identity, error)
}

// noneAuth is a no-op authenticator that always succeeds.
type noneAuth struct{}

func (noneAuth) Authenticate(r *http.Request) (*auth.Identity, error) {
	return &auth.Identity{Subject: "anonymous"}, nil
}

// Handler serves the config UI over HTTP.
type Handler struct {
	schema        ast.Schema
	cfg           any
	path          string
	title         string
	brand         bool
	readOnly      bool
	store         *sync.Store
	authenticator Authenticator
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
	Authenticator Authenticator
	Actions       []actions.Def
}

// New creates a new HTTP handler for the config UI.
func New(c Config) *Handler {
	if c.Authenticator == nil {
		c.Authenticator = noneAuth{}
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
