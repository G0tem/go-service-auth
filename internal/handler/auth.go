package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/G0tem/go-servise-auth/internal"
	"github.com/G0tem/go-servise-auth/internal/model"
	"github.com/G0tem/go-servise-auth/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Login
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.LoginRequest true "username or email"
// @Success 200 {object} types.LoginSuccessResponse
// @Failure 401 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/login [post]
func (h *Handler) login(c *fiber.Ctx) error {
	var user *model.User
	input := new(types.LoginRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on login request",
			Error:   err.Error(),
		})
	}

	err := *new(error)

	if validateEmail(input.Identity) {
		user, err = h.getUserByEmail(input.Identity)
	} else {
		user, err = h.getUserByUsername(input.Identity)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   err.Error(),
		})
	} else if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(types.FailureResponse{
			Status:  "error",
			Message: "Invalid identity or password",
		})
	}

	if !CheckPasswordHash(input.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(types.FailureResponse{
			Status:  "error",
			Message: "Invalid identity or password",
		})
	}

	token, err := h.GetJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.LoginSuccessResponse{
		Status: "ok",
		Data:   types.LoginSuccessData{Token: token},
	})
}

// Register
// @Summary Register
// @Description Register
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "username or email"
// @Success 200 {object} types.LoginSuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/register [post]
func (h *Handler) register(c *fiber.Ctx) error {
	input := new(types.RegisterRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on register request",
			Error:   err.Error(),
		})
	}

	if !validateEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureResponse{
			Status:  "error",
			Message: "Invalid email address",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureResponse{
			Status:  "error",
			Message: "Password and confirm password must be same",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
	}

	// Получаем роль по умолчанию
	var defaultRole model.UserRole
	if err := h.db.Where("name = ?", model.DefaultUserRole).First(&defaultRole).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Failed to get default role",
			Error:   err.Error(),
		})
	}

	user := model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		RoleID:       defaultRole.ID,
	}

	tx := h.db.Create(&user)
	if err := tx.Error; err != nil {
		if err.(*pgconn.PgError).Code == "23505" {
			return c.Status(fiber.StatusBadRequest).JSON(types.FailureResponse{
				Status:  "error",
				Message: "Username or email already taken",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error (Create(&user))",
			Error:   err.Error(),
		})
	}

	err = h.SendEmailConfirm(&user, input.Password)
	if err != nil {
		log.Error().Msgf("SendEmailConfirm error %v", err)
	}

	token, err := h.GetJWT(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error (user.GetJWT())",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(types.LoginSuccessResponse{
		Status: "Ok",
		Data:   types.LoginSuccessData{Token: token},
	})
}

// Password Change
// @Summary Password Change
// @Description Password Change
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.PasswordChangeRequest true "password change"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Security ApiKeyAuth
// @Router /auth/password/change [post]
func (h *Handler) passwordChange(c *fiber.Ctx) error {

	var (
		user model.User
	)
	input := new(types.PasswordChangeRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on password reset request",
			Error:   err.Error(),
		})
	}

	if input.NewPassword != input.NewPasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureResponse{
			Status:  "error",
			Message: "New password and new confirm password must be same",
		})
	}

	claims := c.Locals("claims").(*types.JwtClaims)

	tx := h.db.Where("email = ?", claims.Email).First(&user)
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (passwordReset())",
			Error:   err.Error(),
		})
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
	}

	tx = h.db.Model(&user).Update("PasswordHash", string(hashedNewPassword))
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (Update PasswordHash)",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.SuccessResponse{
		Status:  "ok",
		Message: "Password changed.",
	})
}

// Refresh jwt token
// @Summary Refresh jwt token
// @Description Refresh jwt token
// @Tags auth
// @Produce json
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Security ApiKeyAuth
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *fiber.Ctx) error {
	authorization := string(c.Request().Header.Peek("Authorization"))
	tokenString := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (refresh)",
			Error:   err.Error(),
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (email field not present in jwt)",
			Error:   "email field not present in jwt",
		})
	}
	email := claims["email"]

	var (
		user model.User
	)
	// Find user by email
	tx := h.db.Find(&user, "Email = ?", email)
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: fmt.Sprintf("Internal server error (user not found by email %v)", email),
			Error:   err.Error(),
		})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (user not found)",
			Error:   fmt.Sprintf("user with email %v not found", email),
		})
	}

	newToken, err := h.GetJWT(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.LoginSuccessResponse{
		Status: "ok",
		Data:   types.LoginSuccessData{Token: newToken},
	})
}

// Password Reset
// @Summary Password Reset
// @Description Send link to email for reset password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.PasswordChangeResetRequest true "password reset"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/password/reset [post]
func (h *Handler) passwordReset(c *fiber.Ctx) error {
	var (
		user model.User
	)

	input := new(types.PasswordChangeResetRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on password reset request",
			Error:   err.Error(),
		})
	}

	// Find user by email
	tx := h.db.Find(&user, "Email = ?", input.Email)
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (passwordReset())",
			Error:   err.Error(),
		})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (user not found)",
			Error:   fmt.Sprintf("user with email %v not found", input.Email),
		})
	}

	err := h.SendPasswordResetConfirm(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error (h.SendPasswordResetConfirm())",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.SuccessResponse{
		Status:  "ok",
		Message: "Sent email for password reset.",
	})
}

// Password Reset Confirm
// @Summary Password Reset Confirm
// @Description Password Reset Confirm
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.PasswordChangeResetConfirmRequest true "password reset confirm"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/password/reset/confirm [post]
func (h *Handler) passwordResetConfirm(c *fiber.Ctx) error {
	var (
		user model.User
	)

	input := new(types.PasswordChangeResetConfirmRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on password reset request",
			Error:   err.Error(),
		})
	}

	if input.NewPassword != input.NewPasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureResponse{
			Status:  "error",
			Message: "New password and new confirm password must be same",
		})
	}

	jwtToken, err := jwt.Parse(input.Token, nil)
	if jwtToken == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Failed to parse jwt token",
			Error:   err.Error(),
		})
	}

	claimsMap := jwtToken.Claims.(jwt.MapClaims)
	claims := &types.JwtClaims{
		Email: fmt.Sprint(claimsMap["email"]),
	}

	tx := h.db.Where("email = ?", claims.Email).First(&user)
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (passwordReset())",
			Error:   err.Error(),
		})
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
	}

	tx = h.db.Model(&user).Update("PasswordHash", string(hashedNewPassword))
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (Update PasswordHash)",
			Error:   err.Error(),
		})
	}

	log.Info().Msgf("Sending notification to user service about password change with new password")
	// Notify external service about password change with new password
	h.userService.OnUpdateUser(&UserCredentials{
		Email:    user.Email,
		Password: input.NewPassword,
	})

	return c.Status(fiber.StatusOK).JSON(types.SuccessResponse{
		Status:  "ok",
		Message: "New password set, after reset.",
	})
}

// Email Confirm
// @Summary Email Confirm
// @Description Email confirm via emailed link
// @Tags auth
// @Produce json
// @Param string query string false "string valid" minlength(100) maxlength(100) extensions(x-nullable,x-abc=def)
// @Success 302 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/email/confirm [get]
func (h *Handler) emailConfirmFromExternalLink(c *fiber.Ctx) error {
	params := c.Queries()
	log.Info().Msgf("Email confirmation params: %v", params)
	if len(params) != 1 {
		log.Error().Msgf("only one parameter required but given %v", params)
		return c.Redirect(h.GetPublicErrorUrl())
	}

	// Get confirmation request with user data
	var confirmation model.EmailConfirmation
	tx := h.db.Preload("User").First(&confirmation, "Token = ?", internal.FirstKey(params))
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to get confirmation request")
		return c.Redirect(h.GetPublicErrorUrl())
	}

	token := internal.FirstKey(params)

	err := h.ConfirmEmailBySlug(token)
	if err != nil {
		log.Error().Msgf("Internal Server Error (h.ConfirmEmailBySlug()): %v", err)
		return c.Redirect(h.GetPublicErrorUrl())
	}

	log.Info().Msgf("Sending notification to user service about confirmed user with original password")
	// Notify external service about confirmed user with original password
	var ctx = context.Background()
	password, err := h.redis.Get(ctx, token).Result()
	if err != nil {
		log.Error().Msgf("Error retrieving password from Redis: %v", err)
		return c.Redirect(h.GetPublicErrorUrl())
	}
	h.userService.OnCreateUser(&UserCredentials{
		Email:    confirmation.User.Email,
		Password: password,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": h.GetPublicUrl(),
		"status":  "ok",
	})
}

// Email Resend Confirm
// @Summary Email Resend Confirm
// @Description Send email message with email confirmation link
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.ResendEmailConfirmationRequest true "email resend confirm"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Router /auth/register/resend-confirm-email [post]
func (h *Handler) emailResendConfirm(c *fiber.Ctx) error {

	var (
		user model.User
	)

	input := new(types.ResendEmailConfirmationRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Error on password reset request",
			Error:   err.Error(),
		})
	}

	// Find user by email
	tx := h.db.Find(&user, "Email = ?", input.Email)
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (passwordReset())",
			Error:   err.Error(),
		})
	}

	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error (user not found)",
			Error:   fmt.Sprintf("user with email %v not found", input.Email),
		})
	}

	err := h.SendEmailConfirm(&user, "")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal Server Error (h.SendEmailConfirm())",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.SuccessResponse{
		Status:  "ok",
		Message: "Email confirmation sent",
	})
}

// GetMe
// @Summary Get current user info
// @Description Get current user ID and email
// @Tags auth
// @Produce json
// @Success 200 {object} types.GetMeResponse
// @Failure 401 {object} types.FailureResponse
// @Failure 500 {object} types.FailureErrorResponse
// @Security ApiKeyAuth
// @Router /auth/get-me [get]
func (h *Handler) getMe(c *fiber.Ctx) error {
	var user model.User

	claims := c.Locals("claims").(*types.JwtClaims)
	log.Debug().
		Str("email", claims.Email).
		Time("exp", claims.Exp).
		Msg("Attempting to get user")

	tx := h.db.Where("email = ?", claims.Email).First(&user)
	if err := tx.Error; err != nil {
		log.Error().
			Err(err).
			Str("email", claims.Email).
			Msg("Failed to get user from database")
		return c.Status(fiber.StatusInternalServerError).JSON(types.FailureErrorResponse{
			Status:  "error",
			Message: "Internal server error",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.GetMeResponse{
		ID:    user.ID.String(),
		Email: user.Email,
	})
}