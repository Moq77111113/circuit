package main

import (
	"net/http"

	"github.com/moq77111113/circuit"
)

type Config struct {
	Host string `yaml:"host" circuit:"type:text,help:Server hostname"`
	Port int    `yaml:"port" circuit:"type:number,help:Server port,required"`
	TLS  bool   `yaml:"tls" circuit:"type:checkbox,help:Enable TLS"`
}

func main() {
	var cfg Config

	handler, err := circuit.UI(&cfg,
		circuit.WithPath("config.yaml"),
		circuit.WithTitle("Basic Settings"),
	)
	if err != nil {
		panic(err)
	}

	println("Basic example running on :8080")
	http.ListenAndServe(":8080", handler)
}
