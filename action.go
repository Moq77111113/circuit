package circuit

import (
	"context"
	"time"
)

// Action defines an executable user action displayed as a button in the Circuit UI.
//
// Actions enable operators to trigger safe, application-defined operations like:
//   - Restarting workers or background jobs
//   - Flushing caches or clearing queues
//   - Running maintenance tasks
//   - Triggering manual sync operations
//
// Actions are registered via WithActions() and appear in the Actions section of the UI.
//
// Example (creating actions with constructor):
//
//	restart := circuit.NewAction("restart_worker", "Restart Worker", func(ctx context.Context) error {
//	    return worker.Restart(ctx)
//	}).Describe("Safely restarts the background worker").Confirm().WithTimeout(10 * time.Second)
//
//	flush := circuit.NewAction("flush_cache", "Flush Cache", func(ctx context.Context) error {
//	    cache.Clear()
//	    return nil
//	}).Describe("Clears all cached data")
//
//	h, _ := circuit.From(&cfg, circuit.WithActions(restart, flush))
type Action struct {
	Name                string
	Label               string
	Description         string
	Run                 func(context.Context) error
	Timeout             time.Duration
	RequireConfirmation bool
}

// NewAction creates a new action with required fields.
// Optional configuration can be added via fluent builder methods:
//
//	action := circuit.NewAction("restart", "Restart Worker", runFunc).
//	    Describe("Safely restarts the worker").
//	    Confirm().
//	    WithTimeout(10 * time.Second)
func NewAction(name, label string, run func(context.Context) error) Action {
	return Action{
		Name:    name,
		Label:   label,
		Run:     run,
		Timeout: 30 * time.Second, // default timeout
	}
}

// Describe sets the action description shown in the UI.
// Helps operators understand what the action does before triggering it.
func (a Action) Describe(desc string) Action {
	a.Description = desc
	return a
}

// Confirm enables a confirmation dialog before execution.
// Use for destructive or irreversible operations (restarts, deletions, flushes).
// The operator must explicitly confirm before the action runs.
func (a Action) Confirm() Action {
	a.RequireConfirmation = true
	return a
}

// WithTimeout sets custom execution timeout. Default is 30 seconds.
// Use for long-running operations that need more time (migrations, large flushes).
func (a Action) WithTimeout(d time.Duration) Action {
	a.Timeout = d
	return a
}
