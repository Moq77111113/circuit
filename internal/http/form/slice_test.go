package form

import (
	"testing"

	"github.com/moq77111113/circuit/internal/schema"
	"github.com/moq77111113/circuit/internal/tags"
)

// TestAddSliceItem_NestedSlice tests Bug #1 fix: add item to nested slice like "Services.0.Endpoints"
func TestAddSliceItem_NestedSlice(t *testing.T) {
	type Endpoint struct {
		Path   string
		Method string
	}

	type Service struct {
		Name      string
		Endpoints []Endpoint
	}

	type Config struct {
		Services []Service
	}

	cfg := &Config{
		Services: []Service{
			{
				Name:      "UserService",
				Endpoints: []Endpoint{{Path: "/users", Method: "GET"}},
			},
		},
	}

	// Extract schema
	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	// Try to add an item to Services.0.Endpoints (this was failing before Bug #1 fix)
	err = AddSliceItemNode(cfg, s.Nodes, "Services.0.Endpoints")
	if err != nil {
		t.Fatalf("AddSliceItem failed: %v", err)
	}

	// Check that Endpoints slice grew
	if len(cfg.Services[0].Endpoints) != 2 {
		t.Errorf("expected 2 endpoints, got %d", len(cfg.Services[0].Endpoints))
	}
}

// TestAddSliceItem_TopLevel tests add item to top-level slice
func TestAddSliceItem_TopLevel(t *testing.T) {
	type Config struct {
		Tags []string
	}

	cfg := &Config{
		Tags: []string{"go", "web"},
	}

	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	err = AddSliceItemNode(cfg, s.Nodes, "Tags")
	if err != nil {
		t.Fatalf("AddSliceItem failed: %v", err)
	}

	if len(cfg.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(cfg.Tags))
	}
}

// TestRemoveSliceItem_NestedSlice tests Bug #1 fix: remove item from nested slice
func TestRemoveSliceItem_NestedSlice(t *testing.T) {
	type Endpoint struct {
		Path   string
		Method string
	}

	type Service struct {
		Name      string
		Endpoints []Endpoint
	}

	type Config struct {
		Services []Service
	}

	cfg := &Config{
		Services: []Service{
			{
				Name: "UserService",
				Endpoints: []Endpoint{
					{Path: "/users", Method: "GET"},
					{Path: "/users", Method: "POST"},
				},
			},
		},
	}

	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	err = RemoveSliceItemNode(cfg, s.Nodes, "Services.0.Endpoints", 0)
	if err != nil {
		t.Fatalf("RemoveSliceItem failed: %v", err)
	}

	if len(cfg.Services[0].Endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d", len(cfg.Services[0].Endpoints))
	}

	if cfg.Services[0].Endpoints[0].Method != "POST" {
		t.Errorf("expected POST endpoint, got %s", cfg.Services[0].Endpoints[0].Method)
	}
}

// TestAddSliceItem_DeepNested tests deeply nested slices (Services.0.Endpoints.1.AllowedRoles)
func TestAddSliceItem_DeepNested(t *testing.T) {
	type Endpoint struct {
		Path         string
		AllowedRoles []string
	}

	type Service struct {
		Name      string
		Endpoints []Endpoint
	}

	type Config struct {
		Services []Service
	}

	cfg := &Config{
		Services: []Service{
			{
				Name: "UserService",
				Endpoints: []Endpoint{
					{
						Path:         "/users",
						AllowedRoles: []string{"admin"},
					},
				},
			},
		},
	}

	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	err = AddSliceItemNode(cfg, s.Nodes, "Services.0.Endpoints.0.AllowedRoles")
	if err != nil {
		t.Fatalf("AddSliceItem failed: %v", err)
	}

	if len(cfg.Services[0].Endpoints[0].AllowedRoles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(cfg.Services[0].Endpoints[0].AllowedRoles))
	}
}

// TestAddSliceItem_NotFound tests error handling for non-existent field
func TestAddSliceItem_NotFound(t *testing.T) {
	type Config struct {
		Tags []string
	}

	cfg := &Config{Tags: []string{"go"}}

	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	err = AddSliceItemNode(cfg, s.Nodes, "NonExistent")
	if err == nil {
		t.Error("expected error for non-existent field, got nil")
	}
}

// TestAddSliceItem_NotASlice tests error handling for non-slice field
func TestAddSliceItem_NotASlice(t *testing.T) {
	type Config struct {
		Name string
		Tags []string
	}

	cfg := &Config{Name: "test", Tags: []string{"go"}}

	fields, err := tags.Extract(cfg)
	if err != nil {
		t.Fatal(err)
	}

	s := schema.Schema{Nodes: schema.FromTags(fields)}

	err = AddSliceItemNode(cfg, s.Nodes, "Name")
	if err == nil {
		t.Error("expected error for non-slice field, got nil")
	}
}
