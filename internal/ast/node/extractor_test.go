package node

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

	if len(schema.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(schema.Nodes))
	}

	if schema.Name != "Config" {
		t.Errorf("expected name Config, got %s", schema.Name)
	}

	host := schema.Nodes[0]
	if host.Name != "Host" {
		t.Errorf("expected Host, got %s", host.Name)
	}
	if host.Kind != KindPrimitive {
		t.Errorf("expected KindPrimitive, got %v", host.Kind)
	}
	if host.ValueType != ValueString {
		t.Errorf("expected ValueString, got %v", host.ValueType)
	}

	port := schema.Nodes[1]
	if port.Name != "Port" {
		t.Errorf("expected Port, got %s", port.Name)
	}
	if port.Kind != KindPrimitive {
		t.Errorf("expected KindPrimitive, got %v", port.Kind)
	}
	if port.ValueType != ValueInt {
		t.Errorf("expected ValueInt, got %v", port.ValueType)
	}

	tls := schema.Nodes[2]
	if tls.Name != "TLS" {
		t.Errorf("expected TLS, got %s", tls.Name)
	}
	if tls.Kind != KindPrimitive {
		t.Errorf("expected KindPrimitive, got %v", tls.Kind)
	}
	if tls.ValueType != ValueBool {
		t.Errorf("expected ValueBool, got %v", tls.ValueType)
	}
}

func TestExtract_EmptyStruct(t *testing.T) {
	type Empty struct{}

	e := Empty{}
	schema, err := Extract(&e)
	if err != nil {
		t.Fatal(err)
	}

	if len(schema.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(schema.Nodes))
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
