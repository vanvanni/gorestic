package models

import (
	"time"

	"gorm.io/gorm"
)

type StatReport struct {
	ID        uint   `gorm:"primarykey"`
	SourceID  uint   `gorm:"not null"`
	Data      []byte `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
