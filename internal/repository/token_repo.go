package repository

import (
	"errors"

	"github.com/mnabielap/starter-kit-restapi-gofasthttp/internal/model"
	"gorm.io/gorm"
)

type tokenRepo struct {
	db *gorm.DB
}

// NewTokenRepository creates a new instance of TokenRepository
func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Create(token *model.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenRepo) FindOne(tokenString string, tokenType string, blacklisted bool) (*model.Token, error) {
	var token model.Token
	err := r.db.Where("token = ? AND type = ? AND blacklisted = ?", tokenString, tokenType, blacklisted).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepo) Delete(tokenString string, tokenType string) error {
	return r.db.Where("token = ? AND type = ?", tokenString, tokenType).Delete(&model.Token{}).Error
}

func (r *tokenRepo) DeleteByUserID(userID uint, tokenType string) error {
	return r.db.Where("user_id = ? AND type = ?", userID, tokenType).Delete(&model.Token{}).Error
}