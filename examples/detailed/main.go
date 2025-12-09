package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moq77111113/circuit"
)

type ServerConfig struct {
	Host string `yaml:"host" circuit:"type:text,help:Server hostname,required"`
	Port int    `yaml:"port" circuit:"type:number,help:Server port,required,min:1,max:65535"`
	TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`
}

type UIConfig struct {
	Theme string `yaml:"theme" circuit:"type:radio,help:UI Theme,options:light=Light Mode;dark=Dark Mode"`
}

type Rule struct {
	Name   string `yaml:"name" circuit:"type:text,help:Rule name"`
	Active bool   `yaml:"active" circuit:"type:checkbox,help:Rule enabled"`
}

type Config struct {
	Server         ServerConfig `yaml:"server"`
	UI             UIConfig     `yaml:"ui"`
	Password       string       `yaml:"password" circuit:"type:password,help:Admin password"`
	LogLevel       string       `yaml:"log_level" circuit:"type:select,help:Log verbosity,options:debug=Debug;info=Info;warn=Warning;error=Error"`
	MaxConnections int          `yaml:"max_connections" circuit:"type:range,help:Maximum concurrent connections,min:1,max:100,step:1"`
	StartDate      string       `yaml:"start_date" circuit:"type:date,help:Service start date"`
	BackupTime     string       `yaml:"backup_time" circuit:"type:time,help:Daily backup time"`
	AllowedIPs     []string     `yaml:"allowed_ips" circuit:"type:text,help:Allowed IP addresses"`
	BlockedPorts   []int        `yaml:"blocked_ports" circuit:"type:number,help:Blocked ports"`
	Rules          []Rule       `yaml:"rules"`
	AdminUser      *string      `yaml:"admin_user" circuit:"type:text,help:Optional admin username"`
}

func main() {
	var cfg Config

	handler, err := circuit.From(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Detailed Configuration"),
		circuit.WithOnChange(func(e circuit.ChangeEvent) {
			source := "unknown"
			switch e.Source {
			case circuit.SourceFormSubmit:
				source = "form submit"
			case circuit.SourceFileChange:
				source = "file change"
			case circuit.SourceManual:
				source = "manual reload"
			}
			log.Printf("Config changed via %s: %d allowed IPs, %d blocked ports\n",
				source, len(cfg.AllowedIPs), len(cfg.BlockedPorts))
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Detailed example running on http://localhost:8080")
	fmt.Println("Features:")
	fmt.Println("  - Nested structs (Server, UI)")
	fmt.Println("  - String slices (AllowedIPs)")
	fmt.Println("  - Int slices (BlockedPorts)")
	fmt.Println("  - Struct slices (Rules)")
	fmt.Println("  - Pointer fields (AdminUser)")
	fmt.Println("  - Change events with source tracking")
	fmt.Println()

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
