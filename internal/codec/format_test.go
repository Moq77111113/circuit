package codec

import "testing"

// mockCodec implements Codec interface for testing
type mockCodec struct {
	parseErr  error
	encodeErr error
}

func (m mockCodec) Parse(data []byte, dst any) error {
	return m.parseErr
}

func (m mockCodec) Encode(src any) ([]byte, error) {
	if m.encodeErr != nil {
		return nil, m.encodeErr
	}
	return []byte("encoded"), nil
}

func TestCodecInterface(t *testing.T) {
	var c Codec = mockCodec{}

	if err := c.Parse([]byte("test"), nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	data, err := c.Encode(struct{}{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "encoded" {
		t.Errorf("expected 'encoded', got %q", string(data))
	}
}

func TestExtension(t *testing.T) {
	tests := []struct {
		ext      Extension
		expected string
	}{
		{ExtYAML, ".yaml"},
		{ExtYML, ".yml"},
		{ExtTOML, ".toml"},
		{ExtJSON, ".json"},
	}

	for _, tt := range tests {
		if string(tt.ext) != tt.expected {
			t.Errorf("Extension %v: expected %q, got %q",
				tt.ext, tt.expected, string(tt.ext))
		}
	}
}
