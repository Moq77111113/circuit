package actions

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecute_Success(t *testing.T) {
	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			return nil
		},
		Timeout: 100 * time.Millisecond,
	}

	err := Execute(context.Background(), action)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestExecute_Error(t *testing.T) {
	expectedErr := errors.New("action failed")
	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			return expectedErr
		},
		Timeout: 100 * time.Millisecond,
	}

	err := Execute(context.Background(), action)
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestExecute_Timeout(t *testing.T) {
	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			select {
			case <-time.After(100 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		Timeout: 10 * time.Millisecond,
	}

	err := Execute(context.Background(), action)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestExecute_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			return ctx.Err()
		},
		Timeout: 100 * time.Millisecond,
	}

	err := Execute(ctx, action)
	if err == nil {
		t.Fatal("expected context canceled error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestExecute_UsesProvidedTimeout(t *testing.T) {
	var receivedCtx context.Context
	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			receivedCtx = ctx
			return nil
		},
		Timeout: 50 * time.Millisecond,
	}

	err := Execute(context.Background(), action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deadline, ok := receivedCtx.Deadline()
	if !ok {
		t.Fatal("expected context to have deadline")
	}

	remainingTime := time.Until(deadline)
	if remainingTime < 0 || remainingTime > 50*time.Millisecond {
		t.Fatalf("expected deadline ~50ms in future, got %v", remainingTime)
	}
}

func TestExecute_DefaultTimeout(t *testing.T) {
	var receivedCtx context.Context
	action := Def{
		Name: "test",
		Run: func(ctx context.Context) error {
			receivedCtx = ctx
			return nil
		},
		Timeout: 0,
	}

	err := Execute(context.Background(), action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deadline, ok := receivedCtx.Deadline()
	if !ok {
		t.Fatal("expected context to have deadline")
	}

	remainingTime := time.Until(deadline)
	if remainingTime < 25*time.Second || remainingTime > 30*time.Second {
		t.Fatalf("expected deadline ~30s in future, got %v", remainingTime)
	}
}
