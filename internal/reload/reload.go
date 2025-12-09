package reload

import (
	"fmt"
	"os"

	"github.com/moq77111113/circuit/internal/yaml"
)

func (l *Loader) reload() {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return
	}

	l.mu.Lock()
	err = yaml.Parse(data, l.cfg)
	l.mu.Unlock()

	if err != nil {
		return
	}

	if l.onChange != nil {
		l.onChange(ChangeEvent{
			Source: SourceFileChange,
			Path:   l.path,
		})
	}
}

func (l *Loader) Reload() error {
	data, err := os.ReadFile(l.path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	l.mu.Lock()
	err = yaml.Parse(data, l.cfg)
	l.mu.Unlock()

	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	if l.onChange != nil {
		l.onChange(ChangeEvent{
			Source: SourceManual,
			Path:   l.path,
		})
	}

	return nil
}
