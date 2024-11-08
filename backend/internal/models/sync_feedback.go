package models

type SyncFeedback struct {
	ID             uint `gorm:"primary_key"`
	UserID         uint `gorm:"uniqueIndex"`
	SyncDataID     uint
	HumanFeedback  string
	HumanReaction  string
	NewContentType string
	QuestionsRelly int
}
