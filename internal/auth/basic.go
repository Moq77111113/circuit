package auth

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/moq77111113/circuit/internal/crypto"
)

// Basic authenticates via HTTP Basic Auth.
type Basic struct {
	Username string
	Password string
}

// Authenticate validates Basic Auth credentials.
func (b Basic) Authenticate(r *http.Request) (*Identity, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, fmt.Errorf("basic auth required")
	}

	if subtle.ConstantTimeCompare([]byte(username), []byte(b.Username)) != 1 {
		return nil, fmt.Errorf("invalid credentials")
	}

	if crypto.IsArgon2(b.Password) {
		if !crypto.VerifyArgon2(b.Password, password) {
			return nil, fmt.Errorf("invalid credentials")
		}
	} else {
		if subtle.ConstantTimeCompare([]byte(password), []byte(b.Password)) != 1 {
			return nil, fmt.Errorf("invalid credentials")
		}
	}

	return &Identity{Subject: username}, nil
}
