# circuit

[![Go Report Card](https://goreportcard.com/badge/github.com/moq77111113/circuit)](https://goreportcard.com/report/github.com/moq77111113/circuit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Circuit is a runtime control surface for Go processes: modify configuration, trigger actions, see changes live. No DB, no infra, no magic.**

Stop SSH’ing into servers. Change config and trigger actions safely, in-process, from a simple web interface.

Circuit does not manage users, orchestrate multiple services, or act as a distributed system.

## Demo 

A demo is available at [https://circuit.up.railway.app/](https://circuit.up.railway.app/), use credentials: admin / admin.

## The Pitch

You have a Go service running somewhere. You want to tweak a log level, flip a feature flag, or update a rate limit.

Usually, you have two bad options:

1.  **The "YOLO" approach**: SSH in, `vim config.yaml`, `systemctl restart`. Hope it comes back up.
2.  **The "Enterprise" approach**: Build a full admin API, a React frontend, set up auth, deploy a separate service... and now you have two problems.

**Circuit is the third option.** It lives _inside_ your binary. It reads your existing config struct. It serves a tiny, safe web UI to control it.

It also lets you trigger safe, application-defined actions like restarting a worker or flushing caches.

## Why Circuit?

- **In-Process**: No sidecars, no agents, no external databases. It's just a library.
- **Minimal**: Zero dependencies. No npm, no webpack, no build steps. Just Go.
- **Safe**: It validates input based on your struct types. No more typos crashing production.
- **Live**: Changes persist to disk and trigger callbacks instantly.

## Authentication

Circuit supports three auth modes via `WithAuth()`. No auth by default.

**No Auth** - Omit `WithAuth()` for local dev or when behind external auth:
```go
ui, _ := circuit.From(&cfg, circuit.WithPath("config.yaml"))
```

**Basic Auth** - Simple username/password (plaintext or argon2id hash):
```go
ui, _ := circuit.From(&cfg,
    circuit.WithAuth(circuit.NewBasicAuth("admin", "secret")),
)
```

**Forward Auth** - Reads headers from OAuth2 Proxy, Traefik, Cloudflare Access:
```go
ui, _ := circuit.From(&cfg,
    circuit.WithAuth(circuit.NewForwardAuth("X-Forwarded-User", nil)),
)
```

Your proxy handles OAuth redirects. Circuit validates the forwarded headers.

## Install

```bash
go get github.com/moq77111113/circuit
```

## Quick Start

1.  **Define your config:**

    ```go
    type Config struct {
        LogLevel string `yaml:"log_level" circuit:"type:select,options:debug|info|error"`
        MaxConns int    `yaml:"max_conns" circuit:"type:number,min:1"`
        FeatureX bool   `yaml:"feature_x" circuit:"type:checkbox"`
    }
    ```

2.  **Serve it (with optional auth):**

    ```go
    func main() {
        var cfg Config

        // Create the circuit
        auth := circuit.NewBasicAuth("admin", "secret") // optional
        c, _ := circuit.From(&cfg,
            circuit.WithPath("config.yaml"),
            circuit.WithAuth(auth), // remove for no auth
            circuit.WithOnChange(func(e circuit.ChangeEvent) {
                logger.SetLevel(cfg.LogLevel)
                pool.Resize(cfg.MaxConns)
            }),
        )

        // Serve it on a private port
        http.ListenAndServe(":9090", c)
    }
    ```

3.  **Control it:**
    Open `http://localhost:9090`. Tweak values. Hit Save. Watch your app adapt.

## Options

Configure behavior via `From(cfg, options...)`:

**Essential:**
- `WithPath(path)` - **Required.** Sets YAML file path to load and watch
- `WithAuth(auth)` - Add authentication (Basic or Forward Auth)

**UI Customization:**
- `WithTitle(title)` - Customize page title (default: "Circuit")
- `WithBrand(false)` - Hide Circuit footer branding
- `WithReadOnly(true)` - View-only mode: disable all edits

**Behavior:**
- `WithAutoWatch(false)` - Disable file watching (changes won't auto-reload)
- `WithAutoApply(false)` - **Preview mode:** form submits show preview, call `handler.Apply()` to confirm
- `WithAutoSave(false)` - **Manual save:** changes stay in memory, call `handler.Save()` to persist

**Callbacks:**
- `WithOnChange(fn)` - Called after config changes (form submit, file change, or manual apply)
- `WithOnError(fn)` - Called when file watch or reload fails

**Advanced:**
- `WithSaveFunc(fn)` - Custom persistence (database, S3, etc.) instead of flat files
- `WithActions(actions...)` - Register executable action buttons (restarts, cache flushes, etc.)

**Manual Control (when AutoApply/AutoSave disabled):**
```go
handler, _ := circuit.From(&cfg, circuit.WithAutoApply(false))

// In your HTTP handler:
if err := handler.Apply(r.Form); err != nil {
    // preview approved, changes applied to memory
}
if err := handler.Save(); err != nil {
    // changes persisted to disk
}
```

## Actions

Add executable action buttons to trigger operations from the UI:

```go
restart := circuit.NewAction("restart", "Restart Worker", func(ctx context.Context) error {
    return worker.Restart(ctx)
}).Describe("Safely restarts the background worker").Confirm().WithTimeout(10 * time.Second)

flush := circuit.NewAction("flush", "Flush Cache", func(ctx context.Context) error {
    cache.Clear()
    return nil
}).Describe("Clears all cached data")

h, _ := circuit.From(&cfg,
    circuit.WithPath("config.yaml"),
    circuit.WithActions(restart, flush),
)
```

**Builder methods:**
- `Describe(text)` - Add description shown in UI
- `Confirm()` - Require confirmation dialog before execution
- `WithTimeout(duration)` - Set execution timeout (default: 30s)

## Features

- **Zero Setup**: One function call, no configuration files.
- **Auto-Generated**: Struct tags → beautiful UI.
- **Hot-Reload**: File changes reflect instantly in the UI.
- **Actions Callbacks**: UI buttons linked to code.
- **Validation**: Min/max/required constraints.
- **Read-only Fields**: Prevent accidental changes.
- **Apply Hook**: Rollback on error.
- **Zero Dependencies**: Embedded CSS/JS. No build step.
- **Framework Agnostic**: Works with stdlib, Echo, Gin, Fiber, Chi.

## Non-Goals

Circuit is strictly a single-process control surface. It explicitly avoids:

- Multi-service orchestration
- User management / DB
- Distributed system logic
- Arbitrary shell execution

## Roadmap

- ✅ Auth support (Basic, Forward)
- ✅ Action buttons
- ⏳ Validation + read-only fields
- ⏳ Apply hooks with rollback

## Contributing

PRs welcome. Keep it simple. Write tests. No "service" or "manager" files.

## Security

**Always use HTTPS in production.** HTTP Basic Auth sends credentials in base64 (easily decoded).

**Don't expose Circuit to the public internet without auth.** Use:

- Basic Auth for simple deployments
- Forward Auth behind OAuth2 Proxy/Traefik for production
- Network isolation (VPN, internal networks)

## License

MIT. Use responsibly.
