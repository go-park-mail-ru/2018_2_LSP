package utils

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidateStringForEmptiness(strs ...string) error {
	for _, s := range strs {
		if len(s) == 0 {
			return errors.New("Found empty parameter")
		}
	}
	return nil
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
