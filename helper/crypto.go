package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// func HashPassword(password string) string {
// 	hasher := sha256.New()
// 	hasher.Write([]byte(password))
// 	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
// 	return hashedPassword
// }

// func VerifyPassword(hashedPassword, password string) bool {
// 	return hashedPassword == HashPassword(password)
// }

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies if a provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
