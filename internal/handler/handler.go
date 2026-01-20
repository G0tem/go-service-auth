package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/G0tem/go-servise-auth/internal"
	"github.com/G0tem/go-servise-auth/internal/config"
	"github.com/G0tem/go-servise-auth/internal/handler/rbac"
	"github.com/G0tem/go-servise-auth/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

type Handler struct {
	rbac        *rbac.RBACLayer
	db          *gorm.DB
	cfg         *config.Config
	userService UserService
	redis       *redis.Client
}

func NewHandler(db *gorm.DB, rbac *rbac.RBACLayer, cfg *config.Config) *Handler {
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr, // Адрес Redis (например, "localhost:6379")
		DB:   cfg.RedisDB,   // Номер базы данных Redis
	})

	log.Println("Successfully connected to Redis")
	return &Handler{
		rbac:        rbac,
		db:          db,
		cfg:         cfg,
		userService: NewHTTPUserService(cfg.UserServiceBaseUrl),
		redis:       redisClient,
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

func (h *Handler) GetJWT(user *model.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id":     user.ID.String(),
		"username":    user.Username,
		"email":       user.Email,
		"role":        user.Role.Name,
		"permissions": h.GetPermissions(user),
		"exp":         time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.cfg.SecretKey))

	return t, err
}

func (h *Handler) GetPermissions(user *model.User) []string {
	permissions, err := h.rbac.GetRolePermissions(user.Role.Name)
	if err != nil {
		return []string{}
	}
	return internal.Mapping(permissions, func(x model.UserPermission) string {
		return fmt.Sprintf("%v:%v", x.Model, x.Action)
	})
}

func (h *Handler) GetPublicUrl() string {
	return h.cfg.PublicUrl
}

func (h *Handler) GetPublicErrorUrl() string {
	return h.cfg.PublicErrorUrl
}
