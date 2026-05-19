package conversations

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
	conv := e.Group("/conversations")
	conv.Use(auth.AuthMiddleware)
	conv.POST("", h.createConversation)
}

func (h *Handler) createConversation(c *echo.Context) error {
	var payload CreateConversationPayload

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}

	userSession := c.Get("user").(auth.SessionStore)

	userID, _ := uuid.Parse(userSession.UserId)

	conv, err := h.service.CreateConversation(c.Request().Context(), userID, payload.Type, payload.ParticipantEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create conversation")
	}

	return c.JSON(http.StatusCreated, conv)
}
