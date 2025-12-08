package reload

import (
	"fmt"
	"os"
	"sync"

	"github.com/moq77111113/circuit/internal/yaml"
)

// Loader handles config loading and reloading.
type Loader struct {
	path    string
	cfg     any
	onApply func()
	watcher *Watcher
	mu      sync.RWMutex
}

// Load reads a config file and starts watching for changes.
// When the file changes, it reloads the config and calls onApply.
func Load(path string, cfg any, onApply func()) (*Loader, error) {
	// Initial load
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	err = yaml.Parse(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	loader := &Loader{
		path:    path,
		cfg:     cfg,
		onApply: onApply,
	}

	// Start watching
	watcher, err := Watch(path, loader.reload)
	if err != nil {
		return nil, fmt.Errorf("watch config: %w", err)
	}

	loader.watcher = watcher

	return loader, nil
}

// Stop stops watching the config file.
func (l *Loader) Stop() {
	if l.watcher != nil {
		l.watcher.Stop()
	}
}

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

	if l.onApply != nil {
		l.onApply()
	}
}

// WithLock executes a function while holding a read lock on the config.
func (l *Loader) WithLock(fn func()) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	fn()
}
