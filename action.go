package circuit

import (
	"context"
	"time"
)

// Action defines an executable server-side operation displayed as a button in the Circuit UI.
//
// Actions enable operators to trigger safe, application-defined operations like:
//   - Restarting workers or background jobs
//   - Flushing caches or clearing queues
//   - Running maintenance tasks
//   - Triggering manual sync operations
//
// SAFETY REQUIREMENTS:
//   - Actions run server-side code - ensure they are safe by default
//   - Use .Confirm() for destructive operations (restarts, deletions, flushes)
//   - Always use timeouts to prevent hanging (default: 30s, can be customized)
//   - Avoid shelling out unless necessary - prefer native Go APIs
//   - Never use actions for privileged operations (system updates, user management)
//   - Validate inputs if actions accept parameters
//
// Actions are registered via WithActions() and appear in the Actions section of the UI.
//
// Example (creating actions with constructor):
//
//	// Safe action - no confirmation needed
//	flush := circuit.NewAction("flush_cache", "Flush Cache", func(ctx context.Context) error {
//	    cache.Clear()
//	    return nil
//	}).Describe("Clears all cached data")
//
//	// Destructive action - requires confirmation
//	restart := circuit.NewAction("restart_worker", "Restart Worker", func(ctx context.Context) error {
//	    return worker.Restart(ctx)
//	}).Describe("Safely restarts the background worker").Confirm().WithTimeout(10 * time.Second)
//
//	h, _ := circuit.From(&cfg, circuit.WithActions(flush, restart))
type Action struct {
	Name                string
	Label               string
	Description         string
	Run                 func(context.Context) error
	Timeout             time.Duration
	RequireConfirmation bool
}

// NewAction creates a new action with required fields.
//
// Parameters:
//   - name: Unique identifier (used in URLs and logs)
//   - label: Display name shown on the button
//   - run: Function to execute (receives context for cancellation/timeout)
//
// The run function MUST respect the provided context:
//   - Check ctx.Done() for long-running operations
//   - Return when context is cancelled
//   - Use context-aware APIs (http.NewRequestWithContext, db.QueryContext, etc.)
//
// Default timeout is 30 seconds. Use .WithTimeout() for longer operations.
//
// Optional configuration via fluent builder methods:
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
//
// Use for destructive or irreversible operations:
//   - Restarting services or workers
//   - Deleting data or clearing caches
//   - Flushing queues or buffers
//   - Triggering expensive operations
//
// When enabled, the operator must explicitly click "Confirm" before the action runs.
// This helps prevent accidental execution of dangerous operations.
//
// Example:
//
//	circuit.NewAction("restart", "Restart Service", restartFunc).Confirm()
func (a Action) Confirm() Action {
	a.RequireConfirmation = true
	return a
}

// WithTimeout sets custom execution timeout. Default is 30 seconds.
//
// Use for long-running operations that need more time:
//   - Database migrations or large queries
//   - Cache warming operations
//   - Batch processing tasks
//   - External API calls with slow response times
//
// The timeout is enforced via context cancellation. Your action's run function
// MUST respect the context and return promptly when ctx.Done() is signaled.
//
// Example:
//
//	circuit.NewAction("migrate", "Run Migration", migrateFunc).WithTimeout(5 * time.Minute)
func (a Action) WithTimeout(d time.Duration) Action {
	a.Timeout = d
	return a
}
