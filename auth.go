package circuit

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/crypto"
)

// Authenticator validates HTTP requests and returns identity information.
//
// Implementations must be safe for concurrent use.
//
// The Authenticate method is called on EVERY request to the Circuit handler.
// Return an error to reject the request with 401 Unauthorized.
//
// Built-in implementations:
//   - BasicAuth - HTTP Basic Authentication
//   - ForwardAuth - Reverse proxy header authentication
type Authenticator interface {
	Authenticate(r *http.Request) (*auth.Identity, error)
}

// Identity represents an authenticated user.
type Identity = auth.Identity

// BasicAuth implements HTTP Basic Authentication with support for plaintext
// and argon2id hashed passwords.
//
// IMPORTANT: For production use, always use argon2id password hashes, not plaintext.
// Plaintext passwords should ONLY be used in development/testing environments.
//
// Argon2id hashes are automatically detected by the "$argon2id$v=19$" prefix.
// Circuit supports the PHC string format output by golang.org/x/crypto/argon2.
//
// Operational security recommendations:
//   - Store credentials in a separate file (e.g., /etc/myapp/auth.conf)
//   - Set file permissions to 0640 or stricter (readable only by app user)
//   - Never store credentials in your config struct (they would be editable via Circuit UI)
//   - Use environment variables or secret management systems in production
type BasicAuth struct {
	Username string
	Password string // plaintext or argon2id PHC hash
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
// The password can be:
//   - Plaintext (for dev/testing ONLY) - never use in production
//   - Argon2id PHC hash (for production) - auto-detected by "$argon2id$v=19$" prefix
//
// Circuit supports argon2id hashes in PHC string format as output by the
// golang.org/x/crypto/argon2 package. The hash format is:
//
//	$argon2id$v=19$m=MEMORY,t=TIME,p=PARALLELISM$SALT$HASH
//
// Example with plaintext (DEVELOPMENT ONLY):
//
//	auth := circuit.NewBasicAuth("admin", "dev-password")
//	ui, _ := circuit.From(&cfg, circuit.WithAuth(auth))
//
// Example with argon2id hash (PRODUCTION):
//
//	// Generate hash with: golang.org/x/crypto/argon2
//	// Example hash:
//	hash := "$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$..."
//	auth := circuit.NewBasicAuth("admin", hash)
//	ui, _ := circuit.From(&cfg, circuit.WithAuth(auth))
//
// Operational security:
//   - Never store credentials in your app's config struct
//   - Store credentials in /etc/myapp/auth.conf with 0640 permissions
//   - Use a secret management system (Vault, AWS Secrets Manager, etc.)
//   - Rotate passwords regularly
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}
