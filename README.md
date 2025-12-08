# circuit

[![Go Report Card](https://goreportcard.com/badge/github.com/moq77111113/circuit)](https://goreportcard.com/report/github.com/moq77111113/circuit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Because SSH-ing into servers to teach operators which YAML key to edit is nobody's idea of a good time.

**circuit** generates a clean web UI for your Go config structs. Zero JavaScript. Zero build step. Zero drama.

Add a few tags to your struct, call one function, and you get a production-ready admin panel that doesn't look like it's from 2005.

## Why?

Built this after deploying one too many Go apps to embedded systems and industrial endpoints. Operators needed to tweak configs (timeouts, endpoints, thresholds), but asking them to SSH in and edit YAML was a recipe for disaster.

Most solutions involve:
- Installing a whole framework
- Learning a new DSL
- Writing 500 lines of boilerplate
- Teaching people YAML syntax..

**circuit** is one function call:

```go
handler, _ := circuit.UI(&cfg, circuit.WithPath("config.yaml"))
```

That's it. You're done.

## Features

- Auto-generates forms from struct tags
- Hot-reload on file changes
- Clean UI (not your granddad's PHP panel)
- Embedded CSS, no build tools
- Plug-and-play with any Go HTTP server (stdlib, Echo, Gin, Fiber...)

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

handler, err := circuit.UI(&cfg,
    circuit.WithPath("config.yaml"),
    circuit.WithTitle("My App Settings"),
)
if err != nil {
    panic(err)
}

http.ListenAndServe(":8080", handler)
```

### 3. Visit `http://localhost:8080`

Edit. Save. Done. The file updates, your app reloads (if you want).

## Options

```go
circuit.UI(&cfg,
    circuit.WithPath("config.yaml"),        // Required: path to YAML file
    circuit.WithTitle("Admin Panel"),       // Optional: page title
    circuit.OnApply(func() {                // Optional: callback on save
        log.Println("Config updated!")
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

## Examples

Check [`examples/`](./examples) for:
- **basic**: Simple config with 3 fields
- **complex**: Nested structs, all input types
- **reload**: Hot-reload callback example

## Design Philosophy

- **Minimal API**: One function. That's it.
- **No magic**: If it doesn't serve a purpose, it's not here
- **Production-ready**: Used in real systems, not a toy

## How It Works

1. Parses your struct tags to build a schema
2. Loads current values from YAML
3. Renders a clean HTML form (via [gomponents](https://github.com/maragudk/gomponents))
4. Watches the file for external changes
5. Saves edits back to YAML when you submit


## Future Features

- Built-in authentication (basic auth)
- Config change history/audit log
- Multi-file support
- Field validation rules
- Read-only mode for certain fields
- Mobile-friendly responsive layout

PRs welcome for any of these.

## Contributing

PRs welcome. Keep it simple. Write tests. No "service" or "manager" files.

## License

MIT. Do whatever you want. Just don't blame me if you expose this to the internet without auth.
