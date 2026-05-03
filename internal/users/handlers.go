package users

import (
	"chitchat/internal/auth"
	"chitchat/internal/utils"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	service     Service
	logger      *slog.Logger
	redisClient *redis.Client
}

func NewHandler(service Service, logger *slog.Logger, redisClient *redis.Client) *Handler {
	return &Handler{
		service:     service,
		logger:      logger,
		redisClient: redisClient,
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

	userSession := c.Get("user").(auth.SessionStore)

	_, err := h.service.OnboardUser(c.Request().Context(), payload.Name, payload.Password, payload.Image, userSession.Email)
	if err != nil {
		return err
	}

	utils.WriteCookie(c, "onboarding", "", time.Now().AddDate(0, 0, -2))
	if err := auth.RemoveOnboardingToken(c.Request().Context(), h.redisClient, userSession.UserId); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "onboarded successfully",
	})
}
