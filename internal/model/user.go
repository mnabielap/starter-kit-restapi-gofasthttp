package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity in the database
type User struct {
	gorm.Model
	Name            string `json:"name" gorm:"not null"`
	Email           string `json:"email" gorm:"uniqueIndex;not null"`
	Password        string `json:"password" gorm:"not null"` 
	Role            string `json:"role" gorm:"default:'user'"`
	IsEmailVerified bool   `json:"isEmailVerified" gorm:"default:false"`
}

// UserResponse is a DTO for sending user data to the client safely
type UserResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
}

// ToResponse converts a User model to UserResponse DTO
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Role:            u.Role,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
	}
}