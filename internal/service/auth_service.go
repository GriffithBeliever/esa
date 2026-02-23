package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/solomon/ims/internal/domain"
	"github.com/solomon/ims/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         *domain.User `json:"user"`
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input LoginInput) (*AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error)
	Logout(ctx context.Context, userID string) error
	ValidateAccessToken(tokenStr string) (userID string, err error)
}

type authService struct {
	userRepo      repository.UserRepository
	jwtSecret     []byte
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry, refreshExpiry time.Duration) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtSecret:     []byte(jwtSecret),
		jwtExpiry:     jwtExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	if input.Email == "" || input.Username == "" || len(input.Password) < 6 {
		return nil, fmt.Errorf("%w: email, username and password (min 6 chars) are required", domain.ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	u := &domain.User{
		ID:           uuid.NewString(),
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrConflict, err.Error())
	}

	return u, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthTokens, error) {
	u, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	return s.issueTokens(ctx, u)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	hash := hashToken(refreshToken)
	id, userID, expiresAt, revoked, err := s.userRepo.FindRefreshToken(ctx, hash)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	if revoked || time.Now().After(expiresAt) {
		return nil, domain.ErrExpired
	}

	if err := s.userRepo.RevokeRefreshToken(ctx, id); err != nil {
		return nil, err
	}

	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, u)
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	return s.userRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *authService) ValidateAccessToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", domain.ErrUnauthorized
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}

	return userID, nil
}

func (s *authService) issueTokens(ctx context.Context, u *domain.User) (*AuthTokens, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": u.ID,
		"iat": now.Unix(),
		"exp": now.Add(s.jwtExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	rawRefresh := make([]byte, 32)
	if _, err := rand.Read(rawRefresh); err != nil {
		return nil, err
	}
	refreshToken := hex.EncodeToString(rawRefresh)
	refreshHash := hashToken(refreshToken)
	refreshID := uuid.NewString()
	expiresAt := now.Add(s.refreshExpiry)

	if err := s.userRepo.StoreRefreshToken(ctx, refreshID, u.ID, refreshHash, expiresAt); err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.jwtExpiry.Seconds()),
		User:         u,
	}, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
