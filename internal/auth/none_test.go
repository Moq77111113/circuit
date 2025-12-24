package auth

import (
	"net/http"
	"testing"
)

func TestNone_AlwaysSucceeds(t *testing.T) {
	auth := None{}
	req, _ := http.NewRequest("GET", "/", nil)

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Errorf("None.Authenticate() should never return error, got: %v", err)
	}

	if id == nil {
		t.Fatal("None.Authenticate() should return Identity")
	}
}

func TestNone_ReturnsAnonymous(t *testing.T) {
	auth := None{}
	req, _ := http.NewRequest("GET", "/", nil)

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Subject != "anonymous" {
		t.Errorf("Subject = %q, want %q", id.Subject, "anonymous")
	}
}

func TestNone_NoClaims(t *testing.T) {
	auth := None{}
	req, _ := http.NewRequest("POST", "/admin", nil)

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Claims != nil {
		t.Errorf("Claims should be nil, got: %v", id.Claims)
	}
}
