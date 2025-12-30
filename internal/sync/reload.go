package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/moq77111113/circuit/internal/codec"
)

func (s *Store) reload() {
	s.mu.RLock()
	shouldDebounce := time.Since(s.lastFormSubmit) < s.debounceWindow
	s.mu.RUnlock()

	if shouldDebounce {
		return
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if s.onError != nil {
			s.onError(fmt.Errorf("%w: %w", ErrAutoReloadRead, err))
		}
		return
	}

	cdc, err := codec.Detect(s.path)
	if err != nil {
		if s.onError != nil {
			s.onError(fmt.Errorf("%w: %w", ErrAutoReloadParse, err))
		}
		return
	}

	s.mu.Lock()
	err = cdc.Parse(data, s.cfg)
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

	cdc, err := codec.Detect(s.path)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}

	s.mu.Lock()
	err = cdc.Parse(data, s.cfg)
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
