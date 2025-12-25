package sync

import (
	"fmt"
	"os"

	"github.com/moq77111113/circuit/internal/yaml"
)

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
