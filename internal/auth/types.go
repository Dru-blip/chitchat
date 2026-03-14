package auth

type SendOtpPayload struct {
	Email  string `json:"email"`
	Pubkey string `json:"pubkey"`
}

type SendOtpResponse struct {
	Challenge string `json:"challenge"`
	Message   string `json:"message"`
}

type VerifyOtpPayload struct {
	Id    string `json:"request_id"`
	Email string `json:"email"`
	Code  string `json:"code"`
}
