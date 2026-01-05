# circuit

[![Go Report Card](https://goreportcard.com/badge/github.com/moq77111113/circuit)](https://goreportcard.com/report/github.com/moq77111113/circuit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Circuit is a runtime control surface for Go processes: modify configuration, trigger actions, see changes live. No DB, no infra, no magic.**

Stop SSH'ing into servers. Change config and trigger actions safely, in-process, from a simple web interface.

Circuit does not manage users, orchestrate multiple services, or act as a distributed system.

---

**[Try the demo](https://circuit.up.railway.app/)** (username: `admin`, password: `admin`)

---

## The Pitch

You have a Go service running somewhere. You want to tweak a log level, flip a feature flag, or update a rate limit.

Usually, you have two bad options:

1. **The "YOLO" approach**: SSH in, `vim config.yaml`, `systemctl restart`. Hope it comes back up.
2. **The "Enterprise" approach**: Build a full admin API, a React frontend, set up auth, deploy a separate service... and now you have two problems.

**Circuit is the third option.** It lives _inside_ your binary. It reads your existing config struct. It serves a tiny, safe web UI to control it.

```go
type Config struct {
    Workers  int    `yaml:"workers" circuit:"type:number,min:1,max:100"`
    LogLevel string `yaml:"log_level" circuit:"type:select,options:debug|info|warn|error"`
}

func main() {
    var cfg Config

    ui, _ := circuit.From(&cfg,
        circuit.WithPath("config.yaml"),
        circuit.WithAuth(circuit.NewBasicAuth("admin", "secret")),
        circuit.WithOnChange(func(e circuit.ChangeEvent) {
            pool.Resize(cfg.Workers)       // Apply the change
            logger.SetLevel(cfg.LogLevel)  // Your code, your rules
        }),
    )

    http.ListenAndServe(":9090", ui)
}
```

Now open `http://localhost:9090`. Change values. Hit Save. Your app reacts instantly.

It also lets you trigger safe, application-defined actions like restarting a worker or flushing caches.

## Why Circuit?

- **In-Process**: No sidecars, no agents, no external databases. It's just a library.
- **Minimal**: Zero dependencies. No npm, no webpack, no build steps. Just Go.
- **Safe**: It validates input based on your struct types. No more typos crashing production.
- **Live**: Changes persist to disk and trigger callbacks instantly.

## Common Use Cases

**Feature flags**: Toggle features without redeploying
```go
EnableBeta bool `circuit:"type:checkbox,help:Enable beta features"`
```

**Rate limits**: Adjust throttling on the fly
```go
RequestsPerSec int `circuit:"type:number,min:1,max:10000,help:Max requests/sec"`
```

**Worker pools**: Scale background workers dynamically
```go
Workers int `circuit:"type:number,min:1,max:100,help:Worker pool size"`
```

**Log levels**: Debug production without rebuilding
```go
LogLevel string `circuit:"type:select,options:debug=Debug;info=Info;warn=Warning;error=Error"`
```

**Maintenance mode**: Flip a switch to return 503
```go
Maintenance bool `circuit:"type:checkbox,help:Enable maintenance mode"`
```

## Authentication

Three auth modes. Pick what fits your setup.

**No auth** (local dev or behind a trusted proxy):
```go
ui, _ := circuit.From(&cfg, circuit.WithPath("config.yaml"))
```

**Basic Auth** (simple username/password):
```go
// Dev: plaintext
auth := circuit.NewBasicAuth("admin", "secret")

// Production: argon2id hash
auth := circuit.NewBasicAuth("admin", "$argon2id$v=19$m=65536,t=3,p=4$...")

ui, _ := circuit.From(&cfg, circuit.WithAuth(auth))
```

**Forward Auth** (OAuth2 Proxy, Traefik, Cloudflare Access):
```go
auth := circuit.NewForwardAuth("X-Forwarded-User", nil)
ui, _ := circuit.From(&cfg, circuit.WithAuth(auth))
```

Your reverse proxy handles OAuth. Circuit reads the headers.

## Quick Start

```bash
go get github.com/moq77111113/circuit
```

**1. Tag your config struct**

```go
type Config struct {
    Port     int    `yaml:"port" circuit:"type:number,min:1,max:65535,required"`
    LogLevel string `yaml:"log_level" circuit:"type:select,options:debug|info|error"`
    FeatureX bool   `yaml:"feature_x" circuit:"type:checkbox,help:Enable experimental feature"`
}
```

**2. Wire it up**

```go
func main() {
    var cfg Config

    ui, _ := circuit.From(&cfg,
        circuit.WithPath("config.yaml"),
        circuit.WithAuth(circuit.NewBasicAuth("admin", "secret")),
        circuit.WithOnChange(func(e circuit.ChangeEvent) {
            // Apply changes to your running app
            logger.SetLevel(cfg.LogLevel)
            server.UpdatePort(cfg.Port)
        }),
    )

    http.ListenAndServe(":9090", ui)
}
```

**3. Done**

Open `http://localhost:9090`. Change values. Hit Save. Your app reacts.

File changes (manual edits to `config.yaml`) also trigger the callback automatically.

## Options

Pass options to `From(cfg, options...)` to customize behavior:

| Option | What it does |
|--------|--------------|
| `WithPath(path)` | **Required.** Config file path (YAML/JSON/TOML auto-detected) |
| `WithAuth(auth)` | Enable authentication (Basic or Forward Auth) |
| `WithOnChange(fn)` | Callback fired after config changes (apply updates here) |
| `WithOnError(fn)` | Callback for file watch or reload errors |
| `WithTitle(title)` | Custom page title (default: "Configuration") |
| `WithReadOnly(true)` | View-only mode (no edits allowed) |
| `WithAutoWatch(false)` | Disable file watching (manual reload only) |
| `WithAutoApply(false)` | Preview mode: call `handler.Apply()` to confirm changes |
| `WithAutoSave(false)` | Manual save: call `handler.Save()` to persist |
| `WithSaveFunc(fn)` | Custom persistence (database, S3, etc.) |
| `WithActions(...)` | Add action buttons (see below) |
| `WithBrand(false)` | Hide Circuit footer |

**Preview mode example** (manual apply):
```go
h, _ := circuit.From(&cfg, circuit.WithAutoApply(false))

http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" && r.FormValue("confirm") == "yes" {
        r.ParseForm()
        h.Apply(r.Form) // Confirm the preview
    }
    h.ServeHTTP(w, r)
})
```

## Actions

Add buttons to trigger server-side operations: restart workers, flush caches, run migrations.

```go
// Safe operation - no confirmation needed
flush := circuit.NewAction("flush", "Flush Cache", func(ctx context.Context) error {
    cache.Clear()
    return nil
}).Describe("Clears all cached data")

// Destructive operation - requires confirmation
restart := circuit.NewAction("restart", "Restart Worker", func(ctx context.Context) error {
    return worker.Restart(ctx)
}).Describe("Safely restarts the background worker").Confirm().WithTimeout(10 * time.Second)

ui, _ := circuit.From(&cfg,
    circuit.WithPath("config.yaml"),
    circuit.WithActions(flush, restart),
)
```

**Action methods:**
- `.Describe(text)` – Help text shown in the UI
- `.Confirm()` – Require confirmation dialog (use for destructive ops)
- `.WithTimeout(duration)` – Execution timeout (default: 30s)

Actions run server-side with context cancellation. Failures are displayed in the UI.

## Struct Tag Reference

Circuit reads `circuit` tags to generate form fields:

```go
type Config struct {
    // Text input with validation
    Name string `circuit:"type:text,help:Service name,required,minlen:2,maxlen:50,pattern:^[a-z]+$"`

    // Number with range
    Port int `circuit:"type:number,min:1,max:65535,required"`

    // Select dropdown
    LogLevel string `circuit:"type:select,options:debug=Debug;info=Info;error=Error"`

    // Checkbox
    Enabled bool `circuit:"type:checkbox,help:Enable this feature"`

    // Hidden field (not shown in UI)
    Secret string `circuit:"-"`
}
```

**Input types:** `text`, `number`, `checkbox`, `select`, `password`, `email`, `url`, `date`, `time`, `color`

**Attributes:** `help`, `min`, `max`, `step`, `minlen`, `maxlen`, `pattern`, `options`, `required`, `readonly`

**Hide fields:** Use `circuit:"-"` to exclude sensitive data like API keys.

## What Circuit Doesn't Do

Circuit is a single-process control panel. It's not:

- A service mesh or orchestrator
- A user management system
- A distributed config store
- A shell command executor

If you need multi-service coordination, use Kubernetes ConfigMaps or Consul. Circuit is for the 90% case: one binary, one config file, quick edits.

## Security Notes

**Use HTTPS in production.** Basic Auth sends credentials in base64. Over HTTP, they're trivial to intercept.

**Protect the endpoint.** Don't expose Circuit on `0.0.0.0:80` without auth. Options:

- Enable `WithAuth()` with Basic or Forward Auth
- Put it behind a reverse proxy (Traefik, Caddy, nginx)
- Bind to `127.0.0.1` and access via SSH tunnel or VPN
- Run on a separate admin port and firewall it

**Use argon2id for passwords.** Never use plaintext passwords in production:
```go
// Bad: plaintext password
auth := circuit.NewBasicAuth("admin", "secret")

// Good: argon2id hash
auth := circuit.NewBasicAuth("admin", "$argon2id$v=19$m=65536,t=3,p=4$...")
```

Generate hashes with `golang.org/x/crypto/argon2`.

## Contributing

PRs welcome. Keep it minimal. Write tests. No "service" or "manager" files.


## License

MIT. Use responsibly.
