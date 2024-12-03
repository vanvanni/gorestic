package storage

import (
	"time"
)

type BackupStats struct {
	TotalSize      int64     `json:"total_size"`
	TotalFileCount int       `json:"total_file_count"`
	SnapshotsCount int       `json:"snapshots_count"`
	CreatedAt      time.Time `json:"created_at"`
	APIKeyName     string    `json:"api_key_name"`
}

type Storage struct {
	Stats []BackupStats `json:"stats"`
}