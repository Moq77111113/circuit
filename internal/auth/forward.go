package auth

import (
	"fmt"
	"net/http"
)

// Forward authenticates via headers set by a reverse proxy.
type Forward struct {
	SubjectHeader string
	ClaimHeaders map[string]string
}

// Authenticate validates the request via proxy headers.
func (f Forward) Authenticate(r *http.Request) (*Identity, error) {
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

	return &Identity{
		Subject: subject,
		Claims:  claims,
	}, nil
}
