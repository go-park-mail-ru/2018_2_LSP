package user

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2018_2_LSP/utils"

	jwt "github.com/dgrijalva/jwt-go"
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
	FirstName   string
	LastName    string
	Group       int
	DateUpdate  string
	DateUreated string
}

// Register Function that sign ups user
func (u *User) Register() error {
	var err error
	// TODO чуть поправить валидацию
	if err := validateRegisterUnique(u); err != nil {
		return err
	}

	if u.Password, err = hashPassword(u.Password); err != nil {
		return nil
	}

	if err := u.createUser(); err != nil {
		return err
	}

	if err := u.generateToken(); err != nil {
		return err
	}

	return nil
}

// Auth Function that authenticates user
func (u *User) Auth(c Credentials) error {
	rows, err := utils.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", c.Email)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()

	if err := rows.Scan(&u.ID, &u.Password); err != nil {
		return err
	}

	if !validatePassword(u.Password, c.Password) {
		return err
	}

	if err := u.generateToken(); err != nil {
		return err
	}

	return nil
}

func validateRegisterUnique(u *User) error {
	rows, err := utils.Query("SELECT EXISTS (SELECT * FROM users WHERE email = $1 LIMIT 1) AS email, EXISTS (SELECT * FROM users WHERE username = $2 LIMIT 1) AS username", u.Email, u.Username)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()

	emailTaken, usernameTaken := false, false
	if err = rows.Scan(&emailTaken, &usernameTaken); err != nil {
		return err
	}

	if emailTaken {
		return errors.New("Email is already taken")
	}
	if usernameTaken {
		return errors.New("Username is already taken")
	}

	return nil
}

func (u *User) createUser() error {
	rows, err := utils.Query("INSERT INTO users (first_name, last_name, email, password, username) VALUES ($1, $2, $3, $4, $5) RETURNING id;", u.FirstName, u.LastName, u.Email, u.Password, u.Username)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()

	err = rows.Scan(&u.ID)
	return err
}

func (u *User) generateToken() error {
	var err error
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        u.ID,
		"generated": time.Now(),
	})
	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
	return err
}
