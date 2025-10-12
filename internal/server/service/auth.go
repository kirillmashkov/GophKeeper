package service

import (
	"context"

	"github.com/kirillmashkov/GophKeeper.git/internal/util"
	"github.com/kirillmashkov/GophKeeper.git/pkg/hasher"
	"go.uber.org/zap"
)

type IUserRepository interface {
	PutUser(ctx context.Context, email string, passwordHash string) (string, error)
	GetUser(ctx context.Context, email string, password_hash string) (string, error)
}

type AuthService struct {
	userRepository IUserRepository
	securityUtil *util.SecurityUtil
	logger *zap.Logger
}

func NewAuthService(ur IUserRepository, logger *zap.Logger, su *util.SecurityUtil) *AuthService {
	return &AuthService{userRepository: ur, logger: logger, securityUtil: su}
}

func (a *AuthService) SignUp(ctx context.Context, email string, password string) (string, error) {
	hash, err := hasher.Hash(password)

	if err != nil {
		a.logger.Error("error generated hash", zap.Error(err))
		return "", err
	}

	userID, err := a.userRepository.PutUser(ctx, email, hash)
	if err != nil {
		a.logger.Error("error create user in db", zap.Error(err))
		return "", err
	}
	a.logger.Debug("%s", zap.String("userID", userID))

	token, err := a.securityUtil.BuildJWTString(userID)
	if err != nil {
		a.logger.Error("error generate token", zap.Error(err))
		return "", nil
	}

	return token, nil
}

func (a *AuthService) SignIn(ctx context.Context, email string, password string) (bool, string, error) {
	hash, err := hasher.Hash(password)

	if err != nil {
		a.logger.Error("error generated hash", zap.Error(err))
		return false, "", err
	}

	userID, err := a.userRepository.GetUser(ctx, email, hash)
	if err != nil {
		a.logger.Error("error get user %s", zap.String("email", email))
		return false, "", err
	}

	token, err := a.securityUtil.BuildJWTString(userID)
	if err != nil {
		a.logger.Error("error generate token", zap.Error(err))
		return true, "", nil
	}

	return true, token, nil
}