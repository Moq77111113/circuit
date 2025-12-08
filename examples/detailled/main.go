package main

import (
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

type Config struct {
	Server         ServerConfig `yaml:"server"`
	UI             UIConfig     `yaml:"ui"`
	Password       string       `yaml:"password" circuit:"type:password,help:Admin password"`
	LogLevel       string       `yaml:"log_level" circuit:"type:select,help:Log verbosity,options:debug=Debug;info=Info;warn=Warning;error=Error"`
	MaxConnections int          `yaml:"max_connections" circuit:"type:range,help:Maximum concurrent connections,min:1,max:100,step:1"`
	StartDate      string       `yaml:"start_date" circuit:"type:date,help:Service start date"`
	BackupTime     string       `yaml:"backup_time" circuit:"type:time,help:Daily backup time"`
}

func main() {
	var cfg Config

	handler, err := circuit.UI(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Detailled Settings"),
	)
	if err != nil {
		panic(err)
	}

	println("Detailled example running on :8080")
	http.ListenAndServe(":8080", handler)
}
