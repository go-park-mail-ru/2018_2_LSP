package user

import (
	"errors"
	"time"

	"github.com/go-park-mail-ru/2018_2_LSP/utils"

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
	FirstName   string
	LastName    string
	Group       int
	DateUpdate  string
	DateUreated string
}

type RegisterError struct {
	Field   string
	Message string
}

func (e *RegisterError) Error() string {
	return e.Message
}

func validateAuthInput(u *User) []error {
	errs := make([]error, 0)

	if len(u.Password) == 0 {
		errs = append(errs, &RegisterError{"password", "No password was specified"})
	}
	if len(u.Email) == 0 {
		errs = append(errs, &RegisterError{"email", "No email was specified"})
	}
	return errs
}

func validateRegisterInput(u *User) []error {
	errs := make([]error, 0)

	if len(u.Password) == 0 {
		errs = append(errs, &RegisterError{"password", "No password was specified"})
	}
	if len(u.Email) == 0 {
		errs = append(errs, &RegisterError{"email", "No email was specified"})
	}
	if len(u.Username) == 0 {
		errs = append(errs, &RegisterError{"username", "No username was specified"})
	}
	return errs
}

func validateRegisterUnique(u *User) []error {
	errs := make([]error, 0)
	rows, err := utils.Query("SELECT EXISTS (SELECT * FROM users WHERE email = $1 LIMIT 1) AS email, EXISTS (SELECT * FROM users WHERE username = $2 LIMIT 1) AS username", u.Email, u.Username)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	rows.Next()
	emailTaken := false
	usernameTaken := false
	err = rows.Scan(&emailTaken, &usernameTaken)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	if emailTaken {
		errs = append(errs, &RegisterError{"email", "Email is already taken"})
	}
	if usernameTaken {
		errs = append(errs, &RegisterError{"username", "Username is already taken"})
	}

	return nil
}

func (u *User) createUser() []error {
	errs := make([]error, 0)
	rows, err := utils.Query("INSERT INTO users (first_name, last_name, email, password, username) VALUES ($1, $2, $3, $4, $5) RETURNING id;", u.FirstName, u.LastName, u.Email, u.Password, u.Username)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	rows.Next()
	err = rows.Scan(&u.ID)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	return nil
}

func (u *User) generateToken() []error {
	var err error
	errs := make([]error, 0)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        u.ID,
		"generated": time.Now(),
	})
	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	return nil
}

// Register Function that sign ups user
func (u *User) Register() []error {
	var err error

	if errs := validateRegisterInput(u); errs != nil {
		return errs
	}

	u.Password, err = hashAndSalt(u.Password)
	if err != nil {
		errs := make([]error, 0)
		errs = append(errs, err)
		return errs
	}

	if errs := validateRegisterUnique(u); errs != nil {
		return errs
	}

	if errs := u.createUser(); errs != nil {
		return errs
	}

	if errs := u.generateToken(); errs != nil {
		return errs
	}

	return nil
}

// Auth Function that authenticates user
func (u *User) Auth(c Credentials) []error {
	if errs := validateAuthInput(u); errs != nil {
		return errs
	}

	rows, err := utils.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", c.Email)
	if err != nil {
		errs := make([]error, 0)
		errs = append(errs, err)
		return errs
	}
	rows.Next()

	if err := rows.Scan(&u.ID, &u.Password); err != nil {
		errs := make([]error, 0)
		errs = append(errs, errors.New("User not found"))
		return errs
	}

	if !comparePasswords(u.Password, c.Password) {
		errs := make([]error, 0)
		errs = append(errs, errors.New("Wrong password"))
		return errs
	}

	if errs := u.generateToken(); errs != nil {
		return errs
	}

	return nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
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
