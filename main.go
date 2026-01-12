package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/G0tem/go-servise-auth/docs" // swagger docs
	"github.com/G0tem/go-servise-auth/internal/config"
	"github.com/G0tem/go-servise-auth/internal/database"
	"github.com/G0tem/go-servise-auth/internal/handler"
	"github.com/G0tem/go-servise-auth/internal/handler/rbac"
	"github.com/G0tem/go-servise-auth/internal/model"
	"github.com/G0tem/go-servise-auth/internal/queue"
	"github.com/G0tem/go-servise-auth/internal/router"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Local-Template-Auth Swagger
// @version 1.0
// @description This is an API of auth-service
// @schemes http https

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

// @BasePath /api/v1
func main() {
	// Initialize Zerolog logger with output to stdout
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.LoadConfig()
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	db, err := database.Connect(cfg)
	if err != nil {
		return
	}

	app := fiber.New(fiber.Config{
		BodyLimit:         cfg.MaxFileUploadSizeInBytes,
		StreamRequestBody: true,
	})

	swaggerCfg := swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/swagger.yaml",
		Path:     "docs",
		CacheAge: 1,
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(swagger.New(swaggerCfg))
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH")
		c.Set("Access-Control-Allow-Headers", "*")
		c.Set("Access-Control-Expose-Headers", "*")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "*",
		MaxAge:           86400, // 24 часов в секундах
	}))

	rbac := &rbac.RBACLayer{
		DB:  db,
		Ctx: context.Background(),
	}
	err = rbac.InitSafety(map[string]string{
		model.AdminRole:       "admin:all",
		model.DefaultUserRole: "user:read",
	})
	if err != nil {
		log.Error().Msgf("Setup roles error: %v", err)
		return
	}

	handlers := handler.NewHandler(db, queue.NewMailQueueConnector(&cfg), rbac, &cfg)

	router.SetupRoutes(app)
	handlers.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	err = app.Listen(fmt.Sprintf(":%v", cfg.HttpPort))
	if err != nil {
		log.Error().Msgf("Unexpected error: %v", err)
		return
	}
}
