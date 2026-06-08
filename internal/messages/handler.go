package messages

import (
	"net/http"
	"time"

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
	msgs := e.Group("/conversations/:conversationId/messages")
	msgs.Use(auth.AuthMiddleware)
	msgs.POST("", h.sendMessage)
	msgs.GET("", h.getMessages)

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

	for _, envelope := range payload.Envelopes {
		e := MessageEnvelope{
			ConversationID: conversationID,
			MessageID:      msg.ID,
			Context:        envelope.Context,
			SentAt:         msg.SentAt,
			SenderID:       msg.SenderID,
		}
		h.hub.SendEvent(envelope.RecipientDeviceID, ws.EventMessage, e)
	}

	return c.JSON(http.StatusCreated, msg)
}

func (h *Handler) getMessages(c *echo.Context) error {
	userSession := c.Get("user").(auth.SessionStore)
	deviceID, _ := uuid.Parse(userSession.DeviceId)

	conversationID, err := uuid.Parse(c.Param("conversationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid conversation ID")
	}

	timestamp, err := time.Parse(time.RFC3339Nano, c.QueryParam("timestamp"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid timestamp")
	}

	messages, err := h.service.GetMessagesFromTimestamp(c.Request().Context(), conversationID, deviceID, timestamp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get messages")
	}

	return c.JSON(http.StatusOK, messages)
}
