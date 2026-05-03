package auth

import (
	"chitchat/internal/utils"
	"errors"
	"log/slog"
	"net/http"
	"net/netip"
	"time"

	"github.com/alexedwards/scs/v2"
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
		if errors.Is(err, ErrTooManyAttempts) {
			return c.JSON(http.StatusTooManyRequests, magic_link_session)
		}
		return err
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

	magic_session, err := h.service.VerifyMagicLink(c.Request().Context(), payload.Token, addr, c.Request().UserAgent())

	if err != nil {
		return err
	}

	user, created, err := h.service.GetOrCreateUser(c.Request().Context(), magic_session.Email, magic_session.Pubkey)

	if err != nil {
		return err
	}

	//TODO: device creation and prekeys setup
	//TODO: should receive windows and client fingerprints from payload
	device, err := h.service.GetOrCreateDevice(c.Request().Context(), user.ID, magic_session.Pubkey, "Windows 11", c.Request().UserAgent())

	if err != nil {
		return err
	}

	session_manager := c.Get("_session").(*scs.SessionManager)
	session_manager.Put(c.Request().Context(), "user", SessionStore{
		Email:    user.Email,
		Pubkey:   device.Pubkey,
		UserId:   user.ID.String(),
		DeviceId: device.ID.String(),
	})

	if created {
		onboardingToken, err := SetOnboardingToken(c.Request().Context(), h.redisClient, user)
		if err != nil {
			return err
		}

		utils.WriteCookie(c, "onboarding", onboardingToken, time.Now().AddDate(1, 0, 0))
	}

	return c.JSON(http.StatusOK, map[string]any{
		"userId":   user.ID.String(),
		"device":   device.Os,
		"onboard":  created,
		"redirect": true,
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

func (h *Handler) me(c *echo.Context) error {
	user := c.Get("user").(SessionStore)
	return c.JSON(http.StatusOK, user)
}
