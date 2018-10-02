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
	row.Next()
	err = row.Scan(&hashedPassword)
	if err != nil {
		return false, err
	}
	return utils.ComparePasswords(hashedPassword, password), nil
}

// Register Function that sign ups user
// func Register(u User) (User, error) {

// 	err := utils.ValidateStringForEmptiness(u.Password, u.Email, u.Username)
// 	if err != nil {
// 		return u, err
// 	}

// 	u.Password = utils.HashAndSalt([]byte(u.Password))

// 	rows, err := utils.Query("SELECT EXISTS (SELECT * FROM users WHERE email = $1 LIMIT 1) AS email, EXISTS (SELECT * FROM users WHERE username = $2 LIMIT 1) AS username", u.Email, u.Username)
// 	if err != nil {
// 		return u, err
// 	}

// 	rows.Next()
// 	emailTaken := false
// 	usernameTaken := false
// 	err = rows.Scan(&emailTaken, &usernameTaken)
// 	if err != nil {
// 		return u, err
// 	}
// 	if emailTaken && usernameTaken {
// 		return u, &RegisterError{"Is already taken", 1}
// 	}
// 	if emailTaken {
// 		return u, &RegisterError{"Is already taken", 2}
// 	}
// 	if usernameTaken {
// 		return u, &RegisterError{"Is already taken", 3}
// 	}

// 	rows, err = utils.Query("INSERT INTO users (first_name, last_name, email, password, username) VALUES ($1, $2, $3, $4, $5) RETURNING id;", u.FirstName, u.LastName, u.Email, u.Password, u.Username)
// 	if err != nil {
// 		return u, err
// 	}

// 	rows.Next()
// 	err = rows.Scan(&u.ID)
// 	if err != nil {
// 		return u, err
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":        u.ID,
// 		"generated": time.Now(),
// 	})
// 	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
// 	if err != nil {
// 		return u, err
// 	}

// 	return u, nil
// }

// // Auth Function that authenticates user
// func Auth(c Credentials) (User, error) {
// 	var u User

// 	err := utils.ValidateStringForEmptiness(c.Email, c.Password)
// 	if err != nil {
// 		return u, err
// 	}

// 	rows, err := utils.Query("SELECT id, password FROM users WHERE email = $1 LIMIT 1", c.Email)
// 	if err != nil {
// 		return u, err
// 	}
// 	rows.Next()
// 	err = rows.Scan(&u.ID, &u.Password)
// 	if err != nil {
// 		return u, errors.New("User not found")
// 	}

// 	if !utils.ComparePasswords(u.Password, c.Password) {
// 		return u, errors.New("Wrong password")
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":        u.ID,
// 		"generated": time.Now(),
// 	})
// 	u.Token, err = token.SignedString([]byte("HeAdfasdf3ref&^%$Dfrtgauyhia"))
// 	if err != nil {
// 		return u, err
// 	}

// 	return u, nil
// }
