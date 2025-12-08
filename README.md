# circuit

[![Go Report Card](https://goreportcard.com/badge/github.com/moq77111113/circuit)](https://goreportcard.com/report/github.com/moq77111113/circuit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Convert any Go config struct into a web dashboard in 30 seconds.**

**circuit** transforms your configuration structs into production-ready admin interfaces. Zero JavaScript. Zero build step. Zero infrastructure.

Perfect for **microservices**, **IoT devices**, **internal tools**, **sidecars**, **agents**, and **CLI utilities** that need runtime configuration without the overhead.

## Why?

Your Go services need runtime configuration. Traditional approaches mean:
- **SSH access** - Security nightmare for operators
- **Environment variables** - Requires restarts, no validation
- **Config management tools** - Overkill for simple services
- **Custom admin panels** - Weeks of React/Vue boilerplate

**circuit** is different. Add one line:

```go
ui, _ := circuit.From(&cfg, circuit.WithPath("config.yaml"))
```

You now have a production-grade web UI. Deploy to **Kubernetes sidecars**, **edge devices**, **Lambda functions**, **systemd services** - anywhere Go runs.

## Use Cases

- **Microservices**: Feature flags, rate limits, circuit breakers without redeployment
- **IoT/Edge**: Configure sensor thresholds, network settings, update intervals on field devices
- **Internal Tools**: Database connection strings, API keys, integration endpoints
- **Agents/Sidecars**: Monitoring configs, log levels, metric collectors
- **CLIs**: Persistent settings, user preferences, default values
- **Development**: Quick admin panels for prototypes and demos

## Features

- **30-second setup** - One function call, no configuration files
- **Auto-generated forms** - Struct tags → beautiful UI
- **Hot-reload** - File changes reflect instantly
- **Zero dependencies** - Embedded CSS, no npm, no webpack
- **Framework agnostic** - Works with stdlib, Echo, Gin, Fiber, Chi
- **Production ready** - Use it in real industrial systems and cloud deployments
- **Tiny footprint** - Perfect for constrained environments (IoT, edge, containers)

## Install

```bash
go get github.com/moq77111113/circuit
```

## Usage

### 1. Tag your config struct

```go
type Config struct {
    Host string `yaml:"host" circuit:"type:text,help:Server hostname"`
    Port int    `yaml:"port" circuit:"type:number,help:Server port,required"`
    TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`
}
```

### 2. Create the UI

```go
var cfg Config

ui, err := circuit.From(&cfg,
    circuit.WithPath("config.yaml"),
    circuit.WithTitle("My App Settings"),
)
if err != nil {
    panic(err)
}

http.ListenAndServe(":8080", ui)
```

### 3. Visit `http://localhost:8080`

Edit. Save. Done. The file updates, your app reloads (if you want).

## Options

```go
circuit.From(&cfg,
    circuit.WithPath("config.yaml"),        // Required: path to YAML file
    circuit.WithTitle("Admin Panel"),       // Optional: page title
    circuit.OnApply(func() {                // Optional: callback on save
        log.Println("Config updated!")
        // Reload services, reconnect clients, etc.
    }),
)
```

## Supported Types

- `text` - string input
- `number` - int/float input
- `checkbox` - boolean
- `password` - hidden text
- `select` - dropdown (define options in tag)
- `range` - slider
- `date`, `time` - date/time pickers
- `radio` - radio buttons

## Real-World Examples

### Microservice Feature Flags
```go
type Config struct {
    EnableNewAPI    bool   `circuit:"type:checkbox,help:Enable v2 API endpoints"`
    RateLimit       int    `circuit:"type:range,min:10,max:1000,help:Requests per minute"`
    CacheTimeout    int    `circuit:"type:number,help:Cache TTL in seconds"`
}
```

### IoT Device Configuration
```go
type SensorConfig struct {
    SampleRate    int     `circuit:"type:range,min:100,max:10000,help:Sampling rate (ms)"`
    Threshold     float64 `circuit:"type:number,help:Alert threshold"`
    ServerURL     string  `circuit:"type:text,help:Data upload endpoint"`
    EnableOffline bool    `circuit:"type:checkbox,help:Store data offline"`
}
```

### Kubernetes Sidecar
```go
type SidecarConfig struct {
    LogLevel      string `circuit:"type:select,options:debug=Debug;info=Info;error=Error"`
    MetricsPort   int    `circuit:"type:number,help:Prometheus metrics port"`
    HealthCheck   string `circuit:"type:text,help:Health check endpoint"`
}
```

More in [`examples/`](./examples): basic setup, nested structs, hot-reload patterns.

## Design Philosophy

- **Minimal API**: `circuit.From()` - that's the entire public surface
- **No magic**: Reflection for schema extraction, that's it. No code generation, no build steps
- **Production-first**: Battle-tested in industrial IoT deployments and cloud microservices
- **Zero infrastructure**: No database, no message queue, no external dependencies
- **Fail gracefully**: Invalid configs → validation errors, not panics

## How It Works

1. Parses your struct tags to build a schema
2. Loads current values from YAML
3. Renders a clean HTML form (via [gomponents](https://github.com/maragudk/gomponents))
4. Watches the file for external changes
5. Saves edits back to YAML when you submit


## Roadmap

**Authentication & Security**
- Basic auth middleware
- JWT token support
- Role-based field access

**Advanced Features**
- Multi-file configuration support
- Config change audit log
- Webhook notifications on updates
- Advanced validation rules (regex, custom validators)
- Read-only fields
- Grouped/tabbed sections for large configs

**Deployment**
- Docker image with embedded UI
- Kubernetes operator
- Terraform provider integration

PRs welcome. Keep the core simple.

## Contributing

PRs welcome. Keep it simple. Write tests. No "service" or "manager" files.

## License

MIT. Do whatever you want. Just don't blame me if you expose this to the internet without auth.
