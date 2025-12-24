package auth

import (
	"net/http"
	"testing"
)

func TestForward_SubjectHeaderPresent(t *testing.T) {
	auth := Forward{
		SubjectHeader: "X-Forwarded-User",
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-User", "alice@example.com")

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Subject != "alice@example.com" {
		t.Errorf("Subject = %q, want %q", id.Subject, "alice@example.com")
	}
}

func TestForward_SubjectHeaderMissing(t *testing.T) {
	auth := Forward{
		SubjectHeader: "X-Auth-User",
	}

	req, _ := http.NewRequest("GET", "/", nil)

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error when subject header missing, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}

	expectedMsg := "missing auth header: X-Auth-User"
	if err.Error() != expectedMsg {
		t.Errorf("error message = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestForward_ClaimsExtraction(t *testing.T) {
	auth := Forward{
		SubjectHeader: "X-Forwarded-User",
		ClaimHeaders: map[string]string{
			"email": "X-Forwarded-Email",
			"role":  "X-Auth-Role",
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-User", "bob")
	req.Header.Set("X-Forwarded-Email", "bob@corp.com")
	req.Header.Set("X-Auth-Role", "admin")

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Subject != "bob" {
		t.Errorf("Subject = %q, want %q", id.Subject, "bob")
	}

	if len(id.Claims) != 2 {
		t.Fatalf("Claims length = %d, want 2", len(id.Claims))
	}

	if id.Claims["email"] != "bob@corp.com" {
		t.Errorf("Claims[email] = %q, want %q", id.Claims["email"], "bob@corp.com")
	}

	if id.Claims["role"] != "admin" {
		t.Errorf("Claims[role] = %q, want %q", id.Claims["role"], "admin")
	}
}

func TestForward_ClaimsPartiallyMissing(t *testing.T) {
	auth := Forward{
		SubjectHeader: "X-User",
		ClaimHeaders: map[string]string{
			"team":   "X-Team",
			"region": "X-Region",
		},
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-User", "charlie")
	req.Header.Set("X-Team", "platform")
	// X-Region not set

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(id.Claims) != 1 {
		t.Errorf("Claims length = %d, want 1 (only team present)", len(id.Claims))
	}

	if id.Claims["team"] != "platform" {
		t.Errorf("Claims[team] = %q, want %q", id.Claims["team"], "platform")
	}

	if _, exists := id.Claims["region"]; exists {
		t.Error("Claims[region] should not exist when header missing")
	}
}

func TestForward_NoClaimsConfigured(t *testing.T) {
	auth := Forward{
		SubjectHeader: "X-Auth-User",
		ClaimHeaders:  nil,
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Auth-User", "dana")

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(id.Claims) != 0 {
		t.Errorf("Claims should be empty when no ClaimHeaders configured, got: %v", id.Claims)
	}
}

func TestForward_HeaderCaseSensitivity(t *testing.T) {
	auth := Forward{
		SubjectHeader: "x-forwarded-user", // lowercase
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-User", "eve") // canonical case

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// HTTP headers are case-insensitive, so this should work
	if id.Subject != "eve" {
		t.Errorf("Subject = %q, want %q (headers should be case-insensitive)", id.Subject, "eve")
	}
}
