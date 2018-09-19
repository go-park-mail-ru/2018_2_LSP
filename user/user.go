package user

import (
	"2018_2_LSP/utils"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Credentials Structure that stores user credentials for auth
type Credentials struct {
	Email    string
	Password string
}

// User Structure that stores user information retrieved from database or
// entered by user during registration
type User struct {
	Credentials
	ID          int
	Token       string
	Password    string
	FirstName   string
	LastName    string
	Group       int
	DateUpdate  string
	DateUreated string
}

// Register Function that sign ups user
func Register(u User) (User, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		return u, nil
	}
	u.Password = string(hashedPwd)

	_, err = utils.Query("INSERT INTO users (firstName, lastName, email, password) VALUES ($1, $2, $3, $4)", u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return u, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	u.Token, err = token.SignedString([]byte("test"))
	if err != nil {
		return u, err
	}

	return u, nil
}

// Auth Function that authenticates user
func Auth(c Credentials) (User, error) {
	var u User
	// if lenc.Email
	rows, err := utils.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", c.Email)
	if err != nil {
		return u, err
	}
	err = rows.Scan(&u)
	if err != nil {
		return u, errors.New("User not found")
	}

	if !comparePasswords(u.Password, c.Password) {
		return u, errors.New("Wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	u.Token, err = token.SignedString([]byte("test"))
	if err != nil {
		return u, err
	}

	return u, nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
