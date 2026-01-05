// Package circuit provides an embeddable HTTP UI for viewing and editing
// application configuration structs with automatic file persistence.
//
// # What Circuit Does
//
// Circuit generates a web UI from a Go config struct and provides:
//   - Automatic form generation from struct tags
//   - File-backed persistence (YAML, JSON, TOML supported)
//   - Automatic file watching and in-memory reload
//   - Optional authentication (Basic Auth, Forward Auth)
//   - Optional executable actions (restart, reload, flush, etc.)
//
// # What Circuit Doesn't Do
//
// Circuit does not automatically restart your application or apply changes to
// running components. Your application is responsible for:
//   - Restarting services when config changes
//   - Applying new config values to active components
//   - Validating config changes before applying them
//
// Use the WithOnChange callback to hook into config updates and trigger your
// application's reload logic.
//
// # Basic Usage
//
// The typical workflow is:
//  1. Define a config struct with yaml and circuit tags
//  2. Load config from disk (your responsibility)
//  3. Create handler via circuit.From(&cfg, options...)
//  4. Mount handler on your existing http.ServeMux
//
// Example (minimal embed with net/http):
//
//	package main
//
//	import (
//	    "log"
//	    "net/http"
//	    "github.com/moq77111113/circuit"
//	)
//
//	type Config struct {
//	    Host string `yaml:"host" circuit:"type:text,help:Server hostname"`
//	    Port int    `yaml:"port" circuit:"type:number,help:Server port,required,min:1,max:65535"`
//	    TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`
//	}
//
//	func main() {
//	    var cfg Config
//
//	    // Create handler - Circuit will load initial config from file
//	    handler, err := circuit.From(&cfg,
//	        circuit.WithPath("config.yaml"),
//	        circuit.WithTitle("My App Settings"),
//	    )
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Mount on existing mux
//	    mux := http.NewServeMux()
//	    mux.Handle("/config", handler)
//
//	    log.Println("Circuit UI available at http://localhost:8080/config")
//	    http.ListenAndServe(":8080", mux)
//	}
//
// # Struct Tags
//
// Circuit uses the circuit struct tag to configure form fields. Use circuit:"-" to
// hide sensitive fields like passwords or API keys.
//
// Tag format: circuit:"type:INPUT_TYPE,ATTRIBUTE:VALUE,FLAG,..."
//
// Supported input types:
//   - text, password, email, url, tel
//   - number, range
//   - checkbox
//   - select (requires options attribute)
//   - date, time, color
//
// Common attributes:
//   - help:TEXT - help text shown below the field
//   - min:N, max:N, step:N - numeric constraints
//   - minlen:N, maxlen:N - string length constraints
//   - pattern:REGEX - regex validation pattern
//   - options:k1=v1;k2=v2 - select/radio options
//
// Common flags:
//   - required - field must not be empty
//   - readonly - field cannot be edited
//
// Example (struct tags):
//
//	type Config struct {
//	    // Text input with length constraints
//	    Name string `yaml:"name" circuit:"type:text,help:Service name,required,minlen:2,maxlen:50"`
//
//	    // Number input with range constraints
//	    Port int `yaml:"port" circuit:"type:number,help:Listen port,required,min:1,max:65535"`
//
//	    // Select input with options
//	    LogLevel string `yaml:"log_level" circuit:"type:select,options:debug=Debug;info=Info;warn=Warning;error=Error"`
//
//	    // Hidden field - Circuit will ignore this
//	    APIKey string `yaml:"api_key" circuit:"-"`
//	}
//
// # Security
//
// Circuit UIs should be protected. Editing config can be dangerous - protect the
// endpoint with authentication or place it behind a reverse proxy.
//
// Use WithAuth() to enable authentication:
//
//	// Development: plaintext password (DO NOT use in production)
//	auth := circuit.NewBasicAuth("admin", "dev-password")
//
//	// Production: argon2id hash
//	auth := circuit.NewBasicAuth("admin", "$argon2id$v=19$m=65536,t=3,p=4$...")
//
//	handler, _ := circuit.From(&cfg,
//	    circuit.WithPath("config.yaml"),
//	    circuit.WithAuth(auth),
//	)
//
// For reverse proxy setups (OAuth2 Proxy, Traefik ForwardAuth, Cloudflare Access):
//
//	auth := circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
//	    "email": "X-Forwarded-Email",
//	})
//
// # Actions
//
// Actions enable operators to trigger safe, application-defined operations like
// restarting workers or flushing caches.
//
// Actions are registered via WithActions() and appear as buttons in the UI:
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
//
// Safety notes for actions:
//   - Actions run server-side code - ensure they are safe by default
//   - Use .Confirm() for destructive operations (restarts, deletions)
//   - Always use timeouts - default is 30 seconds
//   - Avoid shelling out unless necessary (prefer native Go APIs)
//
// # File Watching and Hot Reload
//
// Circuit automatically watches the config file and reloads the in-memory struct
// when changes are detected. Use WithOnChange to be notified:
//
//	handler, _ := circuit.From(&cfg,
//	    circuit.WithPath("config.yaml"),
//	    circuit.WithOnChange(func(e circuit.ChangeEvent) {
//	        log.Printf("config reloaded from %s", e.Source)
//	        // Your responsibility: apply new config to running components
//	        server.ApplyConfig(cfg)
//	    }),
//	)
//
// Change sources:
//   - SourceFormSubmit - user submitted the form
//   - SourceFileChange - file changed on disk
//   - SourceManual - handler.Apply() was called directly
//
// Disable file watching with WithAutoWatch(false) if you want manual reload only.
package circuit
