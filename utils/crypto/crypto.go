package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
)

// GenerateSalt generates a random salt
func GenerateSalt() string {
	saltBytes := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, saltBytes)
	if err != nil {
		log.Fatal(err)
	}
	salt := make([]byte, 32)
	hex.Encode(salt, saltBytes)
	return string(salt)
}

// HashPassword hashes a string
func HashPassword(password, salt string) string {
	hashedPasswordBytes, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		log.Fatal("Unable to hash password")
	}
	hashedPassword := make([]byte, 64)
	hex.Encode(hashedPassword, hashedPasswordBytes)
	return string(hashedPassword)
}

func GenerateToken(userId int) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"nbf":     time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)

	// b := make([]byte, 64)
	// _, err := rand.Read(b)
	// if err != nil {
	// 	return "", err
	// }
	// str := base64.URLEncoding.EncodeToString(b)
	return tokenString, err
}
