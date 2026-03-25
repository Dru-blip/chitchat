package auth

import (
	"chitchat/internal/db/sqlc"
	"chitchat/internal/utils"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/wneessen/go-mail"
)

type Handler struct {
	repo OtpSessionRespository
}

type OtpData struct {
	Otp string
}

func NewHandler(repo *sqlc.Queries) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	authGroup := e.Group("/auth")

	authGroup.POST("/send-otp", h.sendOtp)
	authGroup.POST("/verify-otp", h.verifyOtp)
	authGroup.GET("/ping", h.Ping)
}

func (h *Handler) Ping(c *echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func (h *Handler) sendOtp(c *echo.Context) error {
	var payload SendOtpPayload

	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Bad Input")
	}

	if err := c.Validate(&payload); err != nil {
		return err
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		return err
	}

	message := mail.NewMsg()
	if err = message.From(os.Getenv("SMTP_FROM")); err != nil {
		//TODO: change error messages. instead of sending same message.
		return err
	}

	if err = message.To(payload.Email); err != nil {
		return err
	}

	message.Subject("Your OTP")

	plain := fmt.Sprintf(`
Your OTP is: %s
This code is valid for 5 minutes.
chitchat
`, otp)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height:1.5;">
  <p>Your OTP is:</p>
  <h2 style="letter-spacing:2px;">%s</h2>
  <p>This code is valid for 5 minutes.</p>
  <p>chitchat</p>
</body>
</html>
`, otp)

	message.SetBodyString(mail.TypeTextPlain, plain)
	message.AddAlternativeString(mail.TypeTextHTML, html)

	//TODO: should make this a global client instance
	client, err := mail.NewClient(
		os.Getenv("SMTP_HOST"),
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(os.Getenv("SMTP_USER")),
		mail.WithPassword(os.Getenv("SMTP_PASS")),
	)

	if err != nil {
		return err
	}

	challenge := rand.Text()
	challengeHashed := utils.SHA256(challenge)

	otpSession, err := h.repo.CreateOtpSession(c.Request().Context(), sqlc.CreateOtpSessionParams{
		Email:     payload.Email,
		Code:      otp,
		Pubkey:    payload.Pubkey,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Challenge: challengeHashed,
	})

	if err != nil {
		//TODO: should log database errors
		return c.String(http.StatusInternalServerError, "Failed to send Otp")
	}

	go func() {
		if err := client.DialAndSend(message); err != nil {
			log.Printf("Failed to send email to %s", payload.Email)
		}
	}()

	return c.JSON(http.StatusOK, SendOtpResponse{
		Id:        otpSession.ID.String(),
		Challenge: challenge,
		Message:   "OTP sent successfully",
	})
}

func (h *Handler) verifyOtp(c *echo.Context) error {
	var payload VerifyOtpPayload

	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Bad Input")
	}

	otpSession, err := h.repo.GetOtpSessionById(c.Request().Context(), uuid.MustParse(payload.Id))

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to verify Otp")
	}

	if otpSession.ExpiresAt.Before(time.Now()) {
		return c.String(http.StatusUnauthorized, "Session timed out")
	}

	if otpSession.Email != payload.Email || otpSession.Code != payload.Code {
		return c.String(http.StatusUnauthorized, "Invalid OTP")
	}

	plainText, err := utils.DecryptAES(otpSession.Pubkey, payload.Challenge, payload.Nonce)

	if subtle.ConstantTimeCompare([]byte(plainText), []byte(otpSession.Challenge)) != 1 {
		return c.String(http.StatusUnauthorized, "Invalid OTP")
	}

	//TODO: verification of digital signature using pubkey
	// and generate a session for the client.
	return c.JSON(http.StatusOK, "OTP verified successfully")
}
