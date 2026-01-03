package auth

import "net/http"

// Authenticator validates HTTP requests and returns identity information.
type Authenticator interface {
	Authenticate(r *http.Request) (*Identity, error)
}

// Identity represents an authenticated user.
type Identity struct {
	Subject string
	Claims  map[string]string
}
