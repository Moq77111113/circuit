package circuit

// SaveFunc is called to persist configuration changes.
// Receives the current config value and path, returns error if persistence fails.
type SaveFunc func(cfg any, path string) error

// Option configures behavior passed to `From`.
type Option func(*config)

type config struct {
	path          string
	title         string
	brand         bool
	readOnly      bool
	onChange      OnChange
	onError       func(error)
	autoReload    bool
	autoApply     bool
	autoSave      bool
	saveFunc      SaveFunc
	authenticator Authenticator
	actions       []Action
}

// WithPath sets the filesystem path to the configuration file.
//
// This option is REQUIRED - From will return an error if no path is provided.
//
// Supported formats are auto-detected by extension:
//   - .yaml, .yml - YAML format
//   - .json - JSON format
//   - .toml - TOML format
//
// Circuit will:
//  1. Load the initial config from this file on startup
//  2. Watch the file for changes (unless WithAutoWatch(false))
//  3. Persist updates to this file (unless WithAutoSave(false))
//
// If the file doesn't exist, From returns an error. Create the file first or
// use a custom SaveFunc to handle initialization.
func WithPath(path string) Option {
	return func(c *config) {
		c.path = path
	}
}

// WithTitle sets the title displayed in the UI header.
//
// If not provided, the UI displays "Configuration" as the default title.
//
// Example:
//
//	circuit.WithTitle("Production Settings")
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// WithBrand controls whether the Circuit footer/brand is shown in the UI.
//
// Default: true (brand is shown).
//
// Set to false to hide the "Powered by Circuit" footer:
//
//	circuit.WithBrand(false)
func WithBrand(b bool) Option {
	return func(c *config) {
		c.brand = b
	}
}

// WithOnChange registers a callback for configuration change events.
//
// The callback is invoked AFTER the in-memory config struct has been updated.
// Your application is responsible for applying the new config to running components.
//
// The ChangeEvent indicates the source of the change:
//   - SourceFormSubmit - user submitted the web form
//   - SourceFileChange - file changed on disk (file watcher)
//   - SourceManual - handler.Apply() was called directly
//
// Example:
//
//	circuit.WithOnChange(func(e circuit.ChangeEvent) {
//	    log.Printf("config updated from %s", e.Source)
//	    server.ApplyConfig(cfg) // Your responsibility
//	})
func WithOnChange(fn OnChange) Option {
	return func(c *config) {
		c.onChange = fn
	}
}

// WithOnError registers a callback for errors during file watching or reload.
//
// Common error scenarios:
//   - File watcher errors (permissions, inotify limits)
//   - Config parse errors (invalid YAML/JSON/TOML)
//   - File read errors (deleted file, network mount issues)
//
// The callback is invoked for non-fatal errors. Fatal errors (initial load failure)
// are returned by From() directly.
//
// Example:
//
//	circuit.WithOnError(func(err error) {
//	    log.Printf("config reload failed: %v", err)
//	})
func WithOnError(fn func(error)) Option {
	return func(c *config) {
		c.onError = fn
	}
}

// WithAutoWatch controls whether file watching and automatic reload are enabled.
//
// Default: true (file watching enabled).
//
// When true: Circuit watches the config file for changes and automatically reloads
// the in-memory struct when the file is modified.
//
// When false: File watching is disabled. Config is only reloaded when:
//   - User submits the web form
//   - handler.Apply() is called manually
//
// Disable auto-watch if you want manual control over when config is reloaded, or
// if you're running in an environment without inotify support.
func WithAutoWatch(enable bool) Option {
	return func(c *config) {
		c.autoReload = enable
	}
}

// WithAuth sets the authenticator for the Circuit UI.
//
// Default: nil (no authentication - UI is publicly accessible).
//
// IMPORTANT: Circuit UIs should always be protected. Editing config can be dangerous.
// Use authentication or place the handler behind a reverse proxy.
//
// Built-in authenticators:
//   - NewBasicAuth(username, password) - HTTP Basic Auth
//   - NewForwardAuth(header, claims) - Reverse proxy headers
//
// Example with Basic Auth:
//
//	auth := circuit.NewBasicAuth("admin", "$argon2id$v=19$...")
//	circuit.WithAuth(auth)
//
// Example with Forward Auth (OAuth2 Proxy):
//
//	auth := circuit.NewForwardAuth("X-Forwarded-User", nil)
//	circuit.WithAuth(auth)
//
// The authenticator is called on EVERY request to the handler.
func WithAuth(a Authenticator) Option {
	return func(c *config) {
		c.authenticator = a
	}
}

// WithAutoApply controls whether form submissions automatically update the
// in-memory config struct.
//
// Default: true (auto-apply enabled).
//
// When true: Form submission (POST) immediately updates the config struct in memory
// and saves to disk (if WithAutoSave(true)).
//
// When false: Form submission renders a preview with the submitted values but does
// NOT modify the config in memory. You must call handler.Apply(formData) manually
// to confirm the changes.
//
// Use preview mode (false) when you want to:
//   - Review changes before applying them
//   - Add custom validation before updates
//   - Implement approval workflows
//
// Example (preview mode):
//
//	h, _ := circuit.From(&cfg, circuit.WithAutoApply(false))
//	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
//	    if r.Method == "POST" && r.FormValue("confirm") == "yes" {
//	        r.ParseForm()
//	        h.Apply(r.Form) // Manually confirm
//	    }
//	    h.ServeHTTP(w, r)
//	})
func WithAutoApply(enable bool) Option {
	return func(c *config) {
		c.autoApply = enable
	}
}

// WithAutoSave controls whether config changes are automatically persisted to disk.
//
// Default: true (auto-save enabled).
//
// When true: After updating the in-memory config, Circuit automatically writes
// the new config to the file specified in WithPath.
//
// When false: The in-memory config is updated, but the file is NOT written.
// You must call handler.Save() manually to persist changes.
//
// Use manual save (false) when you want to:
//   - Batch multiple changes before writing to disk
//   - Add custom validation before persistence
//   - Control when disk I/O happens
//
// Example (manual save):
//
//	h, _ := circuit.From(&cfg, circuit.WithAutoSave(false))
//	// Later, after validating changes:
//	if err := h.Save(); err != nil {
//	    log.Printf("failed to save config: %v", err)
//	}
func WithAutoSave(enable bool) Option {
	return func(c *config) {
		c.autoSave = enable
	}
}

// WithSaveFunc replaces the default file writing with custom persistence logic.
//
// Default: Circuit writes to the file path using the detected format (YAML/JSON/TOML).
//
// Use a custom SaveFunc when you need to:
//   - Store config in a database or remote storage
//   - Add encryption or compression before writing
//   - Validate config before persisting
//   - Notify external systems about config changes
//
// The SaveFunc receives the current config value and path.
// It is called:
//   - After form submission (if WithAutoSave(true))
//   - When handler.Save() is called manually
//
// Example (custom persistence):
//
//	circuit.WithSaveFunc(func(cfg any, path string) error {
//	    // Validate before saving
//	    if err := validateConfig(cfg); err != nil {
//	        return err
//	    }
//	    // Write to database instead of file
//	    return db.SaveConfig(cfg)
//	})
func WithSaveFunc(fn SaveFunc) Option {
	return func(c *config) {
		c.saveFunc = fn
	}
}

// WithReadOnly makes the UI read-only, preventing all edits.
//
// Default: false (UI is editable).
//
// When true:
//   - All input fields are disabled
//   - Save button is hidden
//   - Add/Remove slice item buttons are hidden
//   - Form submission is blocked
//
// Use read-only mode when you want to:
//   - Expose config for viewing without allowing edits
//   - Provide a "status" page for operators
//   - Show config in production environments where edits must go through CI/CD
//
// Example:
//
//	circuit.WithReadOnly(true)
func WithReadOnly(enable bool) Option {
	return func(c *config) {
		c.readOnly = enable
	}
}

// WithActions registers executable server-side actions in the UI.
//
// Actions appear as buttons in the Actions section and allow operators to trigger
// safe, application-defined operations like restarting workers or flushing caches.
//
// Each action is created with NewAction(name, label, run) and can be configured via:
//   - .Describe(text) - add help text
//   - .Confirm() - require user confirmation before execution
//   - .WithTimeout(duration) - set execution timeout (default: 30s)
//
// Example:
//
//	restart := circuit.NewAction("restart", "Restart Worker", func(ctx context.Context) error {
//	    return worker.Restart(ctx)
//	}).Describe("Safely restarts the worker").Confirm().WithTimeout(10 * time.Second)
//
//	circuit.WithActions(restart)
//
// Safety notes:
//   - Actions run server-side code - ensure they are safe by default
//   - Use .Confirm() for destructive operations (restarts, deletions, flushes)
//   - Always use timeouts to prevent hanging operations
//   - Avoid shelling out unless necessary (prefer native Go APIs)
//   - Action failures are displayed in the UI
func WithActions(actions ...Action) Option {
	return func(c *config) {
		c.actions = actions
	}
}
