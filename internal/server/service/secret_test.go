package service

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"go.uber.org/zap"
)

type mockSecretRepo struct {
	CreateSecretFunc func(ctx context.Context, userID string, content []byte, name string, kind int, metadata []byte) (string, error)
	GetSecretFunc    func(ctx context.Context, userID string, name string) (model.SecretDB, error)
	UpdateSecretFunc func(ctx context.Context, secretID string, content []byte, kind int, metadata []byte) error

	// call counters/args for assertions
	CreateCalled int
	UpdateCalled int
}

func (m *mockSecretRepo) CreateSecret(ctx context.Context, userID string, content []byte, name string, kind int, metadata []byte) (string, error) {
	m.CreateCalled++
	if m.CreateSecretFunc != nil {
		return m.CreateSecretFunc(ctx, userID, content, name, kind, metadata)
	}
	return "", nil
}

func (m *mockSecretRepo) GetSecret(ctx context.Context, userID string, name string) (model.SecretDB, error) {
	if m.GetSecretFunc != nil {
		return m.GetSecretFunc(ctx, userID, name)
	}
	return model.SecretDB{}, nil
}

func (m *mockSecretRepo) UpdateSecret(ctx context.Context, secretID string, content []byte, kind int, metadata []byte) error {
	m.UpdateCalled++
	if m.UpdateSecretFunc != nil {
		return m.UpdateSecretFunc(ctx, secretID, content, kind, metadata)
	}
	return nil
}

func TestSecretService_GetSecret_Success(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	userID := "user-1"
	name := "cred1"
	data := []byte("secret-bytes")
	meta := []byte("meta-bytes")

	repo := &mockSecretRepo{GetSecretFunc: func(ctx context.Context, u string, n string) (model.SecretDB, error) {
		if u != userID || n != name {
			t.Fatalf("unexpected args: got user=%s name=%s", u, n)
		}
		return model.SecretDB{ID: "sec-1", Name: name, Data: data, Metadata: meta, Kind: 2}, nil
	}}

	svc := NewSecretService(repo, logger)
	resp, err := svc.GetSecret(ctx, userID, name)
	if err != nil {
		t.Fatalf("GetSecret returned error: %v", err)
	}

	if resp.Name != name {
		t.Fatalf("unexpected name: got %s want %s", resp.Name, name)
	}
	if resp.Kind != 2 {
		t.Fatalf("unexpected kind: got %d want %d", resp.Kind, 2)
	}

	expectedData := base64.StdEncoding.EncodeToString(data)
	expectedMeta := base64.StdEncoding.EncodeToString(meta)
	if resp.Data != expectedData {
		t.Fatalf("unexpected data: got %s want %s", resp.Data, expectedData)
	}
	if resp.Metadata != expectedMeta {
		t.Fatalf("unexpected metadata: got %s want %s", resp.Metadata, expectedMeta)
	}
}

func TestSecretService_GetSecret_RepoError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()

	repo := &mockSecretRepo{GetSecretFunc: func(ctx context.Context, userID string, name string) (model.SecretDB, error) {
		return model.SecretDB{}, errors.New("db error")
	}}

	svc := NewSecretService(repo, logger)
	_, err := svc.GetSecret(ctx, "user", "name")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSecretService_UpdateSecret_Success(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	userID := "user-1"
	name := "cred1"
	secretID := "sec-1"
	contentRaw := []byte("content")
	metadataRaw := []byte("metadata")
	contentB64 := base64.StdEncoding.EncodeToString(contentRaw)
	metadataB64 := base64.StdEncoding.EncodeToString(metadataRaw)
	kind := 3

	repo := &mockSecretRepo{
		GetSecretFunc: func(ctx context.Context, u string, n string) (model.SecretDB, error) {
			if u != userID || n != name {
				t.Fatalf("unexpected args: got user=%s name=%s", u, n)
			}
			return model.SecretDB{ID: secretID, Name: name, Data: []byte("old"), Metadata: []byte("oldm"), Kind: 1}, nil
		},
		UpdateSecretFunc: func(ctx context.Context, id string, content []byte, k int, metadata []byte) error {
			if id != secretID {
				t.Fatalf("unexpected id: got %s want %s", id, secretID)
			}
			if k != kind {
				t.Fatalf("unexpected kind: got %d want %d", k, kind)
			}
			if string(content) != string(contentRaw) {
				t.Fatalf("unexpected content: got %q want %q", string(content), string(contentRaw))
			}
			if string(metadata) != string(metadataRaw) {
				t.Fatalf("unexpected metadata: got %q want %q", string(metadata), string(metadataRaw))
			}
			return nil
		},
	}

	svc := NewSecretService(repo, logger)
	req := model.UpdateSecretRequest{Name: name, Data: contentB64, Metadata: metadataB64, Kind: kind}
	if err := svc.UpdateSecret(ctx, userID, req); err != nil {
		t.Fatalf("UpdateSecret returned error: %v", err)
	}

	if repo.UpdateCalled != 1 {
		t.Fatalf("expected UpdateSecret to be called once, got %d", repo.UpdateCalled)
	}
}

func TestSecretService_UpdateSecret_GetSecretError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{GetSecretFunc: func(ctx context.Context, userID, name string) (model.SecretDB, error) {
		return model.SecretDB{}, errors.New("not found")
	}}

	svc := NewSecretService(repo, logger)
	req := model.UpdateSecretRequest{Name: "n", Data: base64.StdEncoding.EncodeToString([]byte("d")), Metadata: base64.StdEncoding.EncodeToString([]byte("m")), Kind: 1}
	if err := svc.UpdateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if repo.UpdateCalled != 0 {
		t.Fatalf("UpdateSecret should not be called on repo when GetSecret fails")
	}
}

func TestSecretService_UpdateSecret_DecodeDataError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{GetSecretFunc: func(ctx context.Context, userID, name string) (model.SecretDB, error) {
		return model.SecretDB{ID: "id"}, nil
	}}

	svc := NewSecretService(repo, logger)
	req := model.UpdateSecretRequest{Name: "n", Data: "!!!not-base64!!!", Metadata: base64.StdEncoding.EncodeToString([]byte("m")), Kind: 1}
	if err := svc.UpdateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.UpdateCalled != 0 {
		t.Fatalf("repo UpdateSecret should not be called when data decode fails")
	}
}

func TestSecretService_UpdateSecret_DecodeMetadataError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{GetSecretFunc: func(ctx context.Context, userID, name string) (model.SecretDB, error) {
		return model.SecretDB{ID: "id"}, nil
	}}

	svc := NewSecretService(repo, logger)
	req := model.UpdateSecretRequest{Name: "n", Data: base64.StdEncoding.EncodeToString([]byte("d")), Metadata: "!!!not-base64!!!", Kind: 1}
	if err := svc.UpdateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.UpdateCalled != 0 {
		t.Fatalf("repo UpdateSecret should not be called when metadata decode fails")
	}
}

func TestSecretService_UpdateSecret_RepoUpdateError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{
		GetSecretFunc: func(ctx context.Context, userID, name string) (model.SecretDB, error) {
			return model.SecretDB{ID: "id"}, nil
		},
		UpdateSecretFunc: func(ctx context.Context, secretID string, content []byte, kind int, metadata []byte) error {
			return errors.New("update failed")
		},
	}

	svc := NewSecretService(repo, logger)
	req := model.UpdateSecretRequest{Name: "n", Data: base64.StdEncoding.EncodeToString([]byte("d")), Metadata: base64.StdEncoding.EncodeToString([]byte("m")), Kind: 1}
	if err := svc.UpdateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.UpdateCalled != 1 {
		t.Fatalf("expected repo UpdateSecret to be called once")
	}
}

func TestSecretService_CreateSecret_Success(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	userID := "user-1"
	name := "cred1"
	contentRaw := []byte("content")
	metadataRaw := []byte("metadata")
	contentB64 := base64.StdEncoding.EncodeToString(contentRaw)
	metadataB64 := base64.StdEncoding.EncodeToString(metadataRaw)
	kind := 4

	repo := &mockSecretRepo{CreateSecretFunc: func(ctx context.Context, u string, content []byte, n string, k int, metadata []byte) (string, error) {
		if u != userID || n != name || k != kind {
			t.Fatalf("unexpected args: user=%s name=%s kind=%d", u, n, k)
		}
		if string(content) != string(contentRaw) {
			t.Fatalf("unexpected content: got %q want %q", string(content), string(contentRaw))
		}
		if string(metadata) != string(metadataRaw) {
			t.Fatalf("unexpected metadata: got %q want %q", string(metadata), string(metadataRaw))
		}
		return "new-id", nil
	}}

	svc := NewSecretService(repo, logger)
	req := model.CreateSecretRequest{Name: name, Data: contentB64, Metadata: metadataB64, Kind: kind}
	if err := svc.CreateSecret(ctx, userID, req); err != nil {
		t.Fatalf("CreateSecret returned error: %v", err)
	}
	if repo.CreateCalled != 1 {
		t.Fatalf("expected repo CreateSecret to be called once, got %d", repo.CreateCalled)
	}
}

func TestSecretService_CreateSecret_DecodeDataError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{}

	svc := NewSecretService(repo, logger)
	req := model.CreateSecretRequest{Name: "n", Data: "!!!not-base64!!!", Metadata: base64.StdEncoding.EncodeToString([]byte("m")), Kind: 1}
	if err := svc.CreateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.CreateCalled != 0 {
		t.Fatalf("repo CreateSecret should not be called when data decode fails")
	}
}

func TestSecretService_CreateSecret_DecodeMetadataError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{}

	svc := NewSecretService(repo, logger)
	req := model.CreateSecretRequest{Name: "n", Data: base64.StdEncoding.EncodeToString([]byte("d")), Metadata: "!!!not-base64!!!", Kind: 1}
	if err := svc.CreateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.CreateCalled != 0 {
		t.Fatalf("repo CreateSecret should not be called when metadata decode fails")
	}
}

func TestSecretService_CreateSecret_RepoError(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	repo := &mockSecretRepo{CreateSecretFunc: func(ctx context.Context, userID string, content []byte, name string, kind int, metadata []byte) (string, error) {
		return "", errors.New("insert failed")
	}}

	svc := NewSecretService(repo, logger)
	req := model.CreateSecretRequest{Name: "n", Data: base64.StdEncoding.EncodeToString([]byte("d")), Metadata: base64.StdEncoding.EncodeToString([]byte("m")), Kind: 1}
	if err := svc.CreateSecret(ctx, "u", req); err == nil {
		t.Fatalf("expected error, got nil")
	}
	if repo.CreateCalled != 1 {
		t.Fatalf("expected repo CreateSecret to be called once, got %d", repo.CreateCalled)
	}
}
