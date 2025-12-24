package circuit

import "github.com/moq77111113/circuit/internal/auth"

// Authenticator validates HTTP requests for the Circuit UI.
// Implementations: Forward, Basic.
type Authenticator = auth.Authenticator

// Identity represents an authenticated user.
type Identity = auth.Identity

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
func NewForwardAuth(subjectHeader string, claimHeaders map[string]string) Authenticator {
	return auth.Forward{
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
func NewBasicAuth(username, password string) Authenticator {
	return auth.Basic{
		Username: username,
		Password: password,
	}
}
