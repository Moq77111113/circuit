package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/moq77111113/circuit"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Service  ServiceConfig `yaml:"service"`
	HTTP     HTTPConfig    `yaml:"http"`
	Ops      OpsConfig     `yaml:"ops"`
	Limits   LimitsConfig  `yaml:"limits"`
	Backends []Backend     `yaml:"backends"`
	Flags    []Flag        `yaml:"flags"`
}

type ServiceConfig struct {
	Name string `yaml:"name" circuit:"type:text,help:Service name shown in responses,required,minlen:2,maxlen:50,pattern:^[a-z][a-z0-9\-]*$"`
	Env  string `yaml:"env" circuit:"type:select,readonly,options:dev=Development;staging=Staging;prod=Production,help:Environment label"`
}

type HTTPConfig struct {
	PublicMessage string   `yaml:"public_message" circuit:"type:text,help:Message returned by the public endpoint,maxlen:200"`
	AllowedCIDRs  []string `yaml:"allowed_cidrs" circuit:"type:text,help:Optional allowlist for /api (CIDR notation). Empty = allow all"`
}

type OpsConfig struct {
	LogLevel       string  `yaml:"log_level" circuit:"type:select,options:debug=Debug;info=Info;warn=Warning;error=Error,help:Controls server logging verbosity"`
	RequestSample  float64 `yaml:"request_sample" circuit:"type:number,min:0,max:1,step:0.05,help:Percentage of requests to log (0-1)"`
	Maintenance    bool    `yaml:"maintenance" circuit:"type:checkbox,help:If enabled, /api returns 503 but /admin remains accessible"`
	MaintenanceMsg string  `yaml:"maintenance_message" circuit:"type:text,help:Message shown during maintenance,maxlen:150"`
}

type LimitsConfig struct {
	MaxBodyBytes      int `yaml:"max_body_bytes" circuit:"type:number,min:1024,max:10485760,step:1024,help:Maximum request body size for /api (bytes)"`
	ArtificialDelayMS int `yaml:"artificial_delay_ms" circuit:"type:number,min:0,max:2000,step:25,help:Optional latency injection for testing"`
}

type Backend struct {
	Name    string `yaml:"name" circuit:"type:text,required,help:Backend name,minlen:2,maxlen:30,pattern:^[a-z][a-z0-9\-]*$"`
	URL     string `yaml:"url" circuit:"type:url,required,help:Backend base URL,pattern:url"`
	Weight  int    `yaml:"weight" circuit:"type:number,min:1,max:100,help:Relative weight (used by apps that do weighted routing)"`
	Enabled bool   `yaml:"enabled" circuit:"type:checkbox,help:Enable/disable this backend"`
}

type Flag struct {
	Key     string `yaml:"key" circuit:"type:text,required,help:Feature flag key,minlen:2,maxlen:50,pattern:^[a-z][a-z0-9_\-]*$"`
	Enabled bool   `yaml:"enabled" circuit:"type:checkbox,help:Feature flag state"`
}

type liveState struct {
	cfg         atomic.Value // stores Config
	uptimeStart time.Time
	requests    atomic.Uint64
}

func main() {
	log.SetFlags(0)

	configPath, err := resolveConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	if err := ensureConfigFile(configPath); err != nil {
		log.Fatalf("ensure config: %v", err)
	}

	var cfg Config
	st := &liveState{uptimeStart: time.Now()}

	// Circuit UI (this mutates cfg in-process).
	ui, err := circuit.From(&cfg,
		circuit.WithPath(configPath),
		circuit.WithTitle("Circuit Demo"),
		circuit.WithAuth(buildAuth()),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			st.cfg.Store(cfg)
			log.Printf("circuit: config changed (source=%s path=%s)", sourceName(e.Source), e.Path)
		}),
		circuit.WithActions(
			circuit.Action{
				Name:  "toggle-maintenance",
				Label: "Toggle Maintenance",
				Run: func(ctx context.Context) error {
					cfg.Ops.Maintenance = !cfg.Ops.Maintenance
					st.cfg.Store(cfg)
					status := "disabled"
					if cfg.Ops.Maintenance {
						status = "enabled"
					}
					log.Printf("action: maintenance mode %s", status)
					return nil
				},
			}.Describe("Enable or disable maintenance mode").Confirm(),

			circuit.Action{
				Name:  "reset-metrics",
				Label: "Reset Metrics",
				Run: func(ctx context.Context) error {
					st.requests.Store(0)
					st.uptimeStart = time.Now()
					log.Printf("action: metrics reset (requests=0, uptime restarted)")
					return nil
				},
			}.Describe("Reset request counter and uptime"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize live copy after initial load.
	st.cfg.Store(cfg)

	mux := http.NewServeMux()
	mux.Handle("/admin", ui)
	mux.Handle("/admin/", ui)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	mux.Handle("/api", withObservability(st, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snap := st.cfg.Load().(Config)
		if snap.Ops.Maintenance {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status":  "maintenance",
				"message": strings.TrimSpace(nonEmpty(snap.Ops.MaintenanceMsg, "Temporarily unavailable")),
			})
			return
		}

		if !cidrAllowed(r.RemoteAddr, snap.HTTP.AllowedCIDRs) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if snap.Limits.MaxBodyBytes > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, int64(snap.Limits.MaxBodyBytes))
		}

		if snap.Limits.ArtificialDelayMS > 0 {
			time.Sleep(time.Duration(snap.Limits.ArtificialDelayMS) * time.Millisecond)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": map[string]any{
				"name": snap.Service.Name,
				"env":  snap.Service.Env,
			},
			"message":  snap.HTTP.PublicMessage,
			"flags":    snap.Flags,
			"backends": snap.Backends,
			"uptime_s": int(time.Since(st.uptimeStart).Seconds()),
			"requests": st.requests.Load(),
		})
	})))

	addr := ":" + mustPort(os.Getenv("PORT"), "8080")
	log.Printf("listening on %s", addr)
	log.Printf("admin ui: http://localhost%s/admin", addr)
	if isRailway() {
		log.Printf("railway detected: make sure auth is enabled")
	}
	log.Fatal(http.ListenAndServe(addr, mux))
}

func withObservability(st *liveState, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st.requests.Add(1)
		snap := st.cfg.Load().(Config)

		if shouldLogRequest(snap.Ops.RequestSample) {
			lvl := strings.ToLower(strings.TrimSpace(snap.Ops.LogLevel))
			if lvl == "debug" || lvl == "info" {
				log.Printf("http: %s %s from=%s", r.Method, r.URL.Path, r.RemoteAddr)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func shouldLogRequest(sample float64) bool {
	if sample <= 0 {
		return false
	}
	if sample >= 1 {
		return true
	}
	n := time.Now().UnixNano() % 1000
	threshold := int64(sample * 1000)
	return n < threshold
}

func cidrAllowed(remoteAddr string, allowlist []string) bool {
	if len(allowlist) == 0 {
		return true
	}
	ipStr, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		ipStr = remoteAddr
	}
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	for _, cidr := range allowlist {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

func sourceName(s circuit.Source) string {
	switch s {
	case circuit.SourceFormSubmit:
		return "form"
	case circuit.SourceFileChange:
		return "file"
	case circuit.SourceManual:
		return "manual"
	default:
		return "unknown"
	}
}

func buildAuth() circuit.Authenticator {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("CIRCUIT_AUTH_MODE")))
	if mode == "" {
		if isRailway() {
			mode = "basic"
		} else {
			mode = "none"
		}
	}

	switch mode {
	case "none":
		return nil
	case "basic":
		user := os.Getenv("CIRCUIT_BASIC_USER")
		pass := os.Getenv("CIRCUIT_BASIC_PASS")
		if user == "" || pass == "" {
			if isRailway() {
				log.Fatal("CIRCUIT_AUTH_MODE=basic requires CIRCUIT_BASIC_USER and CIRCUIT_BASIC_PASS")
			}
			// Local fallback.
			user, pass = "admin", "admin"
			log.Printf("auth: basic enabled with local default credentials (admin/admin)")
		}
		return circuit.NewBasicAuth(user, pass)
	case "forward":
		subjectHeader := nonEmpty(os.Getenv("CIRCUIT_FORWARD_SUBJECT_HEADER"), "X-Forwarded-User")
		return circuit.NewForwardAuth(subjectHeader, map[string]string{
			"email": nonEmpty(os.Getenv("CIRCUIT_FORWARD_EMAIL_HEADER"), "X-Forwarded-Email"),
		})
	default:
		log.Fatalf("unknown CIRCUIT_AUTH_MODE=%q (use none|basic|forward)", mode)
		return nil
	}
}

func resolveConfigPath() (string, error) {
	if p := strings.TrimSpace(os.Getenv("CONFIG_PATH")); p != "" {
		return p, nil
	}
	// Prefer local config in the example directory when running from repo root.
	if fileExists("examples/demo/config.yaml") {
		return "examples/demo/config.yaml", nil
	}
	if fileExists("config.yaml") {
		return "config.yaml", nil
	}
	// Default to ./config.yaml (and generate it on first run).
	return "config.yaml", nil
}

func ensureConfigFile(path string) error {
	if fileExists(path) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	defaults := Config{
		Service: ServiceConfig{Name: "demo", Env: "dev"},
		HTTP: HTTPConfig{
			PublicMessage: "hello from circuit",
			AllowedCIDRs:  []string{},
		},
		Ops: OpsConfig{
			LogLevel:       "info",
			RequestSample:  0.25,
			Maintenance:    false,
			MaintenanceMsg: "maintenance in progress",
		},
		Limits: LimitsConfig{
			MaxBodyBytes:      1_048_576,
			ArtificialDelayMS: 0,
		},
		Backends: []Backend{
			{Name: "primary", URL: "https://example.com", Weight: 100, Enabled: true},
		},
		Flags: []Flag{
			{Key: "new-ui", Enabled: false},
			{Key: "beta-endpoint", Enabled: true},
		},
	}
	b, err := yaml.Marshal(defaults)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return err
	}
	log.Printf("created default config at %s", path)
	return nil
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func mustPort(v string, fallback string) string {
	p := strings.TrimSpace(v)
	if p == "" {
		return fallback
	}
	_, err := strconv.Atoi(p)
	if err != nil {
		return fallback
	}
	return p
}

func nonEmpty(v string, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

func isRailway() bool {
	for _, k := range []string{
		"RAILWAY_PROJECT_ID",
		"RAILWAY_SERVICE_ID",
		"RAILWAY_ENVIRONMENT_ID",
		"RAILWAY_ENVIRONMENT_NAME",
		"RAILWAY_PUBLIC_DOMAIN",
	} {
		if strings.TrimSpace(os.Getenv(k)) != "" {
			return true
		}
	}
	return false
}
