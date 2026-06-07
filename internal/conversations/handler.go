package conversations

import (
	"net/http"

	"chitchat/internal/auth"
	"chitchat/internal/ws"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
	hub     *ws.Hub
}

func NewHandler(service Service, hub *ws.Hub) *Handler {
	return &Handler{service: service, hub: hub}
}

func (h *Handler) Register(e *echo.Echo) {
	conv := e.Group("/conversations")
	conv.Use(auth.AuthMiddleware)
	conv.POST("", h.createConversation)
	conv.GET("", h.getConversations)
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

	var participantID string

	for _, participant := range conv.Participants {
		if userID != participant.UserID {
			participantID = participant.UserID.String()
		}
	}

	h.hub.SendToUser(participantID, ws.EventNewConversation, conv)

	return c.JSON(http.StatusCreated, conv)
}

func (h *Handler) getConversations(c *echo.Context) error {
	userSession := c.Get("user").(auth.SessionStore)

	userID, _ := uuid.Parse(userSession.UserId)

	conversations, err := h.service.GetConversationsByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch conversations")
	}

	return c.JSON(http.StatusOK, conversations)
}
