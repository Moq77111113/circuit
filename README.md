# circuit

[![Go Report Card](https://goreportcard.com/badge/github.com/moq77111113/circuit)](https://goreportcard.com/report/github.com/moq77111113/circuit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Circuit is a runtime control surface for Go processes: modify configuration, trigger actions, see changes live. No DB, no infra, no magic.**

Stop SSH’ing into servers. Change config and trigger actions safely, in-process, from a simple web interface.

Circuit does not manage users, orchestrate multiple services, or act as a distributed system.

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

Circuit protects your UI and actions with pluggable authentication. **Circuit never handles OAuth flows directly** - it validates requests.

### Three Modes

**1. No Auth (default)**

```go
ui, _ := circuit.From(&cfg, circuit.WithPath("config.yaml"))
// No WithAuth() = open access
```

Use for: Local development, internal networks, when behind other auth layers.

**2. Basic Auth (simple deployments)**

```go
ui, _ := circuit.From(&cfg,
    circuit.WithPath("config.yaml"),
    circuit.WithAuth(circuit.NewBasicAuth("admin", "$argon2id$v=19$m=65536,t=3,p=4$...")),
)
```

**Passwords:** Plaintext or argon2id hash.

**3. Forward Auth (when you already have a proxy)**

```go
auth := circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
    "email": "X-Forwarded-Email",
})
ui, _ := circuit.From(&cfg, circuit.WithAuth(auth))
```

Use for: **When you already have OAuth2 Proxy, Traefik ForwardAuth, or Cloudflare Access deployed.**

Circuit reads headers set by your proxy. **The proxy handles OAuth redirects, not Circuit.**

**Example setup with OAuth2 Proxy:**

```yaml
# oauth2-proxy.cfg
upstreams = ["http://your-app:8080"]
provider = "github"

# Your app
http.Handle("/admin", circuit.From(&cfg,
    circuit.WithAuth(circuit.NewForwardAuth("X-Forwarded-User", nil)),
))
```

The proxy intercepts requests, redirects to GitHub/Google, then forwards with headers.

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

See [godoc](https://pkg.go.dev/github.com/moq77111113/circuit) for full API.

**Control:**
- `WithAutoApply(false)` - Preview mode: require confirmation before applying changes
- `WithAutoSave(false)` - Manual save: call `ui.Save()` to persist
- `WithAutoWatch(false)` - Disable file watching

**Customization:**
- `WithSaveFunc(fn)` - Custom persistence (database, S3, API, etc.)
- `WithOnError(fn)` - Error callback for auto-reload/watch failures

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
- ⏳ Validation + read-only fields
- ⏳ Apply hooks with rollback
- ⏳ Action buttons

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
