package circuit

import (
	"fmt"
	"net/http"
	"reflect"

	internalhttp "github.com/moq77111113/circuit/internal/http"
	"github.com/moq77111113/circuit/internal/reload"
	"github.com/moq77111113/circuit/internal/schema"
)

// From creates an HTTP handler for the config UI.
// cfg must be a pointer to a struct with circuit tags.
func From(cfg any, opts ...Option) (http.Handler, error) {
	if reflect.TypeOf(cfg).Kind() != reflect.Pointer {
		return nil, fmt.Errorf("config must be a pointer")
	}

	conf := &config{
		brand: true,
	}
	for _, opt := range opts {
		opt(conf)
	}

	if conf.path == "" {
		return nil, fmt.Errorf("path is required (use WithPath)")
	}

	s, err := schema.Extract(cfg)
	if err != nil {
		return nil, fmt.Errorf("extract schema: %w", err)
	}

	loader, err := reload.Load(conf.path, cfg, conf.onApply)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return internalhttp.New(s, cfg, conf.path, conf.title, conf.brand, loader), nil
}
