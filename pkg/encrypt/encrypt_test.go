package encrypt

import (
	"testing"
)

func TestHashToken(t *testing.T) {
	token := "my-refresh-token"
	hash1 := HashToken(token)
	hash2 := HashToken(token)

	if hash1 != hash2 {
		t.Error("HashToken should be deterministic")
	}
	if len(hash1) != 64 {
		t.Errorf("SHA-256 hex hash should be 64 chars, got %d", len(hash1))
	}
	if hash1 == token {
		t.Error("hash should not equal original token")
	}
}

func TestHashToken_DifferentInputs(t *testing.T) {
	h1 := HashToken("token1")
	h2 := HashToken("token2")
	if h1 == h2 {
		t.Error("different inputs should produce different hashes")
	}
}

func TestHashPassword_VerifyPassword(t *testing.T) {
	hasher := NewPasswordHasher("pepper")
	password := "myPassword123!"

	hash, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" || hash == password {
		t.Error("hash should be non-empty and different from password")
	}

	if !hasher.VerifyPassword(password, hash) {
		t.Error("VerifyPassword: correct password should verify")
	}
	if hasher.VerifyPassword("wrong", hash) {
		t.Error("VerifyPassword: wrong password should not verify")
	}
}

func TestHashPassword_Empty(t *testing.T) {
	hasher := NewPasswordHasher("pepper")
	_, err := hasher.HashPassword("")
	if err == nil {
		t.Error("expected error for empty password")
	}
}

func TestVerifyPassword_Empty(t *testing.T) {
	hasher := NewPasswordHasher("pepper")
	if hasher.VerifyPassword("", "hash") {
		t.Error("empty password should not verify")
	}
	if hasher.VerifyPassword("pass", "") {
		t.Error("empty hash should not verify")
	}
}
