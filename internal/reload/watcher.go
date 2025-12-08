package reload

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a file for changes.
type Watcher struct {
	watcher  *fsnotify.Watcher
	done     chan struct{}
	callback func()
}

// Watch starts watching a file and calls the callback when it changes.
func Watch(path string, callback func()) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create watcher: %w", err)
	}

	err = fw.Add(path)
	if err != nil {
		if closeErr := fw.Close(); closeErr != nil {
			return nil, fmt.Errorf("watch file: %w; close watcher: %w", err, closeErr)
		}
		return nil, fmt.Errorf("watch file: %w", err)
	}

	w := &Watcher{
		watcher:  fw,
		done:     make(chan struct{}),
		callback: callback,
	}

	go w.run()

	return w, nil
}

// Stop stops watching the file and cleans up resources.
func (w *Watcher) Stop() {
	close(w.done)
	if err := w.watcher.Close(); err != nil {
		_ = err
	}
}

func (w *Watcher) run() {
	for {
		select {
		case <-w.done:
			return
		case event := <-w.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.callback()
			}
		case <-w.watcher.Errors:
			// TODO: handle errors
		}
	}
}
