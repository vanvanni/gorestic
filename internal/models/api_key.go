package models

import (
	"time"

	"gorm.io/gorm"
)

type APIKey struct {
	ID          uint     `gorm:"primarykey"`
	Name        string   `gorm:"varchar(255)"`
	Key         string   `gorm:"uniqueIndex;not null"`
	Description string   `gorm:"varchar(255)"`
	Sources     []Source `gorm:"many2many:api_key_sources;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
