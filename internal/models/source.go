package models

import (
	"time"

	"gorm.io/gorm"
)

type Source struct {
	ID           uint          `gorm:"primarykey"`
	Name         string        `gorm:"not null"`
	APIKeys      []APIKey      `gorm:"many2many:api_key_sources;"`
	StatReports  []StatReport  `gorm:"foreignKey:SourceID"`
	CheckReports []CheckReport `gorm:"foreignKey:SourceID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
