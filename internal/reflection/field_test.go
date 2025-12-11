package reflection

import "testing"

type TestConfig struct {
	Name string
	Port int
}

func TestFieldByName(t *testing.T) {
	cfg := TestConfig{Name: "test", Port: 8080}

	t.Run("extract string field", func(t *testing.T) {
		val := FieldByName(cfg, "Name")
		if val == nil {
			t.Fatal("expected non-nil value")
		}
		if val.(string) != "test" {
			t.Fatalf("expected 'test', got %v", val)
		}
	})

	t.Run("extract int field", func(t *testing.T) {
		val := FieldByName(cfg, "Port")
		if val == nil {
			t.Fatal("expected non-nil value")
		}
		if val.(int) != 8080 {
			t.Fatalf("expected 8080, got %v", val)
		}
	})

	t.Run("pointer to struct", func(t *testing.T) {
		val := FieldByName(&cfg, "Name")
		if val == nil {
			t.Fatal("expected non-nil value")
		}
		if val.(string) != "test" {
			t.Fatalf("expected 'test', got %v", val)
		}
	})

	t.Run("non-existent field", func(t *testing.T) {
		val := FieldByName(cfg, "NonExistent")
		if val != nil {
			t.Fatalf("expected nil, got %v", val)
		}
	})

	t.Run("nil value", func(t *testing.T) {
		val := FieldByName(nil, "Name")
		if val != nil {
			t.Fatalf("expected nil, got %v", val)
		}
	})

	t.Run("non-struct value", func(t *testing.T) {
		val := FieldByName("not a struct", "Name")
		if val != nil {
			t.Fatalf("expected nil, got %v", val)
		}
	})
}
