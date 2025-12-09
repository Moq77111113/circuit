package reload

import (
	"fmt"
	"os"
	"sync"

	"github.com/moq77111113/circuit/internal/yaml"
)

type Source uint8

const (
	SourceFormSubmit Source = iota
	SourceFileChange
	SourceManual
)

type ChangeEvent struct {
	Source Source
	Path   string
}

type OnChange func(ChangeEvent)

// Loader handles config loading and reloading.
type Loader struct {
	path     string
	cfg      any
	onChange OnChange
	watcher  *Watcher
	mu       sync.RWMutex
}

// Load reads a config file and optionally starts watching for changes.
// When the file changes, it reloads the config and calls onChange.
func Load(path string, cfg any, onChange OnChange, autoReload bool) (*Loader, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	err = yaml.Parse(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	loader := &Loader{
		path:     path,
		cfg:      cfg,
		onChange: onChange,
	}

	if autoReload {
		watcher, err := Watch(path, loader.reload)
		if err != nil {
			return nil, fmt.Errorf("watch config: %w", err)
		}
		loader.watcher = watcher
	}

	return loader, nil
}

// Stop stops watching the config file.
func (l *Loader) Stop() {
	if l.watcher != nil {
		l.watcher.Stop()
	}
}

// WithLock executes a function while holding a read lock on the config.
func (l *Loader) WithLock(fn func()) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	fn()
}

// EmitChange emits a change event with the given source.
func (l *Loader) EmitChange(source Source) {
	if l.onChange != nil {
		l.onChange(ChangeEvent{
			Source: source,
			Path:   l.path,
		})
	}
}
