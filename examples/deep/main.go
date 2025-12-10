package main

import (
	"log"
	"net/http"

	"github.com/moq77111113/circuit"
)

// Deep nested configuration with slices at multiple levels
type Config struct {
	AppName string `yaml:"app_name" circuit:"type:text,help:Application name"`
	Version string `yaml:"version" circuit:"type:text,help:Application version"`

	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Services []Service      `yaml:"services"`
}

type ServerConfig struct {
	Host string `yaml:"host" circuit:"type:text,help:Server hostname"`
	Port int    `yaml:"port" circuit:"type:number,help:Server port"`
	TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`

	RateLimiting RateLimitConfig `yaml:"rate_limiting"`
	CORS         CORSConfig      `yaml:"cors"`
	Middlewares  []Middleware    `yaml:"middlewares"`
}

type RateLimitConfig struct {
	Enabled        bool     `yaml:"enabled" circuit:"type:checkbox,help:Enable rate limiting"`
	RequestsPerSec int      `yaml:"requests_per_sec" circuit:"type:number,help:Max requests per second"`
	BurstSize      int      `yaml:"burst_size" circuit:"type:number,help:Burst size"`
	Whitelist      []string `yaml:"whitelist" circuit:"type:text,help:Whitelisted IPs"`
}

type CORSConfig struct {
	Enabled        bool     `yaml:"enabled" circuit:"type:checkbox,help:Enable CORS"`
	AllowedOrigins []string `yaml:"allowed_origins" circuit:"type:text,help:Allowed origins"`
	AllowedMethods []string `yaml:"allowed_methods" circuit:"type:text,help:Allowed HTTP methods"`
	AllowedHeaders []string `yaml:"allowed_headers" circuit:"type:text,help:Allowed headers"`
}

type Middleware struct {
	Name    string `yaml:"name" circuit:"type:text,help:Middleware name"`
	Enabled bool   `yaml:"enabled" circuit:"type:checkbox,help:Enable this middleware"`
	Order   int    `yaml:"order" circuit:"type:number,help:Execution order"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver" circuit:"type:select,options:postgres=PostgreSQL;mysql=MySQL;sqlite=SQLite"`
	Host     string `yaml:"host" circuit:"type:text,help:Database host"`
	Port     int    `yaml:"port" circuit:"type:number,help:Database port"`
	Name     string `yaml:"name" circuit:"type:text,help:Database name"`
	Username string `yaml:"username" circuit:"type:text,help:Database username"`
	Password string `yaml:"password" circuit:"type:password,help:Database password"`

	Pool        PoolConfig  `yaml:"pool"`
	Replicas    []Replica   `yaml:"replicas"`
	Maintenance Maintenance `yaml:"maintenance"`
}

type PoolConfig struct {
	MaxOpenConns    int `yaml:"max_open_conns" circuit:"type:number,help:Maximum open connections"`
	MaxIdleConns    int `yaml:"max_idle_conns" circuit:"type:number,help:Maximum idle connections"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime" circuit:"type:number,help:Connection max lifetime (seconds)"`
}

type Replica struct {
	Name     string `yaml:"name" circuit:"type:text,help:Replica name"`
	Host     string `yaml:"host" circuit:"type:text,help:Replica host"`
	Port     int    `yaml:"port" circuit:"type:number,help:Replica port"`
	ReadOnly bool   `yaml:"read_only" circuit:"type:checkbox,help:Read-only replica"`
	Priority int    `yaml:"priority" circuit:"type:number,help:Priority (higher = preferred)"`
}

type Maintenance struct {
	Enabled       bool     `yaml:"enabled" circuit:"type:checkbox,help:Enable maintenance mode"`
	BackupEnabled bool     `yaml:"backup_enabled" circuit:"type:checkbox,help:Enable automatic backups"`
	BackupHour    int      `yaml:"backup_hour" circuit:"type:number,help:Backup hour (0-23)"`
	RetentionDays int      `yaml:"retention_days" circuit:"type:number,help:Backup retention days"`
	AlertEmails   []string `yaml:"alert_emails" circuit:"type:text,help:Alert email addresses"`
}

type Service struct {
	Name       string           `yaml:"name" circuit:"type:text,help:Service name"`
	Enabled    bool             `yaml:"enabled" circuit:"type:checkbox,help:Enable this service"`
	Type       string           `yaml:"type" circuit:"type:select,options:http=HTTP;grpc=gRPC;websocket=WebSocket"`
	Endpoints  []Endpoint       `yaml:"endpoints"`
	Auth       AuthConfig       `yaml:"auth"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

type Endpoint struct {
	Path         string   `yaml:"path" circuit:"type:text,help:Endpoint path"`
	Method       string   `yaml:"method" circuit:"type:select,options:GET=GET;POST=POST;PUT=PUT;DELETE=DELETE;PATCH=PATCH"`
	RateLimit    int      `yaml:"rate_limit" circuit:"type:number,help:Rate limit (requests/min)"`
	Timeout      int      `yaml:"timeout" circuit:"type:number,help:Timeout (seconds)"`
	RequireAuth  bool     `yaml:"require_auth" circuit:"type:checkbox,help:Require authentication"`
	AllowedRoles []string `yaml:"allowed_roles" circuit:"type:text,help:Allowed roles"`
}

type AuthConfig struct {
	Type      string   `yaml:"type" circuit:"type:select,options:jwt=JWT;oauth=OAuth;basic=Basic Auth;apikey=API Key"`
	Secret    string   `yaml:"secret" circuit:"type:password,help:Auth secret"`
	TokenTTL  int      `yaml:"token_ttl" circuit:"type:number,help:Token TTL (minutes)"`
	Providers []string `yaml:"providers" circuit:"type:text,help:Auth providers"`
}

type MonitoringConfig struct {
	Enabled        bool        `yaml:"enabled" circuit:"type:checkbox,help:Enable monitoring"`
	MetricsPort    int         `yaml:"metrics_port" circuit:"type:number,help:Metrics port"`
	HealthPath     string      `yaml:"health_path" circuit:"type:text,help:Health check path"`
	LogLevel       string      `yaml:"log_level" circuit:"type:select,options:debug=Debug;info=Info;warn=Warning;error=Error"`
	TracingEnabled bool        `yaml:"tracing_enabled" circuit:"type:checkbox,help:Enable distributed tracing"`
	AlertRules     []AlertRule `yaml:"alert_rules"`
}

type AlertRule struct {
	Name      string  `yaml:"name" circuit:"type:text,help:Alert rule name"`
	Metric    string  `yaml:"metric" circuit:"type:text,help:Metric to monitor"`
	Threshold float64 `yaml:"threshold" circuit:"type:number,help:Alert threshold"`
	Enabled   bool    `yaml:"enabled" circuit:"type:checkbox,help:Enable this alert"`
}

func main() {
	ui, err := circuit.From(&Config{},
		circuit.WithPath("deep_config.yaml"),
		circuit.WithTitle("Deep Configuration"),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			log.Printf("Configuration changed via %v\n", e.Source)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on http://localhost:8080")
	log.Println("This example demonstrates:")
	log.Println("  • Multiple levels of nested structs")
	log.Println("  • Slices of primitives ([]string, []int)")
	log.Println("  • Slices of structs ([]Service, []Endpoint, []Replica)")
	log.Println("  • Mixed nesting (structs containing slices of structs)")
	log.Fatal(http.ListenAndServe(":8080", ui))
}
