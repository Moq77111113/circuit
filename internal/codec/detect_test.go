package codec

import "testing"

func TestRegister(t *testing.T) {
	// Clear registry for test isolation
	registry = make(map[Extension]Codec)

	codec := mockCodec{}
	Register(ExtYAML, codec)

	if len(registry) != 1 {
		t.Fatalf("expected 1 codec in registry, got %d", len(registry))
	}

	if _, ok := registry[ExtYAML]; !ok {
		t.Error("expected .yaml extension to be registered")
	}
}

func TestDetect(t *testing.T) {
	// Setup test registry
	registry = make(map[Extension]Codec)
	Register(ExtYAML, mockCodec{})
	Register(ExtYML, mockCodec{})
	Register(ExtTOML, mockCodec{})
	Register(ExtJSON, mockCodec{})

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"yaml extension", "config.yaml", false},
		{"yml extension", "config.yml", false},
		{"toml extension", "config.toml", false},
		{"json extension", "config.json", false},
		{"with path", "/etc/app/config.yaml", false},
		{"unknown extension", "config.xml", true},
		{"no extension", "config", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec, err := Detect(tt.path)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if codec != nil {
					t.Error("expected nil codec on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if codec == nil {
					t.Error("expected non-nil codec")
				}
			}
		})
	}
}

func TestDetectErrorMessage(t *testing.T) {
	registry = make(map[Extension]Codec)

	_, err := Detect("config.xml")
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}

	expected := "unsupported format: .xml"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
