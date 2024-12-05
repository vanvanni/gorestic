package models

type APIKeySource struct {
	APIKeyID uint   `gorm:"primarykey"`
	SourceID uint   `gorm:"primarykey"`
	APIKey   APIKey `gorm:"foreignKey:APIKeyID"`
	Source   Source `gorm:"foreignKey:SourceID"`
}

func (APIKeySource) TableName() string {
	return "api_key_sources"
}
