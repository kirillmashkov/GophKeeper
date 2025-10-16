package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"go.uber.org/zap"
)

type UserRepository struct {
	db  *Database
	log *zap.Logger
}

const timeoutOperationDB = 1 * time.Second

func NewUserRepository(db *Database, log *zap.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

// PutUser сохраняет учетные данные пользователя в базу данных
func (s *UserRepository) PutUser(ctx context.Context, email string, passwordHash string) (string, error) {
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

	id := uuid.NewString()
	_, err = tx.Exec(ctx, "INSERT INTO users (user_id, email, password_hash) VALUES($1, $2, $3)", id, email, passwordHash)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				s.log.Error("duplicate email %s", zap.Error(err))
				return "", model.ErrUserConflict
			}
		}

		return "", err
	}

	return id, err
}

// GetUser возвращает ID пользователя с указанными учетными данными
func (s *UserRepository) GetUser(ctx context.Context, email string) (string, string, error) {
	var id string
	var password_hash string
	err := s.db.Dbpool.QueryRow(ctx, "SELECT user_id, password_hash FROM users WHERE email = $1", email).Scan(&id, &password_hash)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.NoData {
				s.log.Error("no user found %s", zap.String("user email", email))
				return "", "", model.ErrUserNotFound
			}
		}
	}

	return id, password_hash, err
}

