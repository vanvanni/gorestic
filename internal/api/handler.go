package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vanvanni/gorestic/internal/config"
	"github.com/vanvanni/gorestic/internal/storage"
)

type StatsInput struct {
	TotalSize      int64 `json:"total_size"`
	TotalFileCount int   `json:"total_file_count"`
	SnapshotsCount int   `json:"snapshots_count"`
}

type Handler struct {
	config  *config.Config
	storage *storage.Manager
}

func NewHandler(cfg *config.Config, store *storage.Manager) *Handler {
	return &Handler{
		config:  cfg,
		storage: store,
	}
}

func (h *Handler) HandleGetStats(c *fiber.Ctx) error {
	stats := h.storage.GetAllStats()
	return c.JSON(stats)
}

func (h *Handler) HandleUpdateStats(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "API key required",
		})
	}

	// Validate API key
	var foundKey *config.APIKey
	var keyName string
	for name, key := range h.config.APIKeys {
		if key.Key == apiKey {
			foundKey = &key
			keyName = name
			break
		}
	}

	if foundKey == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid API key",
		})
	}

	var input StatsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	stats := storage.BackupStats{
		TotalSize:      input.TotalSize,
		TotalFileCount: input.TotalFileCount,
		SnapshotsCount: input.SnapshotsCount,
		CreatedAt:      time.Now(),
		APIKeyName:     keyName,
	}

	if err := h.storage.AddStats(stats); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save stats",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stats updated successfully",
	})
}
