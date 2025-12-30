package styles

import "testing"

func TestDepthClass(t *testing.T) {
	tests := []struct {
		name  string
		depth int
		want  string
	}{
		{"negative clamped to 0", -1, CollapsibleDepth0},
		{"depth 0", 0, CollapsibleDepth0},
		{"depth 1", 1, CollapsibleDepth1},
		{"depth 2", 2, CollapsibleDepth2},
		{"depth 3", 3, CollapsibleDepth3},
		{"depth 4", 4, CollapsibleDepth4},
		{"depth 5 clamped to 4", 5, CollapsibleDepth4},
		{"depth 10 clamped to 4", 10, CollapsibleDepth4},
		{"depth 100 clamped to 4", 100, CollapsibleDepth4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DepthClass(tt.depth)
			if got != tt.want {
				t.Errorf("DepthClass(%d) = %q, want %q", tt.depth, got, tt.want)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if Collapsible == "" {
		t.Error("Collapsible constant is empty")
	}
	if Field == "" {
		t.Error("Field constant is empty")
	}
	if Button == "" {
		t.Error("Button constant is empty")
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		want   string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"button"}, "button"},
		{"two classes", []string{"button", "button--primary"}, "button button--primary"},
		{"three classes", []string{"button", "button--primary", "button--large"}, "button button--primary button--large"},
		{"with empty strings", []string{"button", "", "button--primary"}, "button button--primary"},
		{"all empty", []string{"", "", ""}, ""},
		{"empty at start", []string{"", "button", "button--primary"}, "button button--primary"},
		{"empty at end", []string{"button", "button--primary", ""}, "button button--primary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.inputs...)
			if got != tt.want {
				t.Errorf("Merge(%v) = %q, want %q", tt.inputs, got, tt.want)
			}
		})
	}
}
