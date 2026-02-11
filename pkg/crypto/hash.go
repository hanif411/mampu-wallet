package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(passwordInput, passwordDB string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(passwordInput))
	return err == nil
}
