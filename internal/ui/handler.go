package ui

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/vanvanni/gorestic/internal/config"
	"github.com/vanvanni/gorestic/internal/repository"
)

type (
	PostKeyRequest struct {
		Name        string `json:"name" form:"name"`
		Key         string `json:"key" form:"key"`
		Description string `json:"description" form:"description"`
	}
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

func (h *Handler) HandleGetKeys(c *fiber.Ctx) error {
	rows, err := h.repo.GetAllAPIKeys()
	if err != nil {
		return c.JSON(fiber.Map{
			"error": "Could not fetch records",
		})
	}

	return c.JSON(fiber.Map{
		"data": rows,
	})
}

func (h *Handler) HandlePostKey(c *fiber.Ctx) error {
	r := new(PostKeyRequest)

	if err := json.Unmarshal(c.BodyRaw(), r); err != nil {
		return c.JSON(fiber.Map{
			"error": "Could not parse request",
		})
	}

	model, err2 := h.repo.CreateAPIKey(r.Name, r.Key, r.Description)
	if err2 != nil {
		return c.JSON(fiber.Map{
			"error": "Could not create key",
		})
	}

	return c.JSON(fiber.Map{
		"data": model,
	})
}
