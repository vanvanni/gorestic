package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vanvanni/gorestic/internal/storage"
)

type Handler struct {
	storage *storage.Manager
}

func NewHandler(store *storage.Manager) *Handler {
	return &Handler{
		storage: store,
	}
}

func (h *Handler) HandleDashboard(c *fiber.Ctx) error {
	return c.Render("views/index", fiber.Map{
		"Title": "GORestic",
	}, "views/layouts/app")
}
