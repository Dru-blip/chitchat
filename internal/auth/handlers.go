package auth

import (
	"errors"
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
		h.logger.Error("failed to send magic link", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send magic link")
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
	addr := getClientIP(c)

	sessionId, userId, err := h.service.VerifyMagicLink(c.Request().Context(), payload.Token, addr, c.Request().UserAgent())

	if err != nil {
		return h.mapAuthError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"sessionId": sessionId,
		"userId":    userId.String(),
	})
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

// TODO: should find another way to map errors
func (h *Handler) mapAuthError(err error) error {
	switch {
	case errors.Is(err, ErrInvalidMagicLink),
		errors.Is(err, ErrMagicLinkUsed),
		errors.Is(err, ErrMagicLinkRevoked),
		errors.Is(err, ErrMagicLinkExpired):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrInvalidRequest):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrInvalidSession),
		errors.Is(err, ErrSessionRevoked),
		errors.Is(err, ErrNoSession):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ErrInternal):
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	default:
		h.logger.Error("verify_magic_link_failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify magic link")
	}
}
