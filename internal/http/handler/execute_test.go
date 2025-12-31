package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/moq77111113/circuit/internal/actions"
	"github.com/moq77111113/circuit/internal/ast"
	"github.com/moq77111113/circuit/internal/auth"
	"github.com/moq77111113/circuit/internal/sync"
)

func TestExecuteAction_Success(t *testing.T) {
	executed := false
	h := &Handler{
		schema:        ast.Schema{Name: "test"},
		cfg:           &struct{}{},
		path:          "/",
		readOnly:      false,
		store:         &sync.Store{},
		authenticator: auth.None{},
		actions: []actions.Def{
			{
				Name:        "test-action",
				Label:       "Test Action",
				Description: "Test action description",
				Run: func(ctx context.Context) error {
					executed = true
					return nil
				},
				Timeout: 100 * time.Millisecond,
			},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	h.executeAction(w, r, "test-action")

	if !executed {
		t.Fatal("expected action to be executed")
	}

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/" {
		t.Fatalf("expected redirect to /, got %s", location)
	}
}

func TestExecuteAction_Error(t *testing.T) {
	expectedErr := errors.New("action failed")
	h := &Handler{
		schema:        ast.Schema{Name: "test"},
		cfg:           &struct{}{},
		path:          "/",
		readOnly:      false,
		store:         &sync.Store{},
		authenticator: auth.None{},
		actions: []actions.Def{
			{
				Name:        "test-action",
				Label:       "Test Action",
				Description: "Test action description",
				Run: func(ctx context.Context) error {
					return expectedErr
				},
				Timeout: 100 * time.Millisecond,
			},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	h.executeAction(w, r, "test-action")

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "?error=") {
		t.Fatalf("expected redirect with error param, got %s", location)
	}

	expectedErrEncoded := url.QueryEscape(expectedErr.Error())
	if !strings.Contains(location, expectedErrEncoded) {
		t.Fatalf("expected error message %q in redirect, got %s", expectedErrEncoded, location)
	}
}

func TestExecuteAction_NotFound(t *testing.T) {
	h := &Handler{
		schema:        ast.Schema{Name: "test"},
		cfg:           &struct{}{},
		path:          "/",
		readOnly:      false,
		store:         &sync.Store{},
		authenticator: auth.None{},
		actions:       []actions.Def{},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	h.executeAction(w, r, "nonexistent")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestExecuteAction_ReadOnlyMode(t *testing.T) {
	executed := false
	h := &Handler{
		schema:        ast.Schema{Name: "test"},
		cfg:           &struct{}{},
		path:          "/",
		readOnly:      true,
		store:         &sync.Store{},
		authenticator: auth.None{},
		actions: []actions.Def{
			{
				Name:        "test-action",
				Label:       "Test Action",
				Description: "Test action description",
				Run: func(ctx context.Context) error {
					executed = true
					return nil
				},
				Timeout: 100 * time.Millisecond,
			},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	h.executeAction(w, r, "test-action")

	if executed {
		t.Fatal("expected action not to be executed in read-only mode")
	}

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestExecuteAction_Timeout(t *testing.T) {
	h := &Handler{
		schema:        ast.Schema{Name: "test"},
		cfg:           &struct{}{},
		path:          "/",
		readOnly:      false,
		store:         &sync.Store{},
		authenticator: auth.None{},
		actions: []actions.Def{
			{
				Name:        "test-action",
				Label:       "Test Action",
				Description: "Test action description",
				Run: func(ctx context.Context) error {
					select {
					case <-time.After(100 * time.Millisecond):
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				},
				Timeout: 10 * time.Millisecond,
			},
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	h.executeAction(w, r, "test-action")

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
	}

	location := w.Header().Get("Location")
	if !strings.Contains(location, "?error=") {
		t.Fatalf("expected redirect with error param for timeout, got %s", location)
	}
}
