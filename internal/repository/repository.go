package repository

import (
	"context"
	"fmt"

	"github.com/vanvanni/gorestic/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAPIKey(name string, key string, description string) (*models.APIKey, error) {
	apiKey := &models.APIKey{Name: name, Key: key, Description: description}
	result := r.db.Create(apiKey)
	return apiKey, result.Error
}

func (r *Repository) GetAllAPIKeys() ([]models.APIKey, error) {
	var keys []models.APIKey
	result := r.db.Find(&keys)
	fmt.Println("Result: ", result)
	return keys, result.Error
}

func (r *Repository) GetAPIKeyByKey(ctx context.Context, key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&apiKey)
	if result.Error != nil {
		return nil, result.Error
	}
	return &apiKey, nil
}
