package models

type SyncData struct {
	ID           uint `gorm:"primary_key"`
	UserID       uint `gorm:"uniqueIndex"`
	ContentType  string
	Age          int
	QueryText    string
	FeelingLevel string
	ImageURL     string
	DocumentURL  string
}
