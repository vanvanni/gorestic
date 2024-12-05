package repository

import (
	"context"

	"github.com/vanvanni/gorestic/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAPIKey(ctx context.Context, key string) (*models.APIKey, error) {
	apiKey := &models.APIKey{Key: key}
	result := r.db.WithContext(ctx).Create(apiKey)
	return apiKey, result.Error
}

func (r *Repository) AssignKeyToSource(ctx context.Context, keyID, sourceID uint) error {
	return r.db.WithContext(ctx).Model(&models.APIKey{ID: keyID}).
		Association("Sources").
		Append(&models.Source{ID: sourceID})
}

func (r *Repository) CreateSource(ctx context.Context, name string) (*models.Source, error) {
	source := &models.Source{Name: name}
	result := r.db.WithContext(ctx).Create(source)
	return source, result.Error
}

func (r *Repository) CreateStatReport(ctx context.Context, sourceID uint, data []byte) (*models.StatReport, error) {
	report := &models.StatReport{
		SourceID: sourceID,
		Data:     data,
	}
	result := r.db.WithContext(ctx).Create(report)
	return report, result.Error
}

func (r *Repository) CreateCheckReport(ctx context.Context, sourceID uint, status string, details []byte) (*models.CheckReport, error) {
	report := &models.CheckReport{
		SourceID: sourceID,
		Status:   status,
		Details:  details,
	}
	result := r.db.WithContext(ctx).Create(report)
	return report, result.Error
}

func (r *Repository) GetSourcesByAPIKey(ctx context.Context, keyID uint) ([]models.Source, error) {
	var sources []models.Source
	err := r.db.WithContext(ctx).
		Model(&models.APIKey{ID: keyID}).
		Association("Sources").
		Find(&sources)
	return sources, err
}

func (r *Repository) GetReportsBySource(ctx context.Context, sourceID uint) ([]models.StatReport, []models.CheckReport, error) {
	var statReports []models.StatReport
	var checkReports []models.CheckReport

	err := r.db.WithContext(ctx).Where("source_id = ?", sourceID).Find(&statReports).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.WithContext(ctx).Where("source_id = ?", sourceID).Find(&checkReports).Error
	return statReports, checkReports, err
}

func (r *Repository) GetAPIKeyByKey(ctx context.Context, key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&apiKey)
	if result.Error != nil {
		return nil, result.Error
	}
	return &apiKey, nil
}
