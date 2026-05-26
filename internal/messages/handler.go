package messages

import (
	"net/http"

	"chitchat/internal/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(e *echo.Echo) {
	msgs := e.Group("/conversations/:conversationId/messages")
	msgs.Use(auth.AuthMiddleware)
	msgs.POST("", h.sendMessage)
}

func (h *Handler) sendMessage(c *echo.Context) error {
	var payload SendMessagePayload

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}

	userSession := c.Get("user").(auth.SessionStore)

	conversationID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid conversation ID")
	}

	userID, _ := uuid.Parse(userSession.UserId)
	deviceID, _ := uuid.Parse(userSession.DeviceId)

	msg, err := h.service.SendMessage(c.Request().Context(), conversationID, userID, deviceID, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send message")
	}

	return c.JSON(http.StatusCreated, msg)
}
