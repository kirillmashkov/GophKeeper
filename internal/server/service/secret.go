package service

import (
	"context"
	"encoding/base64"

	"github.com/kirillmashkov/GophKeeper.git/internal/server/model"
	"go.uber.org/zap"
)

type ISecretRepository interface {
	CreateSecret(ctx context.Context, userID string, content []byte, name string, kind int, metadata []byte) (string, error)
	GetSecret(ctx context.Context, userID string, name string) (model.SecretDB, error)
	UpdateSecret(ctx context.Context, secretID string, content []byte, kind int, metadata []byte) error
}

type SecretService struct {
	secretRepository ISecretRepository
	logger           *zap.Logger
}

func NewSecretService(sr ISecretRepository, l *zap.Logger) *SecretService {
	return &SecretService{
		secretRepository: sr,
		logger:           l,
	}
}

func (s *SecretService) GetSecret(ctx context.Context, userID string, name string) (model.GetSecretResponse, error) {
	secret, err := s.secretRepository.GetSecret(ctx, userID, name)
	if err != nil {
		return model.GetSecretResponse{}, err
	}

	secretDataEncode := base64.StdEncoding.EncodeToString(secret.Data)
	secretMetadataEncode := base64.StdEncoding.EncodeToString(secret.Metadata)

	secretResponse := model.GetSecretResponse{
		Name: name,
		Data: secretDataEncode,
		Metadata: secretMetadataEncode,
		Kind: secret.Kind,
	}

	return secretResponse, nil
}
 
func (s *SecretService) UpdateSecret(ctx context.Context, userID string, updatesecret model.UpdateSecretRequest) error {
	secret, err := s.secretRepository.GetSecret(ctx, userID, updatesecret.Name)
	if err != nil {
		return err
	}

	decodedBytesContent, err := base64.StdEncoding.DecodeString(updatesecret.Data)
	if err != nil {
		s.logger.Error("error convert data from base64", zap.Error(err))
		return err
	}

	decodedBytesMetadata, err := base64.StdEncoding.DecodeString(updatesecret.Metadata)
	if err != nil {
		s.logger.Error("error convert metaData from base64", zap.Error(err))
		return err
	}

	err = s.secretRepository.UpdateSecret(ctx, secret.ID, decodedBytesContent, updatesecret.Kind, decodedBytesMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretService) CreateSecret(ctx context.Context, userID string, createSecret model.CreateSecretRequest) error {
	decodedBytesContent, err := base64.StdEncoding.DecodeString(createSecret.Data)
	if err != nil {
		s.logger.Error("error convert data from base64", zap.Error(err))
		return err
	}

	decodedBytesMetadata, err := base64.StdEncoding.DecodeString(createSecret.Metadata)
	if err != nil {
		s.logger.Error("error convert metaData from base64", zap.Error(err))
		return err
	}

	_, err = s.secretRepository.CreateSecret(ctx, userID, decodedBytesContent, createSecret.Name, createSecret.Kind, decodedBytesMetadata)
	if err != nil {
		s.logger.Error("error insert new secret", zap.Error(err))
		return err
	}

	return nil
}
