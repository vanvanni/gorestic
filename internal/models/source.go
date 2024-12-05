package models

import (
	"time"

	"gorm.io/gorm"
)

type Source struct {
	ID           uint          `gorm:"primarykey"`
	Name         string        `gorm:"not null"`
	StatReports  []StatReport  `gorm:"foreignKey:SourceID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
