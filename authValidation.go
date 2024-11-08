package validations

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(userPassword, loginPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(loginPassword))
	if err != nil {
		return false
	}
	return true
}
