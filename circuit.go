package circuit

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/http/handler"
	"github.com/moq77111113/circuit/internal/reload"
)

// From creates and returns an `http.Handler` that serves a small web UI for
// inspecting and editing a YAML-backed configuration value.
//
// From validates that `cfg` is a pointer to a struct (used to extract schema
// information from struct tags), applies any provided Option values and then
// attempts to load the initial configuration from the path supplied via
// `WithPath`. If successful it starts a file watcher that reloads the
// configuration on changes and returns a handler wired to that loader.
//
// Common errors:
//   - when `cfg` is not a pointer
//   - when no path is provided (use `WithPath`)
//   - when schema extraction, initial load, or watcher setup fails
func From(cfg any, opts ...Option) (http.Handler, error) {
	if reflect.TypeOf(cfg).Kind() != reflect.Pointer {
		return nil, fmt.Errorf("config must be a pointer")
	}

	conf := &config{
		brand:      true,
		autoReload: true,
	}
	for _, opt := range opts {
		opt(conf)
	}

	if conf.path == "" {
		return nil, fmt.Errorf("path is required (use WithPath)")
	}

	s, err := ast.Extract(cfg)
	if err != nil {
		return nil, fmt.Errorf("extract schema: %w", err)
	}

	loader, err := reload.Load(conf.path, cfg, conf.onChange, conf.autoReload)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	h := handler.New(s, cfg, conf.path, conf.title, conf.brand, loader)

	return h, nil
}
