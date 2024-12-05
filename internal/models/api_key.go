package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID        uint     `gorm:"primarykey"`
	Key       string   `gorm:"uniqueIndex;not null"`
	Sources   []Source `gorm:"many2many:api_key_sources;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
