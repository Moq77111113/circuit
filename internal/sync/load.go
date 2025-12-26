package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/moq77111113/circuit/internal/yaml"
)

// Config holds configuration for creating a Store.
type Config struct {
	Path       string
	Cfg        any
	AutoReload bool
	Options    []Option
}

// Load reads a config file and optionally starts watching for changes.
func Load(c Config) (*Store, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	err = yaml.Parse(data, c.Cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	s := &Store{
		path:           c.Path,
		cfg:            c.Cfg,
		autoApply:      true,
		autoSave:       true,
		debounceWindow: 500 * time.Millisecond,
	}

	for _, opt := range c.Options {
		opt(s)
	}

	if c.AutoReload {
		watcher, err := Watch(c.Path, s.reload, s.onError)
		if err != nil {
			return nil, fmt.Errorf("watch config: %w", err)
		}
		s.watcher = watcher
	}

	return s, nil
}
