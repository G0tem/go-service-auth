package handler

import (
	"github.com/G0tem/go-service-auth/internal/config"
	"github.com/G0tem/go-service-auth/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	cfg   *config.Config
	redis *redis.Client
}

func NewHandler(db *gorm.DB, rds *redis.Client, cfg *config.Config) *Handler {
	return &Handler{
		db:    db,
		cfg:   cfg,
		redis: rds,
	}
}

func (h *Handler) SetupRoutes(app *fiber.App) {
	cfg := config.LoadConfig()

	api := app.Group("api")
	v1 := api.Group("v1")

	docs := v1.Group("docs")
	docs.Get("*", fiberSwagger.WrapHandler)

	auth := v1.Group("auth")
	// Публичные маршруты - без проверки JWT
	auth.Post("login", h.login)
	auth.Post("register", h.register)

	// Защищенные маршруты - с middleware JWT
	authProtected := auth.Group("/")
	authProtected.Use(JWTMiddleware(cfg.SecretKey))
	authProtected.Get("get-me", h.getMe)
	authProtected.Post("password/change", h.passwordChange)
	authProtected.Post("refresh", h.refresh)
}

func (h *Handler) ResetPassword(user *model.User, newPasswordHash string) error {
	tx := h.db.Model(&user).Update("Password", newPasswordHash)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
