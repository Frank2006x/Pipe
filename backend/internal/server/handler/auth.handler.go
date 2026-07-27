package handler

import (
	"Frank2006x/Pipe/internal/service"
	"Frank2006x/Pipe/internal/util"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) GetRedirctLink(c fiber.Ctx) error {

	state := util.GenerateRandomState()

	authURL, err := h.AuthService.GetAuthURL(state)
	if err != nil {
		log.Infof("Failed to get GitHub auth URL: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get GitHub auth URL")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
		MaxAge:   60 * 10, // state is only valid for 10 minutes
	})

	return c.JSON(fiber.Map{
		"url": authURL,
	})
}

func (h *AuthHandler) Callback(c fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing state parameter")
	}

	cookieState := c.Cookies("oauth_state")
	if cookieState == "" || cookieState != state {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid state parameter (potential CSRF)")
	}

	// Delete state cookie after validation
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
		MaxAge:   -1,
	})

	code := c.Query("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing code parameter")
	}

	token, err := h.AuthService.Callback(c.Context(), code)
	if err != nil {
		log.Infof("Failed to handle callback: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to handle callback")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,   // MUST be true for SameSite: None
		SameSite: "None", // Required for cross-origin/cross-port cookie transfers
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 3,
	})
	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (h *AuthHandler) GetUserInfo(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	user, err := h.AuthService.GetUserInfo(c.Context(), userID)
	if err != nil {
		log.Infof("Failed to get user info: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get user info")
	}

	return c.JSON(user)
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,   // MUST be true for SameSite: None
		SameSite: "None", // Required for cross-origin/cross-port cookie transfers
		Path:     "/",
		MaxAge:   -1,
	})

	return c.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
