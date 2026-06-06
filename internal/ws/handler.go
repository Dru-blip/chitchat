package ws

import (
	"chitchat/internal/auth"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) Register(echo *echo.Echo) {
	ws := echo.Group("/ws")
	ws.Use(auth.AuthMiddleware)

	ws.GET("", h.Connect)

}

func (h *Handler) Connect(c *echo.Context) error {
	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	userSession := c.Get("user").(auth.SessionStore)
	h.hub.Register(userSession.DeviceId, conn)
	h.hub.SendConnectionEvent(userSession.DeviceId)
	go h.hub.handleClient(conn, userSession.DeviceId)

	return nil
}
