package user

import (
	"2018_2_LSP/utils"
	"errors"
	"log"
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
	Username    string
	Password    string
	FirstName   string
	LastName    string
	Group       int
	DateUpdate  string
	DateUreated string
}

type RegisterError struct {
	s    string
	code int
}

func (e *RegisterError) Error() string {
	return e.s
	// return fmt.Sprintf("radius %0.2f: %s", e.radius, e.err)
}

func (e *RegisterError) Code() int {
	return e.code
	// return fmt.Sprintf("radius %0.2f: %s", e.radius, e.err)
}

// Register Function that sign ups user
func Register(u User) (User, error) {

	err := validateStringForEmptiness(u.Password, u.Email, u.Username)
	if err != nil {
		return u, err
	}

	u.Password = hashAndSalt([]byte(u.Password))

	rows, err := utils.Query("SELECT EXISTS (SELECT * FROM user WHERE email = $1 LIMIT 1) AS `email`, EXISTS (SELECT * FROM user WHERE username = $2 LIMIT 1) AS `username`", u.Email, u.Username)
	if err != nil {
		return u, err
	}

	rows.Next()
	emailTaken := false
	usernameTaken := false
	err = rows.Scan(&emailTaken, &usernameTaken)
	if err != nil {
		return u, err
	}
	if emailTaken {
		return u, &RegisterError{"Username is already taken", 1}
	}
	if usernameTaken {
		return u, &RegisterError{"Email is already taken", 2}
	}

	rows, err = utils.Query("INSERT INTO users (first_name, last_name, email, password, username) VALUES ($1, $2, $3, $4, $5) RETURNING id;", u.FirstName, u.LastName, u.Email, u.Password, u.Username)
	if err != nil {
		return u, err
	}

	rows.Next()
	err = rows.Scan(&u.ID)
	if err != nil {
		return u, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        u.ID,
		"generated": time.Now(),
	})
	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
	if err != nil {
		return u, err
	}

	return u, nil
}

// Auth Function that authenticates user
func Auth(c Credentials) (User, error) {
	var u User

	err := validateStringForEmptiness(c.Email, c.Password)
	if err != nil {
		return u, err
	}

	rows, err := utils.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", c.Email)
	if err != nil {
		return u, err
	}
	rows.Next()
	err = rows.Scan(&u.ID, &u.Password)
	if err != nil {
		return u, errors.New("User not found")
	}

	if !comparePasswords(u.Password, c.Password) {
		return u, errors.New("Wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        u.ID,
		"generated": time.Now(),
	})
	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
	if err != nil {
		return u, err
	}

	return u, nil
}

func validateStringForEmptiness(strs ...string) error {
	for _, s := range strs {
		if len(s) == 0 {
			return errors.New("Found empty parameter")
		}
	}
	return nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
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
