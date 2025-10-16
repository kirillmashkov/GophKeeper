package service

import (
	"context"

	"github.com/kirillmashkov/GophKeeper.git/internal/util"
	"github.com/kirillmashkov/GophKeeper.git/pkg/hasher"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserRepository interface {
	PutUser(ctx context.Context, email string, passwordHash string) (string, error)
	GetUser(ctx context.Context, email string) (string, string, error)
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
	userID, password_hash, err := a.userRepository.GetUser(ctx, email)
	if err != nil {
		a.logger.Error("error get user", zap.String("email", email))
		return false, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	if err != nil {
		a.logger.Error("error verify user", zap.String("email", email))
		return false, "", err
	}

	token, err := a.securityUtil.BuildJWTString(userID)
	if err != nil {
		a.logger.Error("error generate token", zap.Error(err))
		return true, "", nil
	}

	return true, token, nil
}