package auth

import "net/http"

// Authenticator validates HTTP requests and returns identity information.
// Implementations must be safe for concurrent use.
type Authenticator interface {
	// Authenticate validates the request and returns the authenticated identity.
	// Returns (*Identity, nil) on success.
	// Returns (nil, error) on authentication failure.
	Authenticate(r *http.Request) (*Identity, error)
}

// Identity represents an authenticated user.
type Identity struct {
	// Subject is the user identifier (username, email, ID, etc).
	Subject string

	// Claims contains optional metadata about the user.
	// Can be nil or empty.
	Claims map[string]string
}
