package route

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/unitechio/oss-monorepo/server/internal/config"
	"github.com/unitechio/oss-monorepo/server/internal/http/middleware"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

type mockUser struct {
	ID          int      `json:"id"`
	UUID        string   `json:"-"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Password    string   `json:"-"`
	Status      string   `json:"status"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type authTokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         mockUser  `json:"user"`
	ExpiresAt    time.Time `json:"-"`
}

type mockAuthService struct {
	secret        []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	usersByName   map[string]mockUser
	usersByUUID   map[string]mockUser
	refreshTokens map[string]string
	mu            sync.RWMutex
}

func newMockAuthService(cfg *config.Config) *mockAuthService {
	users := []mockUser{
		{
			ID:          1,
			UUID:        "e6f1f1d8-4af3-4b22-9c68-5fc9821d01a1",
			Username:    "superadmin",
			FullName:    "Core Super Admin",
			Email:       "superadmin@example.local",
			Phone:       "0900000001",
			Password:    "Admin@123",
			Status:      "active",
			Roles:       []string{"superadmin"},
			Permissions: []string{"*"},
		},
		{
			ID:          2,
			UUID:        "78f9f9e4-3fc9-4cc1-9b32-8c0700f6de51",
			Username:    "admin",
			FullName:    "Core Administrator",
			Email:       "admin@example.local",
			Phone:       "0900000002",
			Password:    "Admin@123",
			Status:      "active",
			Roles:       []string{"admin"},
			Permissions: []string{"user.read:global", "user.create:global", "user.update:global", "role.read:global", "role.update:global", "menu.read:global", "permission.read:global", "auth.read:global", "audit.read:global", "device.read:global", "client.read:global", "client.create:global", "client.update:global", "channel.read:global", "channel.create:global", "channel.update:global", "policy.read:global", "policy.create:global", "policy.update:global", "option.read:global", "option.create:global", "option.update:global", "service.read:global", "service.create:global", "service.update:global"},
		},
		{
			ID:          3,
			UUID:        "ba3a6962-e78c-49f3-8a18-8ebdbd4aeff3",
			Username:    "operator",
			FullName:    "Core Operator",
			Email:       "operator@example.local",
			Phone:       "0900000003",
			Password:    "Admin@123",
			Status:      "active",
			Roles:       []string{"operator"},
			Permissions: []string{"user.read:self", "auth.read:self", "device.read:self", "policy.read:self", "option.read:self"},
		},
	}

	service := &mockAuthService{
		secret:        []byte(cfg.JWT.Secret),
		accessTTL:     cfg.JWT.AccessTokenExpire,
		refreshTTL:    cfg.JWT.RefreshTokenExpire,
		usersByName:   make(map[string]mockUser, len(users)),
		usersByUUID:   make(map[string]mockUser, len(users)),
		refreshTokens: make(map[string]string),
	}

	if len(service.secret) == 0 {
		service.secret = []byte("change-this-secret-key")
	}
	if service.accessTTL <= 0 {
		service.accessTTL = 15 * time.Minute
	}
	if service.refreshTTL <= 0 {
		service.refreshTTL = 7 * 24 * time.Hour
	}

	for _, user := range users {
		service.usersByName[strings.ToLower(user.Username)] = user
		service.usersByUUID[user.UUID] = user
	}

	return service
}

func (s *mockAuthService) registerRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	auth.POST("/login", s.handleLogin)
	auth.POST("/refresh", s.handleRefresh)
	auth.POST("/logout", s.handleLogout)
	auth.GET("/sso/providers", s.handleSSOProviders)
	auth.GET("/me", middleware.JWTAuth(string(s.secret)), s.handleMe)
	auth.PUT("/change-password", middleware.JWTAuth(string(s.secret)), s.handleChangePassword)
}

func (s *mockAuthService) handleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailError(c, apperr.BadRequest("invalid login payload", err))
		return
	}

	user, ok := s.usersByName[strings.ToLower(strings.TrimSpace(req.Username))]
	if !ok || user.Password != req.Password {
		response.FailError(c, apperr.Unauthorized("invalid username or password"))
		return
	}

	pair, err := s.issueTokens(user)
	if err != nil {
		response.FailError(c, apperr.Internal(fmt.Errorf("issue tokens: %w", err)))
		return
	}

	response.OK(c, "login successful", pair)
}

func (s *mockAuthService) handleRefresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		response.FailError(c, apperr.BadRequest("refresh_token is required", err))
		return
	}

	user, err := s.userFromRefreshToken(req.RefreshToken)
	if err != nil {
		response.FailError(c, apperr.Unauthorized("invalid refresh token", err))
		return
	}

	pair, err := s.issueTokens(user)
	if err != nil {
		response.FailError(c, apperr.Internal(fmt.Errorf("issue tokens: %w", err)))
		return
	}

	response.OK(c, "token refreshed", pair)
}

func (s *mockAuthService) handleLogout(c *gin.Context) {
	response.OK(c, "logged out", gin.H{})
}

func (s *mockAuthService) handleSSOProviders(c *gin.Context) {
	response.OK(c, "sso providers loaded", []gin.H{})
}

func (s *mockAuthService) handleMe(c *gin.Context) {
	user, err := s.userFromContext(c)
	if err != nil {
		response.FailError(c, apperr.Unauthorized("invalid token context", err))
		return
	}

	response.OK(c, "profile loaded", user)
}

func (s *mockAuthService) handleChangePassword(c *gin.Context) {
	user, err := s.userFromContext(c)
	if err != nil {
		response.FailError(c, apperr.Unauthorized("invalid token context", err))
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailError(c, apperr.BadRequest("invalid change password payload", err))
		return
	}
	if strings.TrimSpace(req.NewPassword) == "" || len(req.NewPassword) < 8 {
		response.FailError(c, apperr.BadRequest("new password must be at least 8 characters"))
		return
	}
	if req.OldPassword != user.Password {
		response.FailError(c, apperr.Unauthorized("current password is incorrect"))
		return
	}

	user.Password = req.NewPassword
	s.usersByName[strings.ToLower(user.Username)] = user
	s.usersByUUID[user.UUID] = user

	response.OK(c, "password changed", gin.H{})
}

func (s *mockAuthService) userFromContext(c *gin.Context) (mockUser, error) {
	rawID, ok := c.Get(middleware.ContextKeyUserID)
	if !ok {
		return mockUser{}, errors.New("missing user id")
	}
	id, ok := rawID.(uuid.UUID)
	if !ok {
		return mockUser{}, errors.New("invalid user id type")
	}

	user, ok := s.usersByUUID[id.String()]
	if !ok {
		return mockUser{}, errors.New("user not found")
	}
	return user, nil
}

func (s *mockAuthService) userFromRefreshToken(refreshToken string) (mockUser, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return mockUser{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return mockUser{}, errors.New("invalid token claims")
	}
	if claims["type"] != "refresh" {
		return mockUser{}, errors.New("unexpected token type")
	}

	subject, _ := claims.GetSubject()
	s.mu.RLock()
	storedSubject, ok := s.refreshTokens[refreshToken]
	s.mu.RUnlock()
	if !ok || storedSubject != subject {
		return mockUser{}, errors.New("refresh token revoked")
	}

	user, ok := s.usersByUUID[subject]
	if !ok {
		return mockUser{}, errors.New("user not found")
	}
	return user, nil
}

func (s *mockAuthService) issueTokens(user mockUser) (*authTokenPair, error) {
	now := time.Now()

	accessClaims := middleware.JWTClaims{
		UserID: uuid.MustParse(user.UUID),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UUID,
			Issuer:    "oss-monorepo",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.UUID,
		Issuer:    "oss-monorepo",
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
		ID:        uuid.NewString(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  refreshClaims.Subject,
		"iss":  refreshClaims.Issuer,
		"iat":  refreshClaims.IssuedAt.Unix(),
		"nbf":  refreshClaims.NotBefore.Unix(),
		"exp":  refreshClaims.ExpiresAt.Unix(),
		"jti":  refreshClaims.ID,
		"type": "refresh",
	}).SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.refreshTokens[refreshToken] = user.UUID
	s.mu.Unlock()

	return &authTokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    now.Add(s.accessTTL),
	}, nil
}
