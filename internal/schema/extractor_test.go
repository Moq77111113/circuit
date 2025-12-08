package schema

import "testing"

func TestExtract_BasicTypes(t *testing.T) {
	type Config struct {
		Host string `circuit:"text,help:Server hostname"`
		Port int    `circuit:"number,help:Server port"`
		TLS  bool   `circuit:"checkbox,help:Enable TLS"`
	}

	cfg := Config{}
	schema, err := Extract(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(schema.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(schema.Fields))
	}

	if schema.Name != "Config" {
		t.Errorf("expected name Config, got %s", schema.Name)
	}

	// Verify Host field
	host := schema.Fields[0]
	if host.Name != "Host" {
		t.Errorf("expected Host, got %s", host.Name)
	}
	if host.Type != "string" {
		t.Errorf("expected string, got %s", host.Type)
	}

	// Verify Port field
	port := schema.Fields[1]
	if port.Name != "Port" {
		t.Errorf("expected Port, got %s", port.Name)
	}
	if port.Type != "int" {
		t.Errorf("expected int, got %s", port.Type)
	}

	// Verify TLS field
	tls := schema.Fields[2]
	if tls.Name != "TLS" {
		t.Errorf("expected TLS, got %s", tls.Name)
	}
	if tls.Type != "bool" {
		t.Errorf("expected bool, got %s", tls.Type)
	}
}

func TestExtract_EmptyStruct(t *testing.T) {
	type Empty struct{}

	e := Empty{}
	schema, err := Extract(&e)
	if err != nil {
		t.Fatal(err)
	}

	if len(schema.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(schema.Fields))
	}

	if schema.Name != "Empty" {
		t.Errorf("expected name Empty, got %s", schema.Name)
	}
}

func TestExtract_NonPointer(t *testing.T) {
	type Config struct {
		Host string `circuit:"text"`
	}

	cfg := Config{}
	_, err := Extract(cfg)
	if err == nil {
		t.Fatal("expected error for non-pointer")
	}
}
