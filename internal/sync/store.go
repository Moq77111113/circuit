package sync

import (
	"sync"
	"time"
)

// SaveFunc is called to persist config changes.
// When not provided, the default YAML file write is used.
type SaveFunc func(cfg any, path string) error
type Store struct {
	path     string
	cfg      any
	onChange OnChange
	watcher  *Watcher
	mu       sync.RWMutex

	autoApply bool
	autoSave  bool
	saveFunc  SaveFunc

	lastFormSubmit time.Time
	debounceWindow time.Duration
}

// Stop stops watching the config file.
func (s *Store) Stop() {
	if s.watcher != nil {
		s.watcher.Stop()
	}
}

// WithLock executes a function while holding a read lock on the config.
func (s *Store) WithLock(fn func()) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fn()
}

// EmitChange emits a change event with the given source.
func (s *Store) EmitChange(source Source) {
	if s.onChange != nil {
		s.onChange(ChangeEvent{
			Source: source,
			Path:   s.path,
		})
	}
}

// AutoApply returns whether POST automatically updates memory.
func (s *Store) AutoApply() bool {
	return s.autoApply
}

// AutoSave returns whether changes automatically persist to disk.
func (s *Store) AutoSave() bool {
	return s.autoSave
}

func (s *Store) MarkFormSubmit() {
	s.lastFormSubmit = time.Now()
}
