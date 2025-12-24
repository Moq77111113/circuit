package sync

import "sync"

// Store handles config loading, watching, and reloading.
// It provides thread-safe access to the current configuration and
// notifies listeners when the underlying config file changes.
type Store struct {
	path     string
	cfg      any
	onChange OnChange
	watcher  *Watcher
	mu       sync.RWMutex
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
