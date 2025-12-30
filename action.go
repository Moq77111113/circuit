package circuit

import (
	"context"
	"time"
)

// Action defines an executable user action.
type Action struct {
	Name        string
	Label       string
	Description string
	Run         func(context.Context) error

	timeout             time.Duration
	requireConfirmation bool
}

// Describe sets the action description.
func (a Action) Describe(desc string) Action {
	a.Description = desc
	return a
}

// Confirm enables confirmation dialog before execution.
func (a Action) Confirm() Action {
	a.requireConfirmation = true
	return a
}

// Timeout sets custom execution timeout (default: 30s).
func (a Action) Timeout(d time.Duration) Action {
	a.timeout = d
	return a
}
