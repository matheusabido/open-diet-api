package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/repository"
)

var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type TokenService struct {
	repo      *repository.RefreshTokenRepository
	jwtSecret []byte
}

func NewTokenService(repo *repository.RefreshTokenRepository, jwtSecret string) *TokenService {
	return &TokenService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *TokenService) Issue(ctx context.Context, userID uint64) (*TokenPair, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return nil, err
	}
	refreshToken := base64.RawURLEncoding.EncodeToString(raw)

	err = s.repo.Create(ctx, &model.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: now.Add(refreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenTTL.Seconds()),
	}, nil
}

func (s *TokenService) Consume(ctx context.Context, refreshToken string) (uint64, error) {
	tokenHash := hashToken(refreshToken)

	userID, err := s.repo.RevokeActive(ctx, tokenHash)
	if err == nil {
		return userID, nil
	}
	if !errors.Is(err, repository.ErrRefreshTokenNotFound) {
		return 0, err
	}

	existing, err := s.repo.FindByHash(ctx, tokenHash)
	if errors.Is(err, repository.ErrRefreshTokenNotFound) {
		return 0, ErrInvalidRefreshToken
	}
	if err != nil {
		return 0, err
	}

	if existing.RevokedAt != nil {
		if err := s.repo.RevokeAllByUser(ctx, existing.UserID); err != nil {
			return 0, err
		}
	}

	return 0, ErrInvalidRefreshToken
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
