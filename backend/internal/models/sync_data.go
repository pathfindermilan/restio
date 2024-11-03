package models

import (
	"github.com/jinzhu/gorm"
)

type SyncData struct {
	gorm.Model
	UserID       uint   `gorm:"unique;not null" json:"user_id"`
	ContentType  string `gorm:"not null"`
	Age          int    `gorm:"not null"`
	QueryText    string `gorm:"not null"`
	FeelingLevel string `gorm:"not null"` // happy, sad, angry, confused, confident, tired

	ImageURL    string
	DocumentURL string
}
