package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
)

type apiTokenServiceImpl struct {
	tokenRepo repository.APITokenRepository
}

func NewAPITokenService(tokenRepo repository.APITokenRepository) APITokenService {
	return &apiTokenServiceImpl{tokenRepo: tokenRepo}
}

func (s *apiTokenServiceImpl) Create(ctx context.Context, adminID uint, name string, expiresAt interface{}) (*model.APIToken, string, error) {
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", fmt.Errorf("生成 Token 失败: %w", err)
	}
	rawToken := hex.EncodeToString(rawBytes)

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	var expiresAtTime *time.Time
	if expiresAt != nil {
		switch v := expiresAt.(type) {
		case time.Time:
			expiresAtTime = &v
		case *time.Time:
			expiresAtTime = v
		case string:
			if v != "" {
				parsed, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, "", fmt.Errorf("无效的过期时间格式: %w", err)
				}
				expiresAtTime = &parsed
			}
		}
	}

	token := &model.APIToken{
		Name:      name,
		TokenHash: tokenHash,
		AdminID:   adminID,
		ExpiresAt: expiresAtTime,
	}
	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return nil, "", err
	}
	return token, rawToken, nil
}

func (s *apiTokenServiceImpl) GetByToken(ctx context.Context, token string) (*model.APIToken, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	apiToken, err := s.tokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("无效的 API Token")
	}
	if apiToken.ExpiresAt != nil && apiToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("API Token 已过期")
	}
	return apiToken, nil
}

func (s *apiTokenServiceImpl) List(ctx context.Context, adminID uint) ([]model.APIToken, error) {
	return s.tokenRepo.List(ctx, adminID)
}

func (s *apiTokenServiceImpl) Delete(ctx context.Context, id uint, adminID uint) error {
	token, err := s.tokenRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if token.AdminID != adminID {
		return fmt.Errorf("无权删除此 Token")
	}
	return s.tokenRepo.Delete(ctx, id)
}
