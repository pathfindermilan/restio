package models

type AISummary struct {
	ID                    uint `gorm:"primary_key"`
	UserID                uint `gorm:"uniqueIndex"`
	SyncDataID            uint
	AriaResponse          string
	AriaResponseStatus    SyncStatus
	AllegroResponse       string
	AllegroResponseStatus SyncStatus
}
