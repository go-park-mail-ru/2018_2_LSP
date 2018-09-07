package user

import "errors"

// Credentials Structure that stores user credentials for auth
type Credentials struct {
	Username string
	Password string
}

// User Structure that stores user information retrieved from database or
// entered by user during registration
type User struct {
	Credentials
	firstName   string
	lastName    string
	group       int
	dateUpdate  string
	dateUreated string
}

// Register Function that sign ups user
func Register(u User) error {
	return nil
}

// Auth Function that authenticates user
func Auth(u Credentials) error {
	if u.Username != "test" || u.Password != "test" {
		return errors.New("Wromg username or password")
	}
	return nil
}
