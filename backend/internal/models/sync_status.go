package models

type SyncStatus string

const (
	StatusInProgress SyncStatus = "InProgress"
	StatusDone       SyncStatus = "Done"
	StatusErrored    SyncStatus = "Errored"
	StatusNotFound   SyncStatus = "NotFound"
)
