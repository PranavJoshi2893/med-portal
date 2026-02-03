package encrypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12 // Adjust based on your performance needs
)

type PasswordHasher struct {
	pepper []byte // Load from secrets manager
}

func NewPasswordHasher(pepper string) *PasswordHasher {
	return &PasswordHasher{
		pepper: []byte(pepper),
	}
}

// HashPassword hashes password with HMAC-SHA384 + bcrypt
func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Step 1: HMAC-SHA384 with pepper
	peppered := ph.hmacSHA384(password)

	// Step 2: Base64 encode (bcrypt has 72 byte limit)
	encoded := base64.StdEncoding.EncodeToString(peppered)

	// Step 3: bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(encoded), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword verifies password against hash
func (ph *PasswordHasher) VerifyPassword(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}

	// Step 1: HMAC-SHA384 with pepper
	peppered := ph.hmacSHA384(password)

	// Step 2: Base64 encode
	encoded := base64.StdEncoding.EncodeToString(peppered)

	// Step 3: Compare with bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(encoded))
	return err == nil
}

// hmacSHA384 applies HMAC-SHA384
func (ph *PasswordHasher) hmacSHA384(password string) []byte {
	h := hmac.New(sha512.New384, ph.pepper)
	h.Write([]byte(password))
	return h.Sum(nil)
}

// HashToken returns a SHA-256 hex hash of the token for secure storage.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// GeneratePepper generates a random pepper (run once, store securely)
func GeneratePepper() (string, error) {
	pepper := make([]byte, 32)
	_, err := rand.Read(pepper)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pepper), nil
}
