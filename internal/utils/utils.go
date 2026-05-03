package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

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

func SHA256(payload string) string {
	h := sha256.New()
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func EncryptAES(text string, key string) (string, int, error) {
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", 0, err
	}
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", 0, err
	}
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", 0, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nonceSize, nil
}

func DecryptAES(key string, cipherText string, nonce string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	plainText, err := gcm.Open(nil, []byte(nonce), decodedCipherText, nil)
	if err != nil {
		return "", err
	}

	//TODO: should return bytes instead of string
	return base64.StdEncoding.EncodeToString(plainText), nil
}

func GenerateMagicLinkToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func WriteCookie(c *echo.Context, name string, value string, expiry time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = expiry
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}
