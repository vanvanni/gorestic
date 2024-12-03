package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/vanvanni/gorestic/internal/api"
	"github.com/vanvanni/gorestic/internal/config"
	"github.com/vanvanni/gorestic/internal/storage"
	"github.com/vanvanni/gorestic/internal/web"
)

//go:embed views/*.html
var viewsFS embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	store, err := storage.NewManager(cfg.Storage.Path)
	if err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}

	apiHandler := api.NewHandler(cfg, store)
	webHandler := web.NewHandler(store)
	engine := html.NewFileSystem(http.FS(viewsFS), ".html")

	engine.AddFunc("formatBytes", func(bytes int64) string {
		const unit = 1024
		if bytes < unit {
			return fmt.Sprintf("%d B", bytes)
		}
		div, exp := int64(unit), 0
		for n := bytes / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
	})

	engine.AddFunc("formatTime", func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	})
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	apiGroup := app.Group("/api")
	apiGroup.Get("/stats", apiHandler.HandleGetStats)
	apiGroup.Post("/stats", apiHandler.HandleUpdateStats)

	webGroup := app.Group("/")
	webGroup.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.Server.Username: cfg.Server.Password,
		},
	}))
	webGroup.Get("/", webHandler.HandleDashboard)

	log.Printf("Starting server on port %d", cfg.Server.Port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)))
}
