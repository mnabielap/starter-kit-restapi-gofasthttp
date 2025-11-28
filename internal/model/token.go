package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	TokenTypeAccess        = "access"
	TokenTypeRefresh       = "refresh"
	TokenTypeResetPassword = "resetPassword"
	TokenTypeVerifyEmail   = "verifyEmail"
)

// Token represents authentication tokens in the database
type Token struct {
	gorm.Model
	Token       string    `json:"token" gorm:"index;not null"`
	UserID      uint      `json:"userId" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"`
	Expires     time.Time `json:"expires" gorm:"not null"`
	Blacklisted bool      `json:"blacklisted" gorm:"default:false"`
	
	// Foreign Key Relation
	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}