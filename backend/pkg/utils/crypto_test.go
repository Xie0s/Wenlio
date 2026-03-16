package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "MySecureP@ss1"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == password {
		t.Fatal("HashPassword() returned plaintext password")
	}
}

func TestHashPasswordDifferentHashes(t *testing.T) {
	password := "SamePassword123"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() first call error = %v", err)
	}
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() second call error = %v", err)
	}
	if hash1 == hash2 {
		t.Fatal("HashPassword() produced identical hashes for same password (bcrypt should use random salt)")
	}
}

func TestCheckPasswordCorrect(t *testing.T) {
	password := "CorrectHorse42!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword() returned false for correct password")
	}
}

func TestCheckPasswordIncorrect(t *testing.T) {
	password := "RightPassword1"
	wrong := "WrongPassword1"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if CheckPassword(wrong, hash) {
		t.Fatal("CheckPassword() returned true for wrong password")
	}
}

func TestCheckPasswordEmptyPassword(t *testing.T) {
	hash, err := HashPassword("notempty")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if CheckPassword("", hash) {
		t.Fatal("CheckPassword() returned true for empty password")
	}
}

func TestCheckPasswordInvalidHash(t *testing.T) {
	if CheckPassword("anypassword", "not-a-valid-bcrypt-hash") {
		t.Fatal("CheckPassword() returned true for invalid hash")
	}
}

func TestHashPasswordEmpty(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword() error = %v, want nil for empty password", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash for empty password")
	}
	if !CheckPassword("", hash) {
		t.Fatal("CheckPassword() failed for empty password with its own hash")
	}
}
