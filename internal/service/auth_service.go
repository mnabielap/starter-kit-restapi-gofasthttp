package service

import (
	"errors"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService  *UserService
	tokenService *TokenService
}

func NewAuthService(userService *UserService, tokenService *TokenService) *AuthService {
	return &AuthService{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *AuthService) Login(email, password string) (*model.User, *AuthTokens, error) {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, nil, errors.New("incorrect email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, errors.New("incorrect email or password")
	}

	tokens, err := s.tokenService.GenerateAuthTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) Register(user *model.User) (*model.User, *AuthTokens, error) {
	createdUser, err := s.userService.CreateUser(user)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.tokenService.GenerateAuthTokens(createdUser)
	if err != nil {
		return nil, nil, err
	}

	return createdUser, tokens, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenService.DeleteRefreshToken(refreshToken)
}

func (s *AuthService) RefreshAuth(refreshToken string) (*AuthTokens, error) {
	// Verify token signature and type
	claims, err := s.tokenService.VerifyToken(refreshToken)
	if err != nil || claims.Type != model.TokenTypeRefresh {
		return nil, errors.New("invalid refresh token")
	}

	// Verify token exists in DB (not blacklisted/deleted)
	storedToken, err := s.tokenService.FindRefreshToken(refreshToken)
	if err != nil || storedToken == nil {
		return nil, errors.New("refresh token not found or reused")
	}

	// Get User
	user, err := s.userService.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Delete old token
	s.tokenService.DeleteRefreshToken(refreshToken)

	// Generate new pair
	return s.tokenService.GenerateAuthTokens(user)
}