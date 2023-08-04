package models

import (
	"github.com/lauderdice/ai-finance-tracker/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func CreateUser(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	sqlStatement := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id`
	var id int
	err = database.DB.QueryRow(sqlStatement, email, string(hashedPassword)).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	sqlStatement := "SELECT id, email, password FROM users WHERE email = $1"
	err := database.DB.QueryRow(sqlStatement, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func GetUserByID(userID int64) (*User, error) {
	user := &User{}
	sqlStatement := "SELECT id, email, password FROM users WHERE id = $1"
	err := database.DB.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
