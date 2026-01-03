package circuit

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/crypto"
)

// Authenticator validates HTTP requests and returns identity information.
// Implementations must be safe for concurrent use.
type Authenticator interface {
	Authenticate(r *http.Request) (*auth.Identity, error)
}

// Identity represents an authenticated user.
type Identity = auth.Identity

// BasicAuth implements HTTP Basic Authentication.
type BasicAuth struct {
	Username string
	Password string // plaintext or argon2id hash
}

// Authenticate validates Basic Auth credentials.
func (b *BasicAuth) Authenticate(r *http.Request) (*auth.Identity, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("basic auth required")
	}

	if subtle.ConstantTimeCompare([]byte(username), []byte(b.Username)) != 1 {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Delegate argon2 verification to internal helpers
	if crypto.IsArgon2(b.Password) {
		if !crypto.VerifyArgon2(b.Password, password) {
			return nil, fmt.Errorf("invalid credentials")
		}
	} else {
		if subtle.ConstantTimeCompare([]byte(password), []byte(b.Password)) != 1 {
			return nil, fmt.Errorf("invalid credentials")
		}
	}

	return &auth.Identity{Subject: username}, nil
}

// ForwardAuth implements authentication via reverse proxy headers.
// Common with OAuth2 Proxy, Traefik ForwardAuth, Cloudflare Access.
type ForwardAuth struct {
	SubjectHeader string
	ClaimHeaders  map[string]string
}

// Authenticate validates the request via proxy headers.
func (f *ForwardAuth) Authenticate(r *http.Request) (*auth.Identity, error) {
	subject := r.Header.Get(f.SubjectHeader)
	if subject == "" {
		return nil, fmt.Errorf("missing auth header: %s", f.SubjectHeader)
	}

	claims := make(map[string]string)
	for name, header := range f.ClaimHeaders {
		if value := r.Header.Get(header); value != "" {
			claims[name] = value
		}
	}

	return &auth.Identity{
		Subject: subject,
		Claims:  claims,
	}, nil
}

// NewForwardAuth creates an authenticator that validates via reverse proxy headers.
// Common with OAuth2 Proxy, Traefik ForwardAuth, Cloudflare Access.
//
// The subjectHeader must contain the authenticated user identifier.
// Optional claimHeaders can map claim names to header names for metadata extraction.
//
// Example:
//
//	auth := circuit.NewForwardAuth("X-Forwarded-User", map[string]string{
//	    "email": "X-Forwarded-Email",
//	    "role": "X-Auth-Role",
//	})
//	ui, _ := circuit.From(&cfg, WithAuth(auth))
func NewForwardAuth(subjectHeader string, claimHeaders map[string]string) *ForwardAuth {
	return &ForwardAuth{
		SubjectHeader: subjectHeader,
		ClaimHeaders:  claimHeaders,
	}
}

// NewBasicAuth creates an authenticator that validates via HTTP Basic Auth.
//
// The password can be plaintext (for dev/testing) or argon2id PHC format (for production).
// Argon2id hashes are auto-detected by the "$argon2id$" prefix.
//
// Example with plaintext (dev only):
//
//	auth := circuit.NewBasicAuth("admin", "password123")
//	ui, _ := circuit.From(&cfg, WithAuth(auth))
//
// Example with argon2id hash (production):
//
//	// Generate hash: see plan for hash generation example
//	auth := circuit.NewBasicAuth("admin", "$argon2id$v=19$m=65536,t=3,p=4$...")
//	ui, _ := circuit.From(&cfg, WithAuth(auth))
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}
