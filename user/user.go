package user

import (
	"github.com/go-park-mail-ru/2018_2_LSP/utils"
)

// User Structure that stores user information retrieved from database or
// entered by user during registration
type User struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	OldPassword string `json:"oldpassword"`
	ID          int    `json:"id"`
	Token       string `json:"token"`
	Username    string `json:"username"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Group       int    `json:"group"`
	DateUpdate  string `json:"dateupdated"`
	DateUreated string `json:"datecreated"`
}

// GetUserByID returns all user information by ID
func GetUserByID(id int) (User, error) {
	var u User
	rows, err := utils.Query("SELECT email, first_name, last_name, username FROM users WHERE id = $1", id)
	if err != nil {
		return u, err
	}

	defer rows.Close()
	rows.Next()

	err = rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Username)
	return u, err
}

// ValidateUserPassword validates user password
func ValidateUserPassword(password string, id int) (bool, error) {
	row, err := utils.Query("SELECT password FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	var hashedPassword string

	defer row.Close()
	row.Next()

	err = row.Scan(&hashedPassword)
	if err != nil {
		return false, err
	}
	return utils.ComparePasswords(hashedPassword, password), nil
}
