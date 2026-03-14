package auth

import (
	"chitchat/internal/db"
	"chitchat/internal/db/sqlc"
	"chitchat/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid/v5"
	"github.com/wneessen/go-mail"
)

type Handler struct {
	store *db.Store
}

type OtpData struct {
	Otp string
}

func NewHandler(store *db.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Register() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/send-otp", h.sendOtp)
	r.Get("/ping", h.Ping)
	return r
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, 200, "Pong")
}

func (h *Handler) sendOtp(w http.ResponseWriter, r *http.Request) {
	var payload SendOtpPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	otp, err := utils.GenerateOTPCode(6)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	message := mail.NewMsg()
	if err = message.From(os.Getenv("SMTP_FROM")); err != nil {
		log.Print(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err = message.To(payload.Email); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	message.Subject("Your OTP Code")
	message.SetBodyString(mail.TypeTextPlain, fmt.Sprintf("otp : %s\n", otp))

	client, err := mail.NewClient(os.Getenv("SMTP_HOST"),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(os.Getenv("SMTP_USER")), mail.WithPassword(os.Getenv("SMTP_PASS")),
	)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	id, err := uuid.NewV4()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	h.store.Queries.CreateOtpSession(context.Background(), sqlc.CreateOtpSessionParams{
		ID:        id.String(),
		Email:     payload.Email,
		Code:      otp,
		Pubkey:    payload.Pubkey,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	//TODO: move this to a background worker
	if err := client.DialAndSend(message); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, 200, SendOtpResponse{
		Challenge: "insert random string",
		Message:   "OTP sent successfully",
	})
}
