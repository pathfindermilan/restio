package models

type SyncDescription struct {
	ID              uint `gorm:"primary_key"`
	UserID          uint `gorm:"uniqueIndex"`
	SyncDataID      uint
	ImageURL        string
	ImageStatus     SyncStatus
	ImageSummary    string
	DocumentURL     string
	DocumentStatus  SyncStatus
	DocumentSummary string
}
