package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestBasic_PlaintextSuccess(t *testing.T) {
	auth := Basic{
		Username: "admin",
		Password: "secret123",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("admin", "secret123")

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Subject != "admin" {
		t.Errorf("Subject = %q, want %q", id.Subject, "admin")
	}
}

func TestBasic_Argon2Success(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatal(err)
	}
	hash := argon2.IDKey([]byte("mypassword"), salt, 3, 64*1024, 4, 32)

	encoded := "$argon2id$v=19$m=65536,t=3,p=4$" +
		base64.RawStdEncoding.EncodeToString(salt) + "$" +
		base64.RawStdEncoding.EncodeToString(hash)

	auth := Basic{
		Username: "alice",
		Password: encoded,
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("alice", "mypassword")

	id, err := auth.Authenticate(req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id.Subject != "alice" {
		t.Errorf("Subject = %q, want %q", id.Subject, "alice")
	}
}

func TestBasic_MissingAuthHeader(t *testing.T) {
	auth := Basic{
		Username: "user",
		Password: "pass",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error when Authorization header missing, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}
}

func TestBasic_WrongUsername(t *testing.T) {
	auth := Basic{
		Username: "admin",
		Password: "correct",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("attacker", "correct")

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error for wrong username, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}

	expectedMsg := "invalid credentials"
	if err.Error() != expectedMsg {
		t.Errorf("error message = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestBasic_WrongPasswordPlaintext(t *testing.T) {
	auth := Basic{
		Username: "bob",
		Password: "correctpassword",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("bob", "wrongpassword")

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}
}

func TestBasic_WrongPasswordArgon2(t *testing.T) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		t.Fatal(err)
	}
	hash := argon2.IDKey([]byte("correctpass"), salt, 3, 64*1024, 4, 32)

	encoded := "$argon2id$v=19$m=65536,t=3,p=4$" +
		base64.RawStdEncoding.EncodeToString(salt) + "$" +
		base64.RawStdEncoding.EncodeToString(hash)

	auth := Basic{
		Username: "charlie",
		Password: encoded,
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("charlie", "wrongpass")

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error for wrong argon2 password, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}
}

func TestBasic_MalformedBasicAuth(t *testing.T) {
	auth := Basic{
		Username: "user",
		Password: "pass",
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer invalid-token")

	id, err := auth.Authenticate(req)

	if err == nil {
		t.Fatal("expected error for malformed auth header, got nil")
	}

	if id != nil {
		t.Errorf("Identity should be nil on error, got: %v", id)
	}
}

func TestBasic_ConstantTimeComparison(t *testing.T) {
	auth := Basic{
		Username: "test",
		Password: "password",
	}

	tests := []struct {
		username string
		password string
		wantErr  bool
	}{
		{"test", "password", false}, // exact match
		{"test", "passwor", true},   // prefix wrong
		{"test", "password2", true}, // suffix wrong
		{"tes", "password", true},   // username prefix
		{"test2", "password", true}, // username suffix
	}

	for _, tt := range tests {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(tt.username, tt.password)

		_, err = auth.Authenticate(req)

		if (err != nil) != tt.wantErr {
			t.Errorf("username=%q password=%q: error=%v, wantErr=%v",
				tt.username, tt.password, err, tt.wantErr)
		}
	}
}
