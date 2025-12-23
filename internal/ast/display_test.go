package ast

import "testing"

func TestDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		nodeName string
		want     string
	}{
		{
			name:     "simple name",
			nodeName: "port",
			want:     "port",
		},
		{
			name:     "dot-separated path",
			nodeName: "servers.proxy.port",
			want:     "port",
		},
		{
			name:     "web server host",
			nodeName: "servers.web.host",
			want:     "host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Node{Name: tt.nodeName}
			got := DisplayName(n)
			if got != tt.want {
				t.Errorf("DisplayName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSimplifyPath(t *testing.T) {
	tests := []struct {
		name    string
		pathStr string
		want    string
	}{
		{"simple name", "port", "port"},
		{"single dot", "server.port", "port"},
		{"multiple dots", "servers.proxy.web.port", "port"},
		{"trailing dot", "server.", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SimplifyPath(tt.pathStr)
			if got != tt.want {
				t.Errorf("SimplifyPath(%q) = %q, want %q", tt.pathStr, got, tt.want)
			}
		})
	}
}
