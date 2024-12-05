package models

import "gorm.io/gorm"

// Database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&APIKey{},
		&Source{},
		&StatReport{},
		&CheckReport{},
	)
}
