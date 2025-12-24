// Package circuit provides an embeddable HTTP UI for viewing and reloading
// application configuration stored as YAML.
//
// The primary entry point is From, which accepts a pointer to a configuration
// value (a struct) and returns an `http.Handler` that serves a small UI for
// inspecting and updating the values. The package extracts schema information
// from struct tags to build form fields, and it can watch the YAML file for
// changes and call an optional callback when a new configuration is applied.
//
// Example (usage):
//
//	type Config struct {
//	    Addr string `circuit:",required"`
//	    Port int    `circuit:",min=1,max=65535"`
//	}
//
//	cfg := &Config{}
//	h, err := circuit.From(cfg,
//	    circuit.WithPath("/etc/myapp/config.yaml"),
//	    circuit.WithTitle("My App Config"),
//	    circuit.WithAuth(circuit.NewBasicAuth("admin", "password")),
//	    circuit.OnApply(func(){ log.Println("config reloaded") }),
//	)
//	if err != nil {
//	    // handle error
//	}
//	http.Handle("/config", h)
//
// Notes:
//   - `cfg` must be a pointer to a struct; the package uses reflection to
//     extract schema metadata from struct tags.
//   - `WithPath` is required so the loader can read the initial YAML file.
//   - The returned handler delegates to an internal loader which watches the
//     file for changes; provide `OnApply` to be notified after successful reloads.
package circuit
