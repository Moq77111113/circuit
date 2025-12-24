package auth

import (
	"net/http"
	"testing"
)

// TestIdentity_Construction tests that Identity can be constructed with various configurations.
func TestIdentity_Construction(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		claims  map[string]string
	}{
		{
			name:    "with subject only",
			subject: "user@example.com",
			claims:  nil,
		},
		{
			name:    "with subject and claims",
			subject: "admin",
			claims:  map[string]string{"role": "admin", "team": "platform"},
		},
		{
			name:    "with empty claims map",
			subject: "guest",
			claims:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &Identity{
				Subject: tt.subject,
				Claims:  tt.claims,
			}

			if id.Subject != tt.subject {
				t.Errorf("Subject = %q, want %q", id.Subject, tt.subject)
			}

			if tt.claims == nil && id.Claims != nil {
				t.Errorf("Claims should be nil, got %v", id.Claims)
			}

			if tt.claims != nil && len(id.Claims) != len(tt.claims) {
				t.Errorf("Claims length = %d, want %d", len(id.Claims), len(tt.claims))
			}
		})
	}
}

// TestAuthenticator_InterfaceContract documents the expected behavior of Authenticator implementations.
func TestAuthenticator_InterfaceContract(t *testing.T) {
	// This test documents the Authenticator interface contract:
	//
	// 1. Authenticate(r *http.Request) (*Identity, error)
	//    - Returns (*Identity, nil) on successful authentication
	//    - Returns (nil, error) on failed authentication
	//    - Must not panic
	//    - Must be safe for concurrent use
	//
	// 2. Identity contains:
	//    - Subject: required user identifier
	//    - Claims: optional metadata (can be nil or empty)
	//
	// 3. Implementations must validate requests independently
	//    - No shared state between calls
	//    - Each request authenticated in isolation

	var _ Authenticator = (*mockAuth)(nil) // compile-time interface check

	auth := &mockAuth{valid: true}
	req, _ := http.NewRequest("GET", "/", nil)

	id, err := auth.Authenticate(req)
	if err != nil {
		t.Errorf("valid auth should not error, got: %v", err)
	}
	if id == nil {
		t.Fatal("valid auth should return Identity")
	}
	if id.Subject == "" {
		t.Error("Identity.Subject should not be empty")
	}
}

// mockAuth is a minimal Authenticator implementation for testing the interface contract.
type mockAuth struct {
	valid bool
}

func (m *mockAuth) Authenticate(r *http.Request) (*Identity, error) {
	if !m.valid {
		return nil, http.ErrNoCookie // arbitrary error
	}
	return &Identity{Subject: "test-user"}, nil
}
