package path

import (
	"testing"
)

func TestPath_Simple(t *testing.T) {
	p := NewPath("Port")
	got := p.String()
	want := "Port"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPath_Nested(t *testing.T) {
	p := NewPath("Database").Child("Host")
	got := p.String()
	want := "Database.Host"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPath_Indexed(t *testing.T) {
	p := NewPath("Items").Index(0)
	got := p.String()
	want := "Items.0"

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPath_DeepNested(t *testing.T) {
	tests := []struct {
		name string
		path Path
		want string
	}{
		{
			name: "Services.0.Endpoints.2.Name",
			path: NewPath("Services").Index(0).Child("Endpoints").Index(2).Child("Name"),
			want: "Services.0.Endpoints.2.Name",
		},
		{
			name: "Database.Replicas.1.Host",
			path: NewPath("Database").Child("Replicas").Index(1).Child("Host"),
			want: "Database.Replicas.1.Host",
		},
		{
			name: "Tags.0",
			path: NewPath("Tags").Index(0),
			want: "Tags.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Port", "Port"},
		{"Database.Host", "Database.Host"},
		{"Items.0", "Items.0"},
		{"Services.0.Endpoints.2.Name", "Services.0.Endpoints.2.Name"},
		{"Database.Replicas.1.Host", "Database.Replicas.1.Host"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := ParsePath(tt.input)
			got := p.String()
			if got != tt.want {
				t.Errorf("ParsePath(%q).String() = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPath_HasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		path   Path
		prefix Path
		want   bool
	}{
		{
			name:   "exact match",
			path:   NewPath("Services"),
			prefix: NewPath("Services"),
			want:   true,
		},
		{
			name:   "simple prefix",
			path:   NewPath("Database").Child("Host"),
			prefix: NewPath("Database"),
			want:   true,
		},
		{
			name:   "indexed prefix",
			path:   NewPath("Services").Index(0).Child("Name"),
			prefix: NewPath("Services").Index(0),
			want:   true,
		},
		{
			name:   "deep prefix",
			path:   NewPath("Services").Index(1).Child("Endpoints").Index(2).Child("Path"),
			prefix: NewPath("Services").Index(1).Child("Endpoints"),
			want:   true,
		},
		{
			name:   "no match",
			path:   NewPath("Database").Child("Host"),
			prefix: NewPath("Server"),
			want:   false,
		},
		{
			name:   "longer prefix",
			path:   NewPath("Services"),
			prefix: NewPath("Services").Child("Name"),
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.HasPrefix(tt.prefix)
			if got != tt.want {
				t.Errorf("HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPath_IndexAfter(t *testing.T) {
	tests := []struct {
		name   string
		path   Path
		prefix Path
		want   int
	}{
		{
			name:   "Services.0.Name after Services",
			path:   NewPath("Services").Index(0).Child("Name"),
			prefix: NewPath("Services"),
			want:   0,
		},
		{
			name:   "Services.1.Endpoints.2.Path after Services.1.Endpoints",
			path:   NewPath("Services").Index(1).Child("Endpoints").Index(2).Child("Path"),
			prefix: NewPath("Services").Index(1).Child("Endpoints"),
			want:   2,
		},
		{
			name:   "no index after prefix",
			path:   NewPath("Services").Child("Name"),
			prefix: NewPath("Services"),
			want:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.IndexAfter(tt.prefix)
			if got != tt.want {
				t.Errorf("IndexAfter() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPath_Empty(t *testing.T) {
	p := Path{}
	got := p.String()
	want := ""

	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPath_FieldPath(t *testing.T) {
	tests := []struct {
		name string
		path Path
		want string
	}{
		{
			name: "simple path",
			path: NewPath("Port"),
			want: "Port",
		},
		{
			name: "nested path",
			path: NewPath("Database").Child("Host"),
			want: "Database.Host",
		},
		{
			name: "path with index strips index",
			path: NewPath("Services").Index(0),
			want: "Services",
		},
		{
			name: "deep path with indices strips all indices",
			path: NewPath("Services").Index(0).Child("Endpoints").Index(2).Child("Name"),
			want: "Services.Endpoints.Name",
		},
		{
			name: "multiple indices",
			path: NewPath("Database").Child("Replicas").Index(1).Child("Servers").Index(3).Child("Host"),
			want: "Database.Replicas.Servers.Host",
		},
		{
			name: "empty path",
			path: Path{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.FieldPath()
			if got != tt.want {
				t.Errorf("FieldPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
