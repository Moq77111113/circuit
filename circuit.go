package circuit

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/moq77111113/circuit/internal/actions"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/http/handler"
	"github.com/moq77111113/circuit/internal/sync"
)

// Action defines an executable user action.
type Action struct {
	Name        string
	Label       string
	Description string
	Run         func(context.Context) error

	timeout             time.Duration
	requireConfirmation bool
}

// Describe sets the action description.
func (a Action) Describe(desc string) Action {
	a.Description = desc
	return a
}

// Confirm enables confirmation dialog before execution.
func (a Action) Confirm() Action {
	a.requireConfirmation = true
	return a
}

// Timeout sets custom execution timeout (default: 30s).
func (a Action) Timeout(d time.Duration) Action {
	a.timeout = d
	return a
}

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
			sync.WithOnError(conf.onError),
			sync.WithAutoApply(conf.autoApply),
			sync.WithAutoSave(conf.autoSave),
			sync.WithSaveFunc(conf.saveFunc),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	internalActions := make([]actions.Def, len(conf.actions))
	for i, a := range conf.actions {
		internalActions[i] = actions.Def{
			Name:                a.Name,
			Label:               a.Label,
			Description:         a.Description,
			Run:                 a.Run,
			Timeout:             a.timeout,
			RequireConfirmation: a.requireConfirmation,
		}
	}

	h := handler.New(handler.Config{
		Schema:        s,
		Cfg:           cfg,
		Path:          conf.path,
		Title:         conf.title,
		Brand:         conf.brand,
		ReadOnly:      conf.readOnly,
		Store:         store,
		Authenticator: conf.authenticator,
		Actions:       internalActions,
	})

	return &Handler{h: h}, nil
}
