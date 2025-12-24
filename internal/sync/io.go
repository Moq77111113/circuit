package sync

import (
	"fmt"
	"os"

	"github.com/moq77111113/circuit/internal/yaml"
)

// Load reads a config file and optionally starts watching for changes.
// When the file changes, it reloads the config and calls onChange.
func Load(path string, cfg any, onChange OnChange, autoReload bool) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	err = yaml.Parse(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	r := &Store{
		path:     path,
		cfg:      cfg,
		onChange: onChange,
	}

	if autoReload {
		watcher, err := Watch(path, r.reload)
		if err != nil {
			return nil, fmt.Errorf("watch config: %w", err)
		}
		r.watcher = watcher
	}

	return r, nil
}

// reload is called by the watcher when the file changes.
// It silently ignores errors (TODO: fix in Phase 4).
func (s *Store) reload() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}

	s.mu.Lock()
	err = yaml.Parse(data, s.cfg)
	s.mu.Unlock()

	if err != nil {
		return
	}

	if s.onChange != nil {
		s.onChange(ChangeEvent{
			Source: SourceFileChange,
			Path:   s.path,
		})
	}
}

// Reload manually reloads the config from disk.
func (s *Store) Reload() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	s.mu.Lock()
	err = yaml.Parse(data, s.cfg)
	s.mu.Unlock()

	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	if s.onChange != nil {
		s.onChange(ChangeEvent{
			Source: SourceManual,
			Path:   s.path,
		})
	}

	return nil
}
