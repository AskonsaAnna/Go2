package handlers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CheckPassword(storedPassword, inputPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))

	return err == nil
}