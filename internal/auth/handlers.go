package auth

import (
	"log/slog"
	"net/http"
	"net/netip"

	"github.com/labstack/echo/v5"
)

const (
	SessionCookieName = "session"
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
	auth := e.Group("/auth")
	auth.POST("/send-magic-link", h.sendMagicLink)
	auth.POST("/verify-magic-link", h.verifyMagicLink)
}

func (h *Handler) sendMagicLink(c *echo.Context) error {
	var payload SendMagicLinkPayload

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Input")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}
	addr := getClientIP(c)
	magic_link_session, err := h.service.SendMagicLink(c.Request().Context(), payload.Email, payload.Pubkey, addr, c.Request().UserAgent())
	if err != nil {
		return ErrInternal
	}

	return c.JSON(http.StatusOK, magic_link_session)
}

func (h *Handler) verifyMagicLink(c *echo.Context) error {
	var payload VerifyMagicLinkPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Input")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}
	// addr := getClientIP(c)

	return nil
}

func getClientIP(c *echo.Context) netip.Addr {
	ip := c.RealIP()
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		// Return dummy ip.
		return netip.MustParseAddr("127.0.0.5")
	}
	return addr
}
