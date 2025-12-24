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

**Circuit is the third option.** It lives *inside* your binary. It reads your existing config struct. It serves a tiny, safe web UI to control it.

It also lets you trigger safe, application-defined actions like restarting a worker or flushing caches.

## Why Circuit?

*   **In-Process**: No sidecars, no agents, no external databases. It's just a library.
*   **Minimal**: Zero dependencies. No npm, no webpack, no build steps. Just Go.
*   **Safe**: It validates input based on your struct types. No more typos crashing production.
*   **Live**: Changes persist to disk and trigger callbacks instantly.

### Authentication

Authentication is external. Circuit relies on the request identity to authorize config edits and actions. It supports:
- **Forward Auth** (because you probably have one already)
- **Basic Auth** (because sometimes simple is best)
- **No Auth** (because you live on the edge)

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

2.  **Serve it:**
    ```go
    func main() {
        var cfg Config
        
        // Create the circuit
        c, _ := circuit.From(&cfg,
            circuit.WithPath("config.yaml"),
            circuit.OnApply(func() {
                // This runs when you hit "Save" in the UI
                logger.SetLevel(cfg.LogLevel)
                pool.Resize(cfg.MaxConns)
            }),
            circuit.Action("Restart worker", func() { 
                return restartWorker() 
            }),
        )

        // Serve it on a private port
        http.ListenAndServe(":9090", c)
    }
    ```

3.  **Control it:**
    Open `http://localhost:9090`. Tweak values. Hit Save. Watch your app adapt.

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

## Roadmap (v0.1)

- Auth support
- Validation + read-only fields
- Apply hooks
- Action buttons (in-process, safe)


## Contributing

PRs welcome. Keep it simple. Write tests. No "service" or "manager" files.

## License

MIT. Do whatever you want. Just don't blame me if you expose this to the internet without auth.
