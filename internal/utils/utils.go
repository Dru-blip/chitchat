package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func GenerateOTPCode(length int) (string, error) {
	const digits = "0123456789"
	max := big.NewInt(int64(len(digits)))
	byteSlice := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		byteSlice[i] = digits[num.Int64()]
	}
	return string(byteSlice), nil
}
