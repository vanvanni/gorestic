package web

import (
	"github.com/gofiber/fiber/v2"
)

func HandleDashboard(c *fiber.Ctx) error {
	return c.Render("views/index", fiber.Map{
		"Title": "GORestic",
	}, "views/layouts/app")
}

func HandleSources(c *fiber.Ctx) error {
	return c.Render("views/sources", fiber.Map{}, "views/layouts/app")
}

func HandleKeys(c *fiber.Ctx) error {
	return c.Render("views/keys", fiber.Map{}, "views/layouts/app")
}
