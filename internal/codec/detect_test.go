package codec

import "testing"

func TestRegister(t *testing.T) {
	ext := Extension(".test")
	codec := mockCodec{}
	Register(ext, codec)

	if _, ok := registry[ext]; !ok {
		t.Error("expected .test extension to be registered")
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"unknown extension", "config.xml", true},
		{"no extension", "config", true},
		{"test extension", "config.test", false},
	}

	Register(Extension(".test"), mockCodec{})

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
	_, err := Detect("config.xml")
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}

	expected := "unsupported format: .xml"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
