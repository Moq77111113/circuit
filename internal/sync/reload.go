package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/moq77111113/circuit/internal/yaml"
)

func (s *Store) reload() {
	if time.Since(s.lastFormSubmit) < s.debounceWindow {
		return
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if s.onError != nil {
			s.onError(fmt.Errorf("%w: %w", ErrAutoReloadRead, err))
		}
		return
	}

	s.mu.Lock()
	err = yaml.Parse(data, s.cfg)
	s.mu.Unlock()

	if err != nil {
		if s.onError != nil {
			s.onError(fmt.Errorf("%w: %w", ErrAutoReloadParse, err))
		}
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
