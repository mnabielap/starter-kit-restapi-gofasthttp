package service

import (
	"errors"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	// Check if email exists
	exists, err := s.userRepo.IsEmailTaken(user.Email, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already taken")
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Set Defaults
	if user.Role == "" {
		user.Role = "user"
	}

	err = s.userRepo.Create(user)
	return user, err
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetUsers(page, limit int) ([]model.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Convert to DTOs
	var response []model.UserResponse
	for _, u := range users {
		response = append(response, u.ToResponse())
	}

	return response, total, nil
}

func (s *UserService) UpdateUser(id uint, updateData map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if email, ok := updateData["email"].(string); ok && email != user.Email {
		taken, err := s.userRepo.IsEmailTaken(email, id)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, errors.New("email already taken")
		}
		user.Email = email
	}

	if name, ok := updateData["name"].(string); ok {
		user.Name = name
	}

	if password, ok := updateData["password"].(string); ok {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	err = s.userRepo.Update(user)
	return user, err
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}