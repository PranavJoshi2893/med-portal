package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAccessToken_VerifyAccessToken(t *testing.T) {
	key := "test-secret-key"
	userID, _ := uuid.NewV7()
	role := "user"

	token, err := GenerateAccessToken(key, userID, role)
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := VerifyAccessToken(key, token)
	if err != nil {
		t.Fatalf("VerifyAccessToken: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID: got %v want %v", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("Role: got %q want %q", claims.Role, role)
	}
}

func TestGenerateRefreshToken_VerifyRefreshToken(t *testing.T) {
	key := "test-refresh-key"
	userID, _ := uuid.NewV7()
	role := "admin"

	token, err := GenerateRefreshToken(key, userID, role)
	if err != nil {
		t.Fatalf("GenerateRefreshToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := VerifyRefreshToken(key, token)
	if err != nil {
		t.Fatalf("VerifyRefreshToken: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID: got %v want %v", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("Role: got %q want %q", claims.Role, role)
	}
}

func TestVerifyAccessToken_Invalid(t *testing.T) {
	key := "test-key"
	userID, _ := uuid.NewV7()

	token, _ := GenerateAccessToken(key, userID, "user")

	tests := []struct {
		name string
		key  string
		tok  string
	}{
		{"wrong key", "wrong-key", token},
		{"empty token", key, ""},
		{"garbage", key, "not.a.jwt"},
		{"tampered", key, token + "x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyAccessToken(tt.key, tt.tok)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestVerifyRefreshToken_Invalid(t *testing.T) {
	key := "test-key"
	userID, _ := uuid.NewV7()

	token, _ := GenerateRefreshToken(key, userID, "user")

	tests := []struct {
		name string
		key  string
		tok  string
	}{
		{"wrong key", "wrong-key", token},
		{"empty token", key, ""},
		{"garbage", key, "not.a.jwt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyRefreshToken(tt.key, tt.tok)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestRefreshToken_UniqueJti(t *testing.T) {
	key := "test-key"
	userID, _ := uuid.NewV7()

	t1, _ := GenerateRefreshToken(key, userID, "user")
	t2, _ := GenerateRefreshToken(key, userID, "user")

	if t1 == t2 {
		t.Error("refresh tokens should be unique (different jti)")
	}
}

func TestAccessToken_Expiry(t *testing.T) {
	key := "test-key"
	userID, _ := uuid.NewV7()

	token, err := GenerateAccessToken(key, userID, "user")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := VerifyAccessToken(key, token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt should be set")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("token should not be expired yet")
	}
}
