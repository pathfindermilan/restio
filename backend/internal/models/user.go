package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username         string `gorm:"unique;not null" json:"username" validate:"required,min=8,max=20"`
	Email            string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password         string `gorm:"not null" json:"password" validate:"required,password"`
	Name             string `json:"name"`
	VerificationCode string `json:"verification_code"`
	IsVerified       bool   `json:"is_verified"`
}
