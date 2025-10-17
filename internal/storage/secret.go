package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"go.uber.org/zap"
)

type SecretRepository struct {
	db  *Database
	log *zap.Logger
}

func NewSecretRepository(db *Database, log *zap.Logger) *SecretRepository {
	return &SecretRepository{db: db, log: log}
}

func (s *SecretRepository) GetSecret(ctx context.Context, userID string, name string) (model.SecretDB, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	secret := model.SecretDB{}
	err := s.db.Dbpool.QueryRow(ctx, "select secret_id, kind, data, metadata from secrets where owner_id = $1 and name = $2", userID, name).
		Scan(&secret.ID, &secret.Kind, &secret.Data, &secret.Metadata)
	if errors.Is(err, pgx.ErrNoRows) {
		s.log.Error("no secret found")
		return model.SecretDB{}, model.ErrSecretNoFound
	}

	if err != nil {
		s.log.Error("error when select secret", zap.String("userID", userID), zap.String("name", name), zap.Error(err))
		return model.SecretDB{}, err
	}

	return secret, nil
}

func (s *SecretRepository) UpdateSecret(ctx context.Context, secretID string, content []byte, kind int, metadata []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	tx, err := s.db.Dbpool.Begin(ctx)
	if err != nil {
		s.log.Error("Error open tran", zap.Error(err))
		return err
	}

	defer func() {
		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				s.log.Error("Error commit tran", zap.Error(err))
			}
		} else {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				s.log.Error("Error rollback tx", zap.Error(errRollback))
			}
		}
	}()

	_, err = tx.Exec(ctx, "update secrets set kind = $1, data = $2, metadata = $3 where secret_id = $4", kind, content, metadata, secretID)
	if err != nil {
		s.log.Error("error when update secret", zap.String("secretID", secretID), zap.Error(err))
		return err
	}

	return nil
}

func (s *SecretRepository) CreateSecret(ctx context.Context, userID string, content []byte, name string, kind int, metadata []byte) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	tx, err := s.db.Dbpool.Begin(ctx)
	if err != nil {
		s.log.Error("Error open tran", zap.Error(err))
		return "", err
	}

	defer func() {
		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				s.log.Error("Error commit tran", zap.Error(err))
			}
		} else {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				s.log.Error("Error rollback tx", zap.Error(errRollback))
			}
		}
	}()

	secretID := uuid.NewString()
	_, err = tx.Exec(ctx, "INSERT INTO secrets (secret_id, owner_id, name, kind, metadata, data) VALUES($1, $2, $3, $4, $5, $6)", secretID, userID, name, kind, metadata, content)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				s.log.Error("duplicate secret %s", zap.Error(err))
				return "", model.ErrSecretConflict
			}
		}

		return "", err
	}

	return secretID, nil
}
