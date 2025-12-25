package sync

import (
	"fmt"
	"os"

	"github.com/moq77111113/circuit/internal/yaml"
)

// Save manually persists the current config to disk.
func (s *Store) Save() error {
	var data []byte
	var err error

	s.mu.RLock()
	data, err = yaml.Encode(s.cfg)
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	if s.saveFunc != nil {
		if err := s.saveFunc(s.cfg, s.path); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
	} else {
		if err := os.WriteFile(s.path, data, 0644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}
	}

	return nil
}
