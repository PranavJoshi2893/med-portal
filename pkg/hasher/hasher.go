package hasher

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// bcrypt(base64(hmac-sha384(data:$password, key:$pepper)), $salt, $cost)

func HashPassword(password string) (string, error) {
	cost, err := strconv.Atoi(os.Getenv("COST"))
	if err != nil {
		return "", err
	}
	pepper := os.Getenv("PEPPER")

	h := hmac.New(sha512.New384, []byte(pepper))
	h.Write([]byte(password))
	hmacHash := h.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hmacHash)

	hash, err := bcrypt.GenerateFromPassword([]byte(encoded), cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// func VerifyPassword(password string) string {
// 	cost := os.Getenv("COST")
// 	pepper := os.Getenv("PEPPER")

// 	return ""
// }
