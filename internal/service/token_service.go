package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/config"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/repository"
)

type TokenService struct {
	tokenRepo repository.TokenRepository
	config    *config.Config
}

type AuthTokens struct {
	Access  AuthTokenResponse `json:"access"`
	Refresh AuthTokenResponse `json:"refresh"`
}

type AuthTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func NewTokenService(tokenRepo repository.TokenRepository, cfg *config.Config) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		config:    cfg,
	}
}

func (s *TokenService) GenerateToken(userID uint, role, tokenType string, expires time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *TokenService) GenerateAuthTokens(user *model.User) (*AuthTokens, error) {
	accessExpires := time.Now().Add(s.config.JWTAccessExpirationMinutes)
	accessToken, err := s.GenerateToken(user.ID, user.Role, model.TokenTypeAccess, accessExpires)
	if err != nil {
		return nil, err
	}

	refreshExpires := time.Now().Add(s.config.JWTRefreshExpirationDays)
	refreshToken, err := s.GenerateToken(user.ID, user.Role, model.TokenTypeRefresh, refreshExpires)
	if err != nil {
		return nil, err
	}

	// Save Refresh Token to Database
	err = s.tokenRepo.Create(&model.Token{
		Token:   refreshToken,
		UserID:  user.ID,
		Type:    model.TokenTypeRefresh,
		Expires: refreshExpires,
	})
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		Access:  AuthTokenResponse{Token: accessToken, Expires: accessExpires},
		Refresh: AuthTokenResponse{Token: refreshToken, Expires: refreshExpires},
	}, nil
}

func (s *TokenService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *TokenService) DeleteRefreshToken(token string) error {
	return s.tokenRepo.Delete(token, model.TokenTypeRefresh)
}

func (s *TokenService) FindRefreshToken(token string) (*model.Token, error) {
	return s.tokenRepo.FindOne(token, model.TokenTypeRefresh, false)
}