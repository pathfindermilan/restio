package models

type SyncDescription struct {
	ID              uint `gorm:"primary_key"`
	UserID          uint `gorm:"uniqueIndex"`
	SyncDataID      uint
	ImageURL        string
	ImageSummary    string
	ImageStatus     SyncStatus
	DocumentURL     string
	DocumentSummary string
	DocumentStatus  SyncStatus
	AiStatus        SyncStatus
	AiSummary       string
}
