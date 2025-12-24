package auth

import "net/http"

// None is a no-op authenticator that always succeeds.
type None struct{}

// Authenticate always returns success with subject "anonymous".
func (n None) Authenticate(r *http.Request) (*Identity, error) {
	return &Identity{Subject: "anonymous"}, nil
}
