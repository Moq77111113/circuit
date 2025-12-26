package sync

import "errors"

var (
	ErrAutoReloadRead  = errors.New("auto-reload read failed")
	ErrAutoReloadParse = errors.New("auto-reload parse failed")
	ErrWatcher         = errors.New("watcher error")
)
