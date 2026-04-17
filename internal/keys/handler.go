package keys

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
	keys := e.Group("/keys")
	keys.Use(auth.AuthMiddleware)

	keys.POST("/", h.UploadPrekeys)
	keys.POST("/:userId", h.GetkeyBundle)
	// keys.DELETE("/", h.ExtrackPrekey)

}

func (h *Handler) UploadPrekeys(c *echo.Context) error {
	var payload UploadPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}
	if err := c.Validate(&payload); err != nil {
		return err
	}
	user := c.Get("user").(auth.SessionStore)

	err := h.service.UploadPrekeys(
		c.Request().Context(),
		prekeyUpload{
			DeviceID:    user.DeviceId,
			PrekeyIds:   payload.PrekeyIds,
			Prekeys:     payload.Prekeys,
			Signature:   payload.SignedPreKey.Signature,
			SignedKeyID: payload.SignedPreKey.ID,
			SignedKey:   payload.SignedPreKey.Key,
		},
	)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"message": "successfully stored prekeys",
	})
}

func (h *Handler) GetkeyBundle(c *echo.Context) error {
	userId := c.Param("userId")

	key_bundle, err := h.service.GetKeyBundle(c.Request().Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"bundle": key_bundle,
	})
}
