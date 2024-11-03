package models

import "github.com/jinzhu/gorm"

type SyncStatus string

const (
	StatusInProgress SyncStatus = "inProgress"
	StatusDone       SyncStatus = "done"
	StatusErrored    SyncStatus = "errored"
	StatusNotFound   SyncStatus = "notFound"
)

type SyncDescription struct {
	gorm.Model
	UserID          uint       `gorm:"not null"`
	SyncDataID      uint       `gorm:"not null"`
	ImageSummary    string     `gorm:"type:text"`
	DocumentSummary string     `gorm:"type:text"`
	ImageStatus     SyncStatus `gorm:"not null"`
	DocumentStatus  SyncStatus `gorm:"not null"`
}
