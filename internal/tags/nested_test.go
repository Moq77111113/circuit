package tags

import (
	"testing"
)

func TestExtract_NestedStruct(t *testing.T) {
	type ServerConfig struct {
		Port int `circuit:"type:number"`
	}

	type Config struct {
		Host   string `circuit:"type:text"`
		Server ServerConfig
	}

	cfg := Config{}
	fields, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	host := fields[0]
	if host.Name != "Host" {
		t.Errorf("expected Host, got %s", host.Name)
	}

	server := fields[1]
	if server.Name != "Server" {
		t.Errorf("expected Server, got %s", server.Name)
	}
	if len(server.Fields) != 1 {
		t.Fatalf("expected 1 sub-field, got %d", len(server.Fields))
	}

	port := server.Fields[0]
	if port.Name != "Port" {
		t.Errorf("expected Port, got %s", port.Name)
	}
	if port.InputType != "number" {
		t.Errorf("expected number, got %s", port.InputType)
	}
}
