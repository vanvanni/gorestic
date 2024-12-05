package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vanvanni/gorestic/internal/config"
	"github.com/vanvanni/gorestic/internal/repository"
)

type Handler struct {
	config *config.Config
	repo   *repository.Repository
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
		repo:   repository.NewRepository(cfg.DB),
	}
}

func (h *Handler) HandleGetStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}

func (h *Handler) HandleUpdateStats(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")
	if apiKey == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "API key required",
		})
	}

	key, err := h.repo.GetAPIKeyByKey(c.Context(), apiKey)
	if err != nil || key == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid API key",
		})
	}

	// TODO:

	var input interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// stats := storage.BackupStats{
	// 	TotalSize:      input.TotalSize,
	// 	TotalFileCount: input.TotalFileCount,
	// 	SnapshotsCount: input.SnapshotsCount,
	// 	CreatedAt:      time.Now(),
	// 	APIKeyName:     keyName,
	// }

	// if err := h.storage.AddStats(stats); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"error": "Failed to save stats",
	// 	})
	// }

	return c.JSON(fiber.Map{
		"message": "Stats updated successfully",
	})
}
