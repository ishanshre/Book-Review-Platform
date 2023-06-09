package helpers

import "golang.org/x/crypto/bcrypt"

// EncryptPassowd return the hashed password stirng and error
func EncryptPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword checks if the hashed passwod and given password matches or not
func CheckPassword(hashedPassword, testPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	return err
}
