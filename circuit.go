package circuit

import (
	"fmt"
	"reflect"

	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/http/handler"
	"github.com/moq77111113/circuit/internal/sync"
)

// From creates and returns a Handler that serves a small web UI for
// inspecting and editing a YAML-backed configuration value.
//
// From validates that `cfg` is a pointer to a struct (used to extract schema
// information from struct tags), applies any provided Option values and then
// attempts to load the initial configuration from the path supplied via
// `WithPath`. If successful it starts a file watcher that reloads the
// configuration on changes and returns a handler wired to that loader.
//
// The returned Handler implements http.Handler and exposes manual control methods:
//   - Apply(formData) - manually apply form data (for preview mode)
//   - Save() - manually save config to disk
//
// Common errors:
//   - when `cfg` is not a pointer
//   - when no path is provided (use `WithPath`)
//   - when schema extraction, initial load, or watcher setup fails
func From(cfg any, opts ...Option) (*Handler, error) {
	if reflect.TypeOf(cfg).Kind() != reflect.Pointer {
		return nil, fmt.Errorf("config must be a pointer")
	}

	conf := &config{
		brand:      true,
		autoReload: true,
		autoApply:  true,
		autoSave:   true,
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

	store, err := sync.Load(sync.Config{
		Path:       conf.path,
		Cfg:        cfg,
		AutoReload: conf.autoReload,
		Options: []sync.Option{
			sync.WithOnChange(conf.onChange),
			sync.WithAutoApply(conf.autoApply),
			sync.WithAutoSave(conf.autoSave),
			sync.WithSaveFunc(conf.saveFunc),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	h := handler.New(handler.Config{
		Schema:        s,
		Cfg:           cfg,
		Path:          conf.path,
		Title:         conf.title,
		Brand:         conf.brand,
		Store:         store,
		Authenticator: conf.authenticator,
	})

	return &Handler{h: h}, nil
}
