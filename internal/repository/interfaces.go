package repository

import (
	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
)

// UserFilter contains all possible filters for querying users
type UserFilter struct {
	Page   int
	Limit  int
	Search string
	Scope  string
	Role   string
	SortBy string
}

// UserRepository defines the methods for user database operations
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindAll(filter UserFilter) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
	IsEmailTaken(email string, excludeID uint) (bool, error)
}

// TokenRepository defines the methods for token database operations
type TokenRepository interface {
	Create(token *model.Token) error
	FindOne(token string, tokenType string, blacklisted bool) (*model.Token, error)
	Delete(token string, tokenType string) error
	DeleteByUserID(userID uint, tokenType string) error
}