package sync

import (
	"fmt"
	"os"

	"github.com/moq77111113/circuit/internal/codec"
)

// Save manually persists the current config to disk.
func (s *Store) Save() error {
	cdc, err := codec.Detect(s.path)
	if err != nil {
		return fmt.Errorf("detect format: %w", err)
	}

	var data []byte

	s.mu.RLock()
	data, err = cdc.Encode(s.cfg)
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
