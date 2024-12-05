package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckReport struct {
	ID        uint   `gorm:"primarykey"`
	SourceID  uint   `gorm:"not null"`
	Status    string `gorm:"not null"`
	Details   []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
