package service

import (
	"context"
	"fmt"

	"github.com/netconfighub/netconfighub/internal/model"
	"github.com/netconfighub/netconfighub/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type adminServiceImpl struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminServiceImpl{adminRepo: adminRepo}
}

func (s *adminServiceImpl) Create(ctx context.Context, username, password string) (*model.Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}
	admin := &model.Admin{
		Username:     username,
		PasswordHash: string(hash),
	}
	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *adminServiceImpl) Authenticate(ctx context.Context, username, password string) (*model.Admin, error) {
	admin, err := s.adminRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	return admin, nil
}

func (s *adminServiceImpl) GetByID(ctx context.Context, id uint) (*model.Admin, error) {
	return s.adminRepo.FindByID(ctx, id)
}

func (s *adminServiceImpl) List(ctx context.Context, offset, limit int) ([]model.Admin, int64, error) {
	return s.adminRepo.List(ctx, offset, limit)
}

func (s *adminServiceImpl) Update(ctx context.Context, id uint, username string, password *string) (*model.Admin, error) {
	admin, err := s.adminRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	admin.Username = username
	if password != nil && *password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("密码哈希失败: %w", err)
		}
		admin.PasswordHash = string(hash)
	}
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *adminServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.adminRepo.Delete(ctx, id)
}
