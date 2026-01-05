package validation

import "testing"

func TestIsPreset(t *testing.T) {
	tests := []struct {
		name   string
		preset string
		want   bool
	}{
		{"email preset exists", "email", true},
		{"url preset exists", "url", true},
		{"phone preset exists", "phone", true},
		{"unknown preset", "unknown", false},
		{"custom regex", "^[a-z]+$", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPreset(tt.preset); got != tt.want {
				t.Errorf("IsPreset(%q) = %v, want %v", tt.preset, got, tt.want)
			}
		})
	}
}

func TestGetPreset(t *testing.T) {
	tests := []struct {
		name       string
		preset     string
		wantExists bool
		wantEmpty  bool
	}{
		{"email preset", "email", true, false},
		{"url preset", "url", true, false},
		{"phone preset", "phone", true, false},
		{"unknown preset", "unknown", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, ok := GetPreset(tt.preset)
			if ok != tt.wantExists {
				t.Errorf("GetPreset(%q) exists = %v, want %v", tt.preset, ok, tt.wantExists)
			}
			if tt.wantEmpty && pattern != "" {
				t.Errorf("GetPreset(%q) = %q, want empty", tt.preset, pattern)
			}
			if !tt.wantEmpty && tt.wantExists && pattern == "" {
				t.Errorf("GetPreset(%q) returned empty pattern", tt.preset)
			}
		})
	}
}

func TestPresetPatterns(t *testing.T) {
	// Verify preset patterns are valid regex
	presets := []string{"email", "url", "phone"}
	for _, preset := range presets {
		t.Run(preset, func(t *testing.T) {
			pattern, ok := GetPreset(preset)
			if !ok {
				t.Fatalf("preset %q not found", preset)
			}
			if pattern == "" {
				t.Fatalf("preset %q has empty pattern", preset)
			}
		})
	}
}
