package users

import (
	"chitchat/internal/auth"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	users := e.Group("/users")
	users.Use(auth.AuthMiddleware)
	users.PATCH("/onboard", h.onboard)
}

func (h *Handler) onboard(c *echo.Context) error {
	var payload OnboardUserPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}

	_, err := h.service.OnboardUser(c.Request().Context(), payload.Name, payload.Password, payload.Image, payload.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "onboarded successfully",
	})
}
