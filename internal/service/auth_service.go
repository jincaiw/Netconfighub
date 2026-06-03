package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/netconfighub/netconfighub/internal/middleware"
	"github.com/netconfighub/netconfighub/internal/model"
)

type authServiceImpl struct {
	adminService AdminService
	jwtSecret    string
	jwtExpire    time.Duration
}

func NewAuthService(adminService AdminService, jwtSecret string, jwtExpireHours int) AuthService {
	return &authServiceImpl{
		adminService: adminService,
		jwtSecret:    jwtSecret,
		jwtExpire:    time.Duration(jwtExpireHours) * time.Hour,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, username, password string) (string, time.Time, error) {
	admin, err := s.adminService.Authenticate(ctx, username, password)
	if err != nil {
		return "", time.Time{}, err
	}
	return s.generateToken(admin)
}

func (s *authServiceImpl) Logout(ctx context.Context) error {
	log.Println("用户登出")
	return nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, tokenString string) (string, time.Time, error) {
	claims := &middleware.Claims{}
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", time.Time{}, fmt.Errorf("无效或过期的 JWT Token")
	}

	admin, err := s.adminService.GetByID(ctx, claims.AdminID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("用户不存在")
	}
	return s.generateToken(admin)
}

func (s *authServiceImpl) InitAdmin(ctx context.Context, username, password string) (*model.Admin, error) {
	_, total, err := s.adminService.List(ctx, 0, 1)
	if err != nil {
		return nil, err
	}
	if total > 0 {
		return nil, fmt.Errorf("已存在管理员账户，无法初始化")
	}
	return s.adminService.Create(ctx, username, password)
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, adminID uint, oldPassword, newPassword string) error {
	admin, err := s.adminService.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}
	if _, err := s.adminService.Authenticate(ctx, admin.Username, oldPassword); err != nil {
		return fmt.Errorf("旧密码错误")
	}
	_, err = s.adminService.Update(ctx, adminID, admin.Username, &newPassword)
	return err
}

func (s *authServiceImpl) generateToken(admin *model.Admin) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtExpire)
	claims := &middleware.Claims{
		AdminID:  admin.ID,
		Username: admin.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "netconfighub",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成 JWT Token 失败: %w", err)
	}
	return tokenString, expiresAt, nil
}
