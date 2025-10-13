package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kirillmashkov/GophKeeper.git/internal/util"
	"github.com/kirillmashkov/GophKeeper.git/pkg/hasher"
	"go.uber.org/zap"
)

type mockUserRepo struct {
	PutUserFunc func(ctx context.Context, email string, passwordHash string) (string, error)
	GetUserFunc func(ctx context.Context, email string, passwordHash string) (string, error)
}

func (m *mockUserRepo) PutUser(ctx context.Context, email string, passwordHash string) (string, error) {
	if m.PutUserFunc != nil {
		return m.PutUserFunc(ctx, email, passwordHash)
	}
	return "", nil
}

func (m *mockUserRepo) GetUser(ctx context.Context, email string, passwordHash string) (string, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, email, passwordHash)
	}
	return "", nil
}

func TestAuthService_SignUp_Success(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	email := "user@example.com"
	password := "p@ssw0rd"

	expectedHash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("unexpected hash error: %v", err)
	}

	repo := &mockUserRepo{PutUserFunc: func(ctx context.Context, e string, h string) (string, error) {
		if e != email {
			t.Fatalf("unexpected email: got %s want %s", e, email)
		}
		if h != expectedHash {
			t.Fatalf("unexpected hash: got %s want %s", h, expectedHash)
		}
		return "user-123", nil
	}}

	su := util.NewSecurityUtil(time.Hour, "test-secret-key")

	svc := NewAuthService(repo, logger, su)
	token, err := svc.SignUp(ctx, email, password)
	if err != nil {
		t.Fatalf("SignUp returned error: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
}

func TestAuthService_SignUp_PutUserError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	email := "user@example.com"
	password := "p@ssw0rd"

	repo := &mockUserRepo{PutUserFunc: func(ctx context.Context, e string, h string) (string, error) {
		return "", errors.New("db failure")
	}}

	su := util.NewSecurityUtil(time.Hour, "test-secret-key")
	svc := NewAuthService(repo, logger, su)

	token, err := svc.SignUp(ctx, email, password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}

func TestAuthService_SignIn_Success(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	email := "user@example.com"
	password := "p@ssw0rd"

	expectedHash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("unexpected hash error: %v", err)
	}

	repo := &mockUserRepo{GetUserFunc: func(ctx context.Context, e string, h string) (string, error) {
		if e != email {
			t.Fatalf("unexpected email: got %s want %s", e, email)
		}
		if h != expectedHash {
			t.Fatalf("unexpected hash: got %s want %s", h, expectedHash)
		}
		return "user-123", nil
	}}

	su := util.NewSecurityUtil(time.Hour, "test-secret-key")
	svc := NewAuthService(repo, logger, su)

	ok, token, err := svc.SignIn(ctx, email, password)
	if err != nil {
		t.Fatalf("SignIn returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok==true")
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}
}

func TestAuthService_SignIn_GetUserError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	email := "user@example.com"
	password := "p@ssw0rd"

	repo := &mockUserRepo{GetUserFunc: func(ctx context.Context, e string, h string) (string, error) {
		return "", errors.New("not found")
	}}

	su := util.NewSecurityUtil(time.Hour, "test-secret-key")
	svc := NewAuthService(repo, logger, su)

	ok, token, err := svc.SignIn(ctx, email, password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if ok {
		t.Fatalf("expected ok==false")
	}
	if token != "" {
		t.Fatalf("expected empty token on error, got %q", token)
	}
}
